<?php
namespace Log;

/**
 * unified log handlers
 */
class Handler {
    /**
     *
     * @var log instances
     */
    protected static $loggers = array ();

    /**
     * available loggers of the current version
     *
     * @var array
     */
    protected static $availableLoggers = array (
            'php',
            'file',
            'mongo',
            'jsonfile'
    );

    /**
     * Log record will be prefixed with time.
     * @var boolean
     */
    protected $prefixedWithTime=true;

    protected $useRawMessage = false;

    /**
     * config of the current logger
     *
     * @var array
     */
    protected static $configs;

    /**
     * configuration name of the logger , which defined in \Core\Config\Log
     * @var string
     */
    protected $cfgName;

    /**
     * maps of the E_USER* error types to their comment strings in php error logs.
     * @var array
     */
    protected static $phpUserErrorTypeToStringMap = array(
              E_USER_NOTICE => 'Notice',
              E_USER_WARNING => 'Warning',
              E_USER_DEPRECATED => 'Deprecated',
              E_USER_ERROR => 'Fatal error'
    );

    protected function __construct($cfgName) {
        if(!static::$configs)
        {
            static::$configs = (array) new \Config\Log;
        }
        if (!isset(static::$configs[$cfgName]))
            return false;
        $cfg = static::$configs[$cfgName];
        if (!in_array($cfg ['logger'], self::$availableLoggers ))
        {
            return false;
        }
        $this->loggerCfg = $cfg;
        $this->cfgName = $cfgName;
    }


    /**
     * Set or get configs for the lib.
     *
     * @param string $config
     * @return boolean
     */
    public static function config($config=null)
    {
        if(is_null($config))
        {
            return static::$configs;
        }
        static::$configs = $config;
        return true;
    }

    /**
     *
     * @param string $cfg
     * @return \Core\Lib\Log
     */
    public static function instance($cfgName = 'default') {
        if (! isset ( self::$loggers [$cfgName] )) {
            self::$loggers [$cfgName] = new self ( $cfgName );
        }
        return self::$loggers [$cfgName];
    }
    public function log($msg, $options=array()) {
        if (! $this->loggerCfg)
            return false;
        switch ($this->loggerCfg ['logger']) {
            case 'php' :
                $this->phpLogger($msg, $options);
                continue;
            case 'sys' :
                // @todo implementation
                continue;
            case 'file' :
                $this->fileLogger ( $msg );
                continue;
            case 'mongo' :
                call_user_func_array(array($this, 'mongoLogger'), func_get_args());
                continue;
            case 'jsonfile' :
                $this->useRawMessage = true;
                $this->prefixedWithTime = false;
                $this->jsonfileLogger($msg);
                continue;
        }
    }
    protected function fileLogger($msg) {
        $logDir = static::$configs['FILE_LOG_ROOT'].DIRECTORY_SEPARATOR.$this->cfgName . DIRECTORY_SEPARATOR;
        $timefile = $this->cfgName.'.log';
        if (! is_dir ( $logDir ) && ! @mkdir ( $logDir, 0777, true ) && !is_writable ( $logDir )) {
            return false;
        }
        if($this->useRawMessage)
        {
            $msg = $this->formatMessage($msg);
        }
        file_put_contents ( $logDir . $timefile, $msg, FILE_APPEND );
    }

    /**
     * @param array $msg
     */
    protected function mongoLogger($msg)
    {
        if(!isset($this->loggerCfg['dbConfigName']))
        {
            throw new Exception('Invalid log config of mongo handler, $dbConfigName missing!');
        }
        if(isset($msg['time']))
        {
            $msg['time'] = new \MongoDate(strtotime($msg['time']));
        }
        $db = Mongo::instance($this->loggerCfg['dbConfigName'])->selectDB()->selectCollection($this->loggerCfg['dbCollection']);
        if(is_array($msg))
        {
            if(is_array(current($msg)))
            {
                foreach ($msg as $record)
                {
                    $db->save($record);
                }
            }
            else
            {
                $db->save($msg);
            }
        }
        else
        {
            throw new Exception('$msg for mongo handler should be array!');
        }
    }

    /**
     * send log message to the php system log that defined as error_log in php.ini. messages are formatted as the standard php errors , default error type is E_NOTICE.
     *
     * @param string $msg message body
     * @param array $options availabe options are  'type': error_type which introduce within E_USER_* default is E_USER_NOTICE;'trace_depth': trace depth, default is 1, which means just only record the trace line where you called the "log" method.
     */
    protected function phpLogger($msg, $options)
    {
        if(!is_array($options))
        {
            $options = array();
        }
        $type = isset($options['type']) && in_array($options['type'], array_keys($this::$phpUserErrorTypeToStringMap))? $options['type'] : E_USER_NOTICE;
        $traceDepth = isset($options['trace_depth']) ? (int) $options['trace_depth'] : 1;
        $msg = 'PHP '.$this::$phpUserErrorTypeToStringMap[$type].': '.$msg.' ' ;
        $traceStr = '';
        $traceDepth += 1 ;//count the call of self::Log in trace.
        $traceStr = ' Stack Trace: ';
        if(version_compare(PHP_VERSION, '5.4.0', '>='))
        {
            $trace = debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS);
        }
        else
        {
            $trace = debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS);
        }
        unset($trace[0]);//remove the trace of self::Log
        if (sizeof($trace) > 0)
        {
            foreach ($trace as $k => $v)
            {
                extract($v);
                $traceStr .= $k.'. in '.$file.':'.$line . '\n';
            }
            error_log($msg.$traceStr);
        }
    }

    /**
     * Log message object as json to file.
     *
     * @param array $data
     */
    protected function jsonfileLogger($data)
    {
        $data['log_create_time'] = time();
        $msg = json_encode($data, 256);
        return $this->fileLogger($msg);
    }

    protected function formatMessage($msg)
    {
        $msg = str_replace(array("\n","\r\n"), '\n',$msg);
        if($this->prefixedWithTime)
        {
            $msg = date('Ymd H:i:s ').$msg;
        }
        $msg .="\n";
        return $msg;
    }

    /**
     * 用于外部判断是否记录日志.
     *
     * @param string $cfgName 配置名.
     *
     * return boolean
     */
    public static function isLog($cfgName)
    {
        if (!isset(self::$config[$cfgName]))
            return false;
        $cfg = $logConfig->$cfgName;
        if (! in_array ( $cfg ['logger'], self::$availableLoggers )) {
            return false;
        }
        return true;
    }
}

class Exception extends \Exception{

}
