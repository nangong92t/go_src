<?php
namespace Config;

class Redis{
    /**
     * Configs of Redis.
     * @var array
     */
    public $default = array('nodes' => array(
            array('master' => "192.168.8.230:27003", 'slave' => "192.168.8.231:27003"),
            array('master' => "192.168.8.230:27004", 'slave' => "192.168.8.231:27004"),

    ),
            'db' => 0
    );
    public $fav = array('nodes' => array(
            array('master' => "192.168.25.9:6379", 'slave' => "192.168.25.9:6379"),
    ),
            'db' => 2
    );
}