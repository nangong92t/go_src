<?php
namespace Config;

class Memcache{
    public $default = array(
            array('host'=>'192.168.20.95', 'port'=>6660, 'type'=>'backup'),
            array('host'=>'127.0.0.1', 'port' => 11211)
    );
}