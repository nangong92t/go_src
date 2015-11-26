<?php
namespace Redis;

/*
 * 如果你在一个项目里面用到了很多个集群，那么用这个
 */

/**
 * Description of RedisMultiStorage
 *
 * @author guoxinhua
 */
class RedisMultiStorage {

    private static $instance;
    private static $config;

    public static function getInstance($name) {
        if(!self::$config)
        {
            self::$config = (array) new \Config\Redis();
        }
        if (!isset(self::$instance[$name])) {
            RedisStorage::config(self::$config, $name);
            self::$instance[$name] = RedisStorage::getInstance($name);
        }
        return self::$instance[$name];
    }

    public static function config(array $config) {
        self::$config = $config;
    }

    public static function close(){
        foreach ((array)self::$instance as $inst) {
            $inst->close();
        }
        self::$instance = array();
    }

}
