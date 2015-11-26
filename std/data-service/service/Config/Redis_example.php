<?php
/**
 * Redis的配置文件.
 * 
 * @author yongf <yongf@jumei.com>
 */

namespace Config;

/**
 * Redis的配置文件.
 */
class Redis
{
    /**
     * Configs of Redis.
     * @var array
     */
    public $default = array(
        'nodes' => array(
            array('master' => "192.168.16.42:6379", 'slave' => "192.168.16.42:6379"),
        ),
        'db' => 0
    );
    public $storage = array(
        'nodes' => array(
            array('master' => "192.168.16.42:6379", 'slave' => "192.168.16.42:6379"),
        ),
        'db' => 2
    );

}
