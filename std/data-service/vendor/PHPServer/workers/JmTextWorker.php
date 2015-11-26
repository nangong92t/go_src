<?php 

use MNLogger\MNLogger;

class JmTextWorker extends RpcWorker
{
    public static $rpcMNLogger = null;
    public static $appName = 'JmTextWorker';

    public function onServe()
    {
        require_once SERVER_BASE . 'thirdparty/MNLogger/MNLogger.php';
        $logdir = PHPServerConfig::get('monitor_log_path') ? PHPServerConfig::get('monitor_log_path') : '/home/logs/monitor';
        $config = array(
             'on' => true,
             'app' => 'php-rpc-server',
             'logdir' => $logdir,
        );
        try{
            self::$rpcMNLogger = @thirdparty\MNLogger\MNLogger::instance($config);
        }catch(Exception $e){}
        
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
    
    protected function process($data)
    {
        JmTextStatistic::tick();
        $class_name = '\\Handler\\'.$data['class'];
        $_SERVER['REMOTE_ADDR'] = $this->getRemoteIp();
        try
        {
            // 请求开始时执行的函数，on_request_start一般在bootstrap初始化
            if(function_exists('on_phpserver_request_start'))
            {
                \on_phpserver_request_start();
            }
            if(class_exists($class_name))
            {
                $call_back = array(new $class_name, $data['method']);
                if(is_callable($call_back))
                {
                    $ctx = call_user_func_array($call_back, $data['params']);
                }
                else
                {
                    throw new Exception("method $class_name::{$data['method']} not exist");
                }
            }
            else
            {
                throw new Exception("class $class_name not exist");
            }
        }
        catch (Exception $ex)
        {
            $ctx = array(
                'exception' => array(
                    'class' => get_class($ex),
                    'message' => $ex->getMessage(),
                    'code' => $ex->getCode(),
                    'file' => $ex->getFile(),
                    'line' => $ex->getLine(),
                    'traceAsString' => $ex->getTraceAsString(),
                )
            );
        }
        
        // 请求结束时执行的函数，on_request_start一般在bootstrap中初始化
        if(function_exists('on_phpserver_request_finish'))
        {
            // 这里一般是关闭数据库链接等操作
            \on_phpserver_request_finish();
        }
        
        JmTextStatistic::report($data, $ctx, $this->getRemoteIp());
        
        $this->send($ctx);
    }

    public function handleFatalErrors() {
        if ($errors = error_get_last()) {
            $this->send(array(
                'exception' => array(
                    'class' => 'WebServiceFatalException',
                    'message' => $errors['message'],
                    'file' => $errors['file'],
                    'line' => $errors['line'],
                ),
            ));
        }
    }
}


/**
 * 针对JumeiWorker对统计模块的一层封装
 * @author liangl
 */
class JmTextStatistic
{
    protected static $timeStart = 0;
    
    public static function tick()
    {
        self::$timeStart = StatisticClient::tick();
    }

    public static function report($data, $ctx, $source_ip)
    {
        $module = $data['class'];
        $interface = $data['method'];
        $code = 0;
        $msg = '';
        $success = true;
        if(is_array($ctx) && isset($ctx['exception']))
        {
            $success = false;
            $code = isset($ctx['exception']['code']) ? $ctx['exception']['code'] : 40404;
            $msg = isset($ctx['exception']['class']) ? $ctx['exception']['class'] . "::" : '';
            $msg .= isset($ctx['exception']['message']) ? $ctx['exception']['message'] : '';
            $msg .= "\n" . $ctx['exception']['traceAsString'];
            $msg .= "\nREQUEST_DATA:[" . json_encode($data) . "]\n";
        }
        
        // 格式key[模块名,接口名] val[数量,耗时,是否成功,错误码]
        try {
            if(JmTextWorker::$rpcMNLogger)
            {
                JmTextWorker::$rpcMNLogger->log('RPC-Server,'.JmTextWorker::$appName.','.($success ? 'Success' : 'Failed'), microtime(true)-self::$timeStart);
            }
        }catch(Exception $e){}
        
        StatisticClient::report($module, $interface, $code, $msg, $success, $source_ip);
    }
}

