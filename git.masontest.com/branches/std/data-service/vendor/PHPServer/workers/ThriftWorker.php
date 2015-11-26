<?php
/**
 * thriftWorker
 * @author libingw <libingw@jumei.com>
 * @author liangl  <liangl3@jumei.com>
 * 
 */
require_once SERVER_BASE . 'thirdparty/Thrift/Thrift/Context.php';
require_once SERVER_BASE . 'thirdparty/Thrift/Thrift/ContextSerialize.php';
require_once SERVER_BASE . 'thirdparty/Thrift/Thrift/ContextReader.php';
require_once SERVER_BASE . 'thirdparty/Thrift/Thrift/ClassLoader/ThriftClassLoader.php';

use Thrift\ClassLoader\ThriftClassLoader;

$loader = new ThriftClassLoader();
$loader->registerNamespace('Thrift', SERVER_BASE . 'thirdparty/Thrift');
$loader->register();

define('IN_THRIFT_WORKER', true);

class ThriftWorker extends PHPServerWorker
{
    /**
     * 统一监控日志类
     * @var object
     */
    public static $rpcMNLogger = null;
    
    /**
     * 存放thrift生成文件的目录
     * @var string
     */
    protected $providerDir = null;
    
    /**
     * 存放对thrift生成类的实现的目录
     * @var string
     */
    protected $handlerDir = null;
    
    /**
     * thrift生成类的命名空间
     * @var string
     */
    protected $providerNamespace = 'Provider';
    
    /**
     * thrift生成类实现的命名空间
     * @var string
     */
    protected $handlerNamespace = 'Provider';
    
    /**
     * 服务名
     * @var string
     */
    public static $appName = 'ThriftWorker';
    
    /**
     * 进程启动时的一些初始化
     * @see PHPServerWorker::onServe()
     */
    public function onServe()
    {
        // 初始化thrift生成文件存放目录
        $provider_dir = PHPServerConfig::get('workers.'.$this->serviceName.'.provider');
        if($provider_dir)
        {
            if($this->providerDir = realpath($provider_dir))
            {
                if($path_array = explode('/', $this->providerDir))
                {
                    $this->providerNamespace = $path_array[count($path_array)-1];
                }
            }
            else
            {
                $this->providerDir = $provider_dir;
                $this->notice('provider_dir '.$provider_dir. ' not exsits');
            }
        }
        
        // 初始化thrift生成类业务实现存放目录
        $handler_dir = PHPServerConfig::get('workers.'.$this->serviceName.'.handler');
        if($handler_dir)
        {
            if($this->handlerDir = realpath($handler_dir))
            {
                if($path_array = explode('/', $this->handlerDir))
                {
                    $this->handlerNamespace = $path_array[count($path_array)-1];
                }
            }
            else
            {
                $this->handlerDir = $handler_dir;
                $this->notice('handler_dir' . $handler_dir. ' not exsits');
            }
        }
        else
        {
            $this->handlerDir = $provider_dir;
        }
        
        // 统一日志类初始化
        require_once SERVER_BASE . 'thirdparty/MNLogger/MNLogger.php';
        $logdir = PHPServerConfig::get('monitor_log_path') ? PHPServerConfig::get('monitor_log_path') : '/home/logs/monitor';
        $config = array(
                        'on' => true,
                        'app' => 'php-rpc-server',
                        'logdir' => $logdir,
        );
        try{
            self::$rpcMNLogger = @thirdparty\MNLogger\MNLogger::instance($config);
        }
        catch(Exception $e)
        {
            
        }
        
        // 初始化统计上报地址
        $report_address = PHPServerConfig::get('workers.'.$this->serviceName.'.report_address');
        if($report_address)
        {
            StatisticClient::config(array('report_address'=>$report_address));
        }
        // 没有配置则使用本地StatisticWorker中的配置
        else
        {
            if($config = PHPServerConfig::get('workers.StatisticWorker'))
            {
                if(!isset($config['ip']))
                {
                    $config['ip'] = '127.0.0.1';
                }
                StatisticClient::config(array('report_address'=>'udp://'.$config['ip'].':'.$config['port']));
            }
        }
        
        // 业务引导程序bootstrap初始化（没有则忽略）
        $bootstrap = PHPServerConfig::get('workers.'.$this->serviceName.'.bootstrap');
        if(is_file($bootstrap))
        {
            require_once $bootstrap;
        }
        
        // 服务名
        self::$appName = $this->serviceName;
    }
    
    /**
     * 处理thrift包，判断包是否接收完整
     * 固定使用TFramedTransport，前四个字节是包体长度信息
     * @see PHPServerWorker::dealInput()
     */
    public function dealInput($recv_str) {
        $val = unpack('N', $recv_str);
        $length = $val[1] + 4;
        if ($length <= Thrift\Factory\TStringFuncFactory::create()->strlen($recv_str)) {
            return 0;
        }
        return 1;
    }

    /**
     * 业务处理(non-PHPdoc)
     * @see PHPServerWorker::dealProcess()
     */
    public function dealProcess($recv_str) {
        // 拷贝一份数据包，用来记录日志
        $recv_str_copy = $recv_str;
        // 统计监控记录请求开始时间点，后面用来统计请求耗时
        StatisticHelper::tick();
        // 清除上下文信息
        \Thrift\Context::clear();
        // 服务名
        $serviceName = $this->serviceName;
        // 本地调用方法名
        $method_name = 'none';
        // 来源ip
        $source_ip = $this->getRemoteIp();
        // 尝试读取上下文信息
        try{
            // 去掉TFrameTransport头
            $body_str = substr($recv_str, 4);
            // 读上下文,并且把上下文数据从数据包中去掉
            \Thrift\ContextReader::read($body_str);
            // 再组合成TFrameTransport报文
            $recv_str = pack('N', strlen($body_str)).$body_str;
            // 如果是心跳包
            if (\Thrift\Context::get('isHeartBeat')=='true'){
                $thriftsocket = new \Thrift\Transport\TBufferSocket();
                $thriftsocket->setHandle($this->connections[$this->currentDealFd]);
                $thriftsocket->setBuffer($recv_str);
                $framedTrans = new \Thrift\Transport\TFramedTransport($thriftsocket);
                $protocol = new Thrift\Protocol\TBinaryProtocol($framedTrans, false, false);
                $protocol->writeMessageBegin('#$%Heartbeat', 2, 0);
                $protocol->writeMessageEnd();
                $protocol->getTransport()->flush();
                return;
            }
        }
        catch(Exception $e)
        {
            // 将异常信息发给客户端
            $this->writeExceptionToClient($method_name, $e, self::getProtocol(\Thrift\Context::get('protocol')));
            // 统计上报
            StatisticHelper::report($serviceName, $method_name, $source_ip, $e, $recv_str_copy);
            return;
        }
        
        // 客户端有传递超时参数
        if(($timeout = \Thrift\Context::get("timeout")) && $timeout >= 1)
        {
            pcntl_alarm($timeout);
        }
        
        // 客户端有传递服务名
        if(\Thrift\Context::get('serverName'))
        {
            $serviceName = \Thrift\Context::get('serverName');
        }
        
        // 尝试处理业务逻辑
        try {
            // 服务名为空
            if (!$serviceName){
                throw new \Exception('Context[serverName] empty', 400);
            }
            
            // 如果handler命名空间为provider
            if($this->handlerNamespace == 'Provider')
            {
                $handlerClass = $this->handlerNamespace.'\\'.$serviceName.'\\' . $serviceName . 'Handler';
            }
            else
            {
                $handlerClass = $this->handlerNamespace.'\\' . $serviceName;
            }
            
            // processor
            $processorClass = $this->providerNamespace . '\\' . $serviceName . '\\' . $serviceName . 'Processor';
            
            // 文件不存在尝试从磁盘上读取
            if(!class_exists($handlerClass, false))
            {
                clearstatcache();
                if(!class_exists($processorClass, false))
                {
                    require_once $this->providerDir.'/'.$serviceName.'/Types.php';
                    require_once $this->providerDir.'/'.$serviceName.'/'.$serviceName.'.php';
                }
                
                $handler_file = $this->handlerNamespace == 'Provider' ? $this->handlerDir.'/'.$serviceName.'/'.$serviceName.'Handler.php' : $this->handlerDir.'/'.$serviceName.'.php';
                if(is_file($handler_file))
                {
                    require_once $handler_file;
                }
                
                if(!class_exists($handlerClass))
                {
                    throw new \Exception('Class ' . $handlerClass . ' not found', 404);
                }
            }
            
            // 运行thrift
            $handler = new $handlerClass();
            $processor = new $processorClass($handler);
            $pname = \Thrift\Context::get('protocol') ? \Thrift\Context::get('protocol') : 'binary';
            $protocolName = self::getProtocol($pname);
            $thriftsocket = new \Thrift\Transport\TBufferSocket();
            $thriftsocket->setHandle($this->connections[$this->currentDealFd]);
            $thriftsocket->setBuffer($recv_str);
            $framedTrans = new \Thrift\Transport\TFramedTransport($thriftsocket, true, true);
            $protocol = new $protocolName($framedTrans, false, false);
            $protocol->setTransport($framedTrans);
            // 请求开始时执行的函数，on_request_start一般在bootstrap初始化
            if(function_exists('on_phpserver_request_start'))
            {
                \on_phpserver_request_start();
            }
            $processor->process($protocol, $protocol);
            // 请求结束时执行的函数，on_request_start一般在bootstrap中初始化
            if(function_exists('on_phpserver_request_finish'))
            {
                // 这里一般是关闭数据库链接等操作
                \on_phpserver_request_finish();
            }
            $method_name = $protocol->fname;
        }
        catch (Exception $e)
        {
            // 异常信息返回给客户端
            $method_name = !empty($protocol->fname) ? $protocol->fname : 'none';
            $this->writeExceptionToClient($method_name, $e, !empty($protocolName) ? $protocolName : 'Thrift\Protocol\TBinaryProtocol');
            StatisticHelper::report($serviceName, $method_name, $source_ip, $e, $recv_str_copy);
            return;
        }
        // 统计上报
        StatisticHelper::report($serviceName, $method_name, $source_ip);
    }

    /**
     * 获取协议全名
     * @param string $key
     * @return string
     */
    private static function getProtocol($key=null){
        $protocolArr = array(
          'binary'=>'Thrift\Protocol\TBinaryProtocol',
          'compact'=>'Thrift\Protocol\TCompactProtocol',
          'json'   => 'Thrift\Protocol\TJSONProtocol',
        );
        return isset($protocolArr[$key]) ? $protocolArr[$key] : $protocolArr['binary'];
    }
    
    /**
     * 将异常写会客户端
     * @param string $name
     * @param Exception $e
     * @param string $protocol
     */
    private function writeExceptionToClient($name, $e, $protocol = 'Thrift\Protocol\TBinaryProtocol')
    {
        try {
            $ex = new \Thrift\Exception\TApplicationException($e);
            $thriftsocket = new \Thrift\Transport\TBufferSocket();
            $thriftsocket->setHandle($this->connections[$this->currentDealFd]);
            $framedTrans = new \Thrift\Transport\TFramedTransport($thriftsocket, true, true);
            $protocol = new $protocol($framedTrans, false, false);
            $protocol->writeMessageBegin($name, \Thrift\Type\TMessageType::EXCEPTION, 0);
            $ex->write($protocol);
            $protocol->writeMessageEnd();
            $protocol->getTransport()->flush();
        }
        catch(Exception $e)
        {
            
        }
    }
    
}


/**
 * 针对JumeiWorker对统计模块的一层封装
 * @author liangl
 */
class StatisticHelper
{
    protected static $timeStart = 0;

    public static function tick()
    {
        self::$timeStart = StatisticClient::tick();
    }

    public static function report($serviceName, $method, $source_ip, $exception = null, $request_data = '')
    {
        $success = empty($exception);
        $code = 0;
        $msg = '';
        $success = true;
        if($exception)
        {
            $success = false;
            $code = $exception->getCode();
            $msg = $exception;
            $msg .= "\nREQUEST_DATA:[" . bin2hex($request_data) . "]\n";
        }

        // 格式key[模块名,接口名] val[数量,耗时,是否成功,错误码]
        try {
            if(ThriftWorker::$rpcMNLogger)
                ThriftWorker::$rpcMNLogger->log("RPC-Server,".ThriftWorker::$appName.",".($success ? 'Success' : 'Failed'), microtime(true)-self::$timeStart);
        }catch(Exception $e){
        }

        StatisticClient::report($serviceName, $method, $code, $msg, $success, $source_ip);
    }
}
