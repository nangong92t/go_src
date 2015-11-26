<?php
/**
 * 数据库配置文件.
 * 
 * @author Tonyxu<tonycbcd@gmail.com>
 * @date 2014-10-20
 */

namespace Config;

/**
 * 数据库配置文件.
 */
class Db
{
    /**
     * Configs of database.
     * @var array
     */
    public $read = array(
        'std'    => array(
            'dsn'      => 'mysql:host=localhost;port=3306;dbname=std',
            'user'     => 'std',
            'password' => 'std',
            'confirm_link' => true, // required to set to TRUE in daemons.
            'options'  => array(
                \PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'utf8\'',
                \PDO::MYSQL_ATTR_USE_BUFFERED_QUERY => true,
                \PDO::ATTR_TIMEOUT => 3
            )
        )

    );

    public $write = array(
        'std'    => array(
            'dsn'      => 'mysql:host=localhost;port=3306;dbname=std',
            'user'     => 'std',
            'password' => 'std',
            'confirm_link' => true, // required to set to TRUE in daemons.
            'options'  => array(
                \PDO::MYSQL_ATTR_INIT_COMMAND => 'SET NAMES \'utf8\'',
                \PDO::MYSQL_ATTR_USE_BUFFERED_QUERY => true,
                \PDO::ATTR_TIMEOUT => 3
            )
        )

    );

}
