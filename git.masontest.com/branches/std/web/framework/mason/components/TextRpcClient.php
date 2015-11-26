<?php
/**
 * 版本 1.1 
 * 发布时间 2014-05-09
 * 
 * 新 RPC 文本协议客户端实现
 *
 * @author TonyXu <tonycbcd@gmail.com>
 *
 * @usage:
 *  1, 复制或软链接 TextRpcClient.php 到具体的项目目录中
 *  2, 添加 RpcServer 相关配置, 参考: examples/config/debug.php
 *  3, 在 Controller 中添加 RPC 使用代码, 参考下面的例子
 *
 * @example
 *
 *      $userInfo = RpcClient_User_Info::instance();
 *
 *      # case 1
 *      $result = $userInfo->getInfoByUid(100);
 *      if (!TextRpcClient::hasErrors($result)) {
 *          ...
 *      }
 *
 *      # case 2
 *      $userInfo->getInfoByUid(100, function ($result, $errors) {
 *          if (!$errors) {
 *              ...
 *          }
 *      });
 *
 *      # 其中 RpcClient_ 是接口调用约定
 *      # RpcClient_User_Info::getInfoByUid 映射到
 *      # WebService 中的 \User\Service\Info 类和 getInfoByUid 方法
 *
 * 用户认证算法
 *
 *      # 客户端
 *      $packet = array(
 *          'data' => json_encode(
 *              array(
 *                  'version' => '2.0',
 *                  'user' => $this->rpcUser,
 *                  'password' => md5($this->rpcUser . ':' . $this->rpcSecret),
 *                  'timestamp' => microtime(true)); # 时间戳用于生成不同的签名, 以区分每一个独立请求
 *                  'class' => $this->rpcClass,
 *                  'method' => $method,
 *                  'params' => $arguments,
 *              )
 *          ),
 *      );
 *      $packet['signature'] = $this->encrypt($packet['data'], $secret);
 *
 *      # 服务器端
 *      # $this->encrypt($rawJsonData, $secret) === $packet['signature']
 *
 * 获取网络数据
 *
 *      TextRpcClient::on('send', function ($data) { });
 *      TextRpcClient::on('recv', function ($data) { });
 */

/**
 * 客户端协议实现.
 */
class TextRpcClient
{

    private $connection;
    protected $rpcClass;
    private $rpcUri;
    private $rpcUser;
    private $rpcSecret;
    private $executionTimeStart;

    private static $events = array();

    /**
     * 设置或读取配置信息.
     *
     * @param array $config 配置信息.
     *
     * @return array|void
     */
    public static function config(array $config = array())
    {
        static $_config = array();
        if (empty($config)) {
            return $_config;
        }
        $_config = $config;
    }

    /**
     * 获取RPC对象实例.
     *
     * @param array $config 配置信息, 或配置节点.
     *
     * @return TextRpcClient
     */
    public static function instance($config = array())
    {
        $className = get_called_class();

        static $instances = array();
        $key = $className . '-';
        if (empty($config)) {
            $key .= 'whatever';
        } else {
            $key .= md5(serialize($config));
        }
        if (empty($instances[$key])) {
            $instances[$key] = new $className($config);
            $instances[$key]->rpcClass = $className;
        }
        return $instances[$key];
    }

    /**
     * 检查返回结果是否包含错误信息.
     *
     * @param mixed $ctx 调用RPC接口时返回的数据.
     *
     * @return boolean
     */
    public static function hasErrors(&$ctx)
    {
        if (is_array($ctx)) {
            if (isset($ctx['error'])) {
                $ctx = $ctx['error'];
                return true;
            }
            if (isset($ctx['errors'])) {
                $ctx = $ctx['errors'];
                return true;
            }
        }
        return false;
    }

    /**
     * 注册各种事件回调函数.
     *
     * @param string   $eventName     事件名称, 如: read, recv.
     * @param function $eventCallback 回调函数.
     *
     * @return void
     */
    public static function on($eventName, $eventCallback)
    {
        if (empty(self::$events[$eventName])) {
            self::$events[$eventName] = array();
        }
        array_push(self::$events[$eventName], $eventCallback);
    }

    /**
     * 调用事件回调函数.
     *
     * @param $eventName 事件名称.
     *
     * @return void.
     */
    private static function emit($eventName)
    {
        if (!empty(self::$events[$eventName])) {
            $args = array_slice(func_get_args(), 1);
            foreach (self::$events[$eventName] as $callback) {
                @call_user_func_array($callback, $args);
            }
        }
    }

    /**
     * 构造函数.
     *
     * @param array $config 配置信息, 或配置节点.
     *
     * @throws Exception 抛出开发错误信息.
     */
    private function __construct(array $config = array())
    {
        if (empty($config)) {
            $config = self::config();
        } else {
            self::config($config);
        }

        if (empty($config)) {
            throw new Exception('TextRpcClient: Missing configurations');
        }

        $className = get_called_class();
        if (preg_match('/^[A-Za-z0-9]+_([A-Za-z0-9]+)/', $className, $matches)) {
            $module = $matches[1];
            if (empty($config[$module])) {
                throw new Exception(sprintf('TextRpcClient: Missing configuration for `%s`', $module));
            } else {
                $this->init($config[$module]);
            }
        } else {
            throw new Exception(sprintf('TextRpcClient: Invalid class name `%s`', $className));
        }

        // $this->openConnection();
    }

    /**
     * 析构函数.
     */
    public function __destruct()
    {
        // $this->closeConnection();
    }

    /**
     * 读取初始化配置信息.
     *
     * @param array $config 配置.
     *
     * @return void
     */
    public function init(array $config)
    {
        $this->rpcUri = $config['uri'];
        $this->rpcUser = $config['user'];
        $this->rpcSecret = $config['secret'];
        $this->rpcCompressor = isset($config['compressor']) ? strtoupper($config['compressor']) : null;
    }

    /**
     * 创建网络链接.
     *
     * @throws Exception 抛出链接错误信息.
     *
     * @return void
     */
    private function openConnection()
    {
        $this->connection = @stream_socket_client($this->rpcUri, $errno, $errstr);
        if (!$this->connection) {
            throw new Exception(sprintf('TextRpcClient: %s, %s', $this->rpcUri, $errstr));
        }
        @stream_set_timeout($this->connection, 60);
    }

    /**
     * 关闭网络链接.
     *
     * @return void
     */
    private function closeConnection()
    {
        @fclose($this->connection);
    }

    /**
     * 请求数据签名.
     *
     * @param string $data   待签名的数据.
     * @param string $secret 私钥.
     *
     * @return string
     */
    private function encrypt($data, $secret)
    {
        return md5($data . '&' . $secret);
    }

    /**
     * 调用 RPC 方法.
     *
     * @param string $method    PRC 方法名称.
     * @param mixed  $arguments 方法参数.
     *
     * @throws Exception 抛出开发用的错误提示信息.
     *
     * @return mixed
     */
    public function __call($method, $arguments)
    {
        $sign = '' . $this->rpcSecret;

        $fn = null;
        // if (!empty($arguments) && is_callable($arguments[count($arguments) - 1])) {
        //    $fn = array_pop($arguments);
        // }

        $packet = array(
            'data' => json_encode(
                array(
                    'version' => '1.0',
                    'user' => $this->rpcUser,
                    'password' => md5($this->rpcUser . ':' . $this->rpcSecret),
                    'timestamp' => microtime(true),
                    'class' => $this->rpcClass,
                    'method' => $method,
                    'params' => $arguments,
                )
            ),
        );

        $config = self::config();
        $packet['signature'] = $this->encrypt($packet['data'], $config['rpc_secret_key']);

        $ctx = $this->remoteCall($packet);

        if (isset($ctx['exception']) && is_array($ctx['exception'])) {
            throw new Exception('RPC Exception: ' . var_export($ctx['exception'], true));
        }

        if ($fn === null)
            return $ctx;

        if ($this->hasErrors($ctx)) {
            $fn(null, $ctx);
        } else {
            $fn($ctx, null);
        }
    }

    /**
     * 发起 RPC 调用协议.
     *
     * @param array $data RPC 数据.
     *
     * @throws Exception 抛出开发用的错误提示信息.
     *
     * @return mixed
     */
    private function remoteCall(array $data)
    {
        $this->executionTimeStart = microtime(true);

        $this->openConnection();

        // 用 JSON 序列化请求数据
        if (!$data = json_encode($data)) {
            throw new Exception('TextRpcClient: Cannot serilize $data with json_encode');
        }

        $fp = $this->connection;

        // 压缩数据
        $command = 'RPC';
        if ($this->rpcCompressor === 'GZ') {
            $data = @gzcompress($data);
            $command .= ':GZ';
        } elseif ($this->rpcCompressor) {
            throw new Exception(sprintf('TextRpcClient: Unsupported compress method `%s`', $this->rpcCompressor));
        }

        // 发送 RPC 文本请求协议
        $buffer = sprintf("%d\n%s\n%d\n%s\n", strlen($command), $command, strlen($data), $data);
        if (!@fwrite($fp, $buffer)) {
            throw new Exception(sprintf('TextRpcClient: Network %s disconnected', $this->rpcUri));
        }
        self::emit('send', $data);

        // 读取 RPC 返回数据的长度信息
        if (!$length = @fgets($fp)) {
            throw new Exception(
                sprintf(
                    'TextRpcClient: Network %s may timed out(%.3fs), or there are fatal errors on the RPC server',
                    $this->rpcUri,
                    $this->executionTime()
                )
            );
        }
        $length = trim($length);
        if (!preg_match('/^\d+$/', $length)) {
            throw new Exception(sprintf('TextRpcClient: Got wrong protocol codes: %s', bin2hex($length)));
        }
        $length = 1 + $length; // 1 means \n

        // 读取 RPC 返回的具体数据
        $ctx = '';
        while (strlen($ctx) < $length) {
            $ctx .= fgets($fp);
        }
        self::emit('recv', $ctx);
        $ctx = trim($ctx);

        $this->closeConnection();

        // 反序列化 JSON 数据并返回
        if ($ctx !== '') {
            if ($this->rpcCompressor === 'GZ') {
                $ctx = @gzuncompress($ctx);
            }
            $ctx = json_decode($ctx, true);
            return $ctx;
        }
    }

    /**
     * 计算 RPC 请求时间.
     *
     * @return float
     */
    private function executionTime()
    {
        return microtime(true) - $this->executionTimeStart;
    }

}

spl_autoload_unregister(array('YiiBase','autoload'));
spl_autoload_register(
    function ($className) {
        if (strpos($className, 'RpcClient_') !== 0)
            return false;

        eval(sprintf('class %s extends TextRpcClient {}', $className));
    }
);
spl_autoload_register(array('YiiBase','autoload'));

if (false) {
    $config = array(
        'rpc_secret_key' => '769af463a39f077a0340a189e9c1ec28',
        'User' => array(
            'uri' => 'tcp://127.0.0.1:2201',
            'user' => 'Optool',
            'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
            //'compressor' => 'GZ',
        ),
        'Item' => array(
            'uri' => 'tcp://127.0.0.1:2201',
            'user' => 'Optool',
            'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
        ),
        'Order' => array(
            'uri' => 'tcp://127.0.0.1:2201',
            'user' => 'Optool',
            'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
        ),
    );

    TextRpcClient::on('send', function($data) {
        echo 'Send => ', $data, PHP_EOL;
    });
    TextRpcClient::on('recv', function($data) {
        echo 'Recv <= ', $data, PHP_EOL;
    });

    TextRpcClient::config($config);
    //$test = RpcClient_Item_Iwc::instance();
    //var_export($test->getInventoryByWarehouses(array(100223,100002,100003,100006), array('BJ08','GZ07','SH05')));

    $test = RpcClient_User_Address::instance($config);
    //var_dump($test->getListByUid(5100));
    $test->getListByUid(5100, function () {
        var_dump(func_get_args());
    });
}
