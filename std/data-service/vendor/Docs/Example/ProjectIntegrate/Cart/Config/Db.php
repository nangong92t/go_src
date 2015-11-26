<?php
namespace Config;

class Db{
    /**
     * Configs of database.
     * @var array
     */
    public $read = array(
            'stats' => array('dsn'      => 'mysql:host=192.168.20.72;port=9001;dbname=tuanmei_stats',
                    'user'     => 'dev',
                    'password' => 'jmdevcd',
                    'confirm_link' => true,//required to set to TRUE in daemons.
                    'options'  => array(\PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'latin1\'',
                            \PDO::ATTR_TIMEOUT=>3
                    )
            ),
            'jumei' => array('dsn'      => 'mysql:host=192.168.20.72;port=9001;dbname=jumei',
                    'user'     => 'dev',
                    'password' => 'jmdevcd',
                    'confirm_link' => true,//required to set to TRUE in daemons.
                    'options'  => array(\PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'latin1\'',
                            \PDO::ATTR_TIMEOUT=>3
                    )
            )
    );

    public $write = array(
            'stats'=>array('dsn'      => 'mysql:host=192.168.20.71;port=9001;dbname=tuanmei_stats',
                    'user'     => 'dev',
                    'password' => 'jmdevcd',
                    'confirm_link' => true,//required to set to TRUE in daemons.
                    'options'  => array(\PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'latin1\'',
                            \PDO::ATTR_TIMEOUT=>3
                    )
            ),
            'jumei'=>array('dsn'      => 'mysql:host=192.168.20.71;port=9001;dbname=jumei',
                    'user'     => 'dev',
                    'password' => 'jmdevcd',
                    'confirm_link' => true,//required to set to TRUE in daemons.
                    'options'  => array(\PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'latin1\'',
                            \PDO::ATTR_TIMEOUT=>3
                    )
            )

    );

}