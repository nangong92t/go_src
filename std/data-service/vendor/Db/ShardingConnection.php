<?php
/**
 * Class DbShardingConnection
 *
 * @author Haojie Huang<haojieh@jumei.com>
 */

namespace Db;

use \PDO as PDO;

/**
 * DbShardingConnection.
 */

class ShardingConnection extends Connection
{

    /**
     * instance of the DbConnection
     */
    protected static $instance;
    protected static $configs;
    protected $currentRule = null;
    protected static $writeConnections = array();
    protected static $readConnections  = array();
    protected static $writeAsyncConnections = array();
    protected static $readAsyncConnections  = array();
    protected static $asyncConnections = array();
    protected $isPdo = true;

    protected static $mysqliPdoOptions = array(
        MYSQLI_INIT_COMMAND        => \PDO::MYSQL_ATTR_INIT_COMMAND,
        MYSQLI_OPT_CONNECT_TIMEOUT => \PDO::ATTR_TIMEOUT,
    );

    /**
     * 构造方法.
     *
     * @param string  $host     Host.
     * @param string  $username 用户名.
     * @param string  $passwd   密码.
     * @param string  $db       数据库名.
     * @param integer $port     端口.
     * @param array   $options  需要配置的选项.
     *
     * @access protected
     */
    protected function __construct($host = null, $username = null, $passwd = null, $db = null, $port = null, array $options = array())
    {
        if ( ! self::$configs) {
            self::$configs = (array) new \Config\DbSharding;
        }
        if ( ! is_null($host)) {
            return $this->connect($host, $username, $passwd, $db, $port, $options);
        }
    }

    /**
     * 连接数据库方法.
     *
     * @param string $host     Host.
     * @param string $username 用户名.
     * @param string $passwd   密码.
     * @param string $db       数据库名.
     * @param string $port     端口号.
     * @param array  $options  需要配置的选项.
     *
     * @access public
     * @return object
     * @throws \Exception connecton error.
     */
    public function connect($host, $username = null, $passwd = null, $db = null, $port = null, array $options = array())
    {
        if (is_array($host)) {
            extract($host);
        }

        if ($this->connection) {
            return $this->connection;
        } else {
            $rule = $this->currentRule;
            $db = $rule->getDbName($db);
            if ($this->isPdo) {
                $dsn = "mysql:host={$host};port={$port};dbname=$db";
                try {
                    $this->connection = new \PDO($dsn, $username, $passwd, $this->genPdoOptions($options));
                } catch (\Exception $ex) {
                    if ($ex->getCode() === 2002) {
                        // Retry once on connection error.
                        $this->connection = new \PDO($dsn, $username, $passwd, $this->genPdoOptions($options));
                    } else {
                        throw $ex;
                    }
                }
            } else {
                $this->connection = new Mysqli($host, $username, $passwd, $db, $port, $options);
                if ($this->connection->connect_errno === 2002) {
                    // Retry once on connection error.
                    $this->connection = new Mysqli($host, $username, $passwd, $db, $port, $options);
                } elseif ($this->connection->connect_errno) {
                    throw new \Exception('Mysqli connection error:('.$this->connection->connect_errno.') '.$this->connection->connect_error);
                }
            }
        }
        return $this;
    }

    /**
     * 按照rule获取读库.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     * @throws \Exception   rule不是ShardingRule的实例.
     */
    public function read($rule = null)
    {
        if ($rule != null && ! $rule instanceof ShardingRule) {
            throw new \Exception("sharding rule must be an instance of ShardingRule.");
        }
        if ($rule == null && func_num_args() == 0 && $this->currentRule) {
            $rule = $this->currentRule;
        }

        $name  = $rule->getCfgName();
        if (isset(self::$configs['read'][$name] )) {
            $cfg = self::$configs['read'][$name];
            $db = $rule->getDbName($cfg['db']);

            if (empty ( self::$readConnections[$db]) && ! $this->addReadConnection($rule)) {
                throw new Exception('No available read connections. Please use addReadConnection to initialize  first', 42001);
            }
            return  self::$readConnections[$db];
        } else {
            throw new Exception('Read configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * 按照rule获取读库.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws \Exception  rule不是ShardingRule的实例.
     * @throws Exception 初始化失败.
     */
    public function write($rule = null)
    {
        if ($rule != null && ! $rule instanceof ShardingRule) {
            throw new \Exception("sharding rule must be an instance of ShardingRule.");
        }
        if ($rule == null && func_num_args() == 0 && $this->currentRule) {
            $rule = $this->currentRule;
        }
        $name  = $rule->getCfgName();
        if (isset(self::$configs['write'][$name])) {
            $cfg = self::$configs['write'][$name];
            $db = $rule->getDbName($cfg['db']);
            if (empty( self::$writeConnections[$db]) && ! $this->addWriteConnection($rule)) {
                throw new Exception('No available write connections. Please use addWriteConnection to initialize  first', 42001);
            }
            return self::$writeConnections[$db];
        } else {
            throw new Exception('Write configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * Initialize read connections.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     */
    public function addReadConnection($rule = null)
    {
        $name  = $rule->getCfgName();
        if (isset(self::$configs['read'][$name] )) {
            $cfg = self::$configs['read'][$name];
            $options = isset($cfg['options']) ? $cfg['options'] : array();
            $connection = new self();
            $connection->isPdo = true;
            $connection->currentRule = $rule;
            $connection = $connection->connect($cfg ['host'], $cfg ['username'], $cfg ['passwd'], $cfg['db'], $cfg['port'], $options);
            $connection->connectionCfg = $cfg;
            $db = $rule->getDbName($cfg['db']);
            self::$readConnections [$db] = $connection;
            return self::$readConnections [$db];
        } else {
            throw new Exception('Read configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * Initialize write connections.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     */
    public function addWriteConnection($rule = null)
    {
        $name  = $rule->getCfgName();
        if (isset(self::$configs['write'][$name] )) {
            $cfg = self::$configs['write'][$name];
            $options = isset($cfg['options']) ? $cfg['options'] : array();
            $connection = new self();
            $connection->isPdo = true;
            $connection->currentRule = $rule;
            $connection = $connection->connect($cfg ['host'], $cfg ['username'], $cfg ['passwd'], $cfg['db'], $cfg['port'], $options);
            $connection->connectionCfg = $cfg;
            $db = $rule->getDbName($cfg['db']);
            self::$writeConnections [$db] = $connection;
            return self::$writeConnections [$db];
        } else {
            throw new Exception('Write configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * 按照rule获取读库.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     * @throws \Exception   rule不是ShardingRule的实例.
     */
    public function readAsync($rule = null)
    {
        if ($rule != null && ! $rule instanceof ShardingRule) {
            throw new \Exception("sharding rule must be an instance of ShardingRule.");
        }
        if ($rule == null && func_num_args() == 0 && $this->currentRule) {
            $rule = $this->currentRule;
        }
        if (empty ( self::$readAsyncConnections[$rule->getAtomName()]) && ! $this->addReadAsyncConnection($rule)) {
            throw new Exception('No available read connections. Please use addReadConnection to initialize  first', 42001);
        }
        return  self::$readAsyncConnections[$rule->getAtomName()];
    }

    /**
     * 按照rule获取读库.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws \Exception  rule不是ShardingRule的实例.
     * @throws Exception 初始化失败.
     */
    public function writeAsync($rule = null)
    {
        if ($rule != null && ! $rule instanceof ShardingRule) {
            throw new \Exception("sharding rule must be an instance of ShardingRule.");
        }
        if ($rule == null && func_num_args() == 0 && $this->currentRule) {
            $rule = $this->currentRule;
        }
        if (empty( self::$writeAsyncConnections[$rule->getAtomName()]) && ! $this->addWriteAsyncConnection($rule)) {
            throw new Exception('No available write connections. Please use addWriteConnection to initialize  first', 42001);
        }
        return self::$writeAsyncConnections[$rule->getAtomName()];
    }

    /**
     * Initialize read connections.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     */
    public function addReadAsyncConnection($rule = null)
    {
        $name  = $rule->getCfgName();
        if (isset(self::$configs['read'] [$name] )) {
            $cfg = self::$configs['read'] [$name];
            $options = isset($cfg['options']) ? $cfg['options'] : array();
            $connection = new self();
            $connection->isPdo = false;
            $connection->currentRule = $rule;
            $connection = $connection->connect($cfg ['host'], $cfg ['username'], $cfg ['passwd'], $cfg['db'], $cfg['port'], $options);
            $connection->connectionCfg = $cfg;
            self::$readAsyncConnections [$rule->getAtomName()] = $connection;
            return self::$readAsyncConnections [$rule->getAtomName()];
        } else {
            throw new Exception('Read configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * Initialize write connections.
     *
     * @param object $rule Sharding Rule.
     *
     * @access public
     * @return object
     * @throws Exception 初始化失败.
     */
    public function addWriteAsyncConnection($rule = null)
    {
        $name  = $rule->getCfgName();
        if (isset(self::$configs['write'][$name] )) {
            $cfg = self::$configs['write'][$name];
            $options = isset($cfg['options']) ? $cfg['options'] : array();
            $connection = new self();
            $connection->isPdo = false;
            $connection->currentRule = $rule;
            $connection = $connection->connect($cfg ['host'], $cfg ['username'], $cfg ['passwd'], $cfg['db'], $cfg['port'], $options);
            $connection->connectionCfg = $cfg;
            self::$writeAsyncConnections [$rule->getAtomName($name)] = $connection;
            return self::$writeAsyncConnections [$rule->getAtomName($name)];
        } else {
            throw new Exception('Write configuration of "' . $name . '" is not found.', 42003);
        }
    }

    /**
     * 设置Sharding规则.
     *
     * @param object $rule 规则对象.
     *
     * @return \Core\Lib\DbShardingConnection
     * @throws \Exception  rule不是ShardingRule的实例.
     */
    public function setRule($rule)
    {
        if ( ! $rule instanceof ShardingRule) {
            throw new \Exception("sharding rule must be an instance of ShardingRule.");
        }
        $this->currentRule = $rule;
        return $this;
    }

    /**
     * 设置异步获取数据.
     *
     * @return \Core\Lib\DbShardingConnection
     */
    public function async()
    {
        $this->async = true;
        $name = $this->currentRule->getTableName($this->currentRule->getCfgName());
        static::$asyncConnections[$name] = static::$writeConnections[$name];
        return $this;
    }

    /**
     * 查询.
     *
     * @param string  $sql   查询sql.
     * @param boolean $async 是否异步.
     *
     * @return PDOStatement
     * @throws Exception No WHERE condition in SQL(UPDATE, DELETE) to be executed! please make sure it's safe.
     */
    public function query($sql = null, $async = false)
    {
        static $retryCount = 0;

        $withCache = false;
        $this->withCache = true;

        if (empty($sql)) {
            $this->lastSql = $this->getSelectSql();  // 不需要trim，拼接函数保证以SELECT开头
        } else {
            $this->lastSql = trim($this->buildSql($sql));
        }
        $sqlCmd = strtoupper(substr($this->lastSql, 0, 6));
        if (in_array($sqlCmd, array('UPDATE', 'DELETE')) && stripos($this->lastSql, 'where') === false) {
            throw new Exception('no WHERE condition in SQL(UPDATE, DELETE) to be executed! please make sure it\'s safe', 42005);
        }

        if ($this->allowRealExec || $sqlCmd == 'SELECT') {
            $cacheKey = md5($this->lastSql);
            if ($withCache && isset($this->cachedPageQueries[$cacheKey])) {
                return $this->cachedPageQueries[$cacheKey];
            }
            $this->queryBeginTime = microtime(true);
            if ($async) {
                $this->lastStmt = $this->connection->query($this->lastSql, $async);
            } else {
                $this->lastStmt = $this->connection->query($this->lastSql);
            }
        } else {
            $this->lastStmt = true;
        }
        $this->queryEndTime = microtime(true);
        $this->logQuery($this->lastSql);
        if (false === $this->lastStmt) {
            // connection broken, retry one time
            $errorInfo = $this->connection->errorInfo();
            if ($retryCount < 1 && $this->needConfirmConnection() && 2006 == $errorInfo[1]) {
                $retryCount += 1;
                $this->reConnect();
                $result = $this->query($sql);
                $retryCount = 0;
                return $result;
            }
            $retryCount = 0;
            $this->throwException('Query failure.SQL:' . $this->lastSql . '. ', 42004);
        }

        if (isset($cacheKey)) {
            $this->cachedPageQueries[$cacheKey] = $this->lastStmt;
        }
        return $this->lastStmt;
    }

    /**
     * 获取所有数据.
     *
     * @return array
     * @throws \Exception 查询失败.
     */
    public static function asyncFetchAll()
    {
        $allLinks = self::$asyncConnections;
        $asyncData = array();
        $processed = 0;
        do {
            $links = $errors = $reject = array();
            foreach ($allLinks as $link) {
                $links[] = $errors[] = $reject[] = $link;
            }
            if ( ! mysqli_poll($links, $errors, $reject, 1)) {
                continue;
            }
            foreach ($links as $link) {
                if ($result = $link->reap_async_query()) {
                    $asyncData = array_merge($asyncData,$result->fetch_all(MYSQLI_ASSOC));
                    $result->free();
                } else {
                    throw new \Exception("sync query failed!");
                }
                $processed++;
            }
        } while ($processed < count($allLinks));
        static::clearAsyncConnections();
        return $asyncData;
    }

    /**
     * 异步查询.
     *
     * @return mixed
     */
    public function asyncQuery()
    {
        // 暂时禁止直接SQL查询.
        $sql = null;
        $name = $this->currentRule->getAtomName();
        static::$asyncConnections[$name] = $this->connection;
        return $this->query($sql, true);
    }

    /**
     * 把Mysqli的Option转换为PDO需要的.
     *
     * @param array $options Mysqli设置选项.
     *
     * @return array
     */
    protected function genPdoOptions(array $options)
    {
        $pdoOptions = array();
        foreach ($options as $key => $op) {
            if (isset(static::$mysqliPdoOptions[$key])) {
                $pdoOptions[static::$mysqliPdoOptions[$key]] = $op;
            } else {
                new \Exception("Sharding Db 配置项中没有在 \DbShardingConnection::\$mysqliPdoOptions 找到对应的PDO项. MYSQLI KEY:$key");
            }
        }
        return $pdoOptions;
    }

    /**
     * Clear async connection.
     *
     * @return void
     */
    public static function clearAsyncConnections()
    {
        static::$writeAsyncConnections = array();
        static::$readAsyncConnections  = array();
        static::$asyncConnections = array();
    }

    /**
     * Close all connections.
     *
     * @return boolean
     */
    public function closeAll()
    {
        foreach (static::$readConnections as $k => $v) {
            $v = null;
            unset(static::$readConnections[$k]);
        }
        foreach (static::$writeConnections as $k => $v) {
            $v = null;
            unset(static::$writeConnections[$k]);
        }
        foreach (static::$writeAsyncConnections as $k => $v) {
            $v = null;
            unset(static::$writeAsyncConnections[$k]);
        }
        foreach (static::$readAsyncConnections as $k => $v) {
            $v = null;
            unset(static::$writeAsyncConnections[$k]);
        }
        return true;
    }

}
