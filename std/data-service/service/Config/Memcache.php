<?php
/**
 * Memcache的配置文件.
 * 
 * @author yongf <yongf@jumei.com>
 */

namespace Config;

/**
 * Memcache的配置文件.
 */
class Memcache
{
    public $default = array(
            array('host' => '192.168.20.95', 'port' => 6660, 'type' => 'backup'),
            array('host' => '127.0.0.1', 'port' => 11211)
    );
}
