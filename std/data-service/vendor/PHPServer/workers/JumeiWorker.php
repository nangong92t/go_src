<?php 

use MNLogger\MNLogger;

class JumeiWorker extends RpcWorker
{
    public static $rpcMNLogger = null;
    public static $appName = 'JumeiWorker';

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
        
        $app_name = PHPServerConfig::get('workers.'.$this->serviceName.'.framework.app_name');
        if($app_name)
        {
            self::$appName = $app_name;
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
    }
    
    protected function process($data)
    {
        if (!class_exists('Core\Lib\RpcServer')) {
            $frameworkBootstrap = PHPServerConfig::get('workers.'.$this->serviceName.'.framework.path') .
                '/Serverroot/Autoload.php';
            require_once $frameworkBootstrap;
        }

        StatisticHelper::tick();
        
        $rpcServer = new Core\Lib\RpcServer;
        $ctx = $rpcServer->run($data);
        
        StatisticHelper::report($data, $ctx, $this->getRemoteIp());
        
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
class StatisticHelper
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
            if(JumeiWorker::$rpcMNLogger)
                JumeiWorker::$rpcMNLogger->log('RPC-Server,'.JumeiWorker::$appName.','.($success ? 'Success' : 'Failed'), microtime(true)-self::$timeStart);
        }catch(Exception $e){}
        
        
        StatisticClient::report($module, $interface, $code, $msg, $success, $source_ip);
    }
}

