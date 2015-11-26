<?php
/**
 * 初始化
 * 
 * @author XuRongyi <tonycbcd@gmail.com>
 * @date 2014-06-22
 */

define('ROOT_PATH', __DIR__.DIRECTORY_SEPARATOR);

require ROOT_PATH.'../vendor/Bootstrap/Autoloader.php';
\Bootstrap\Autoloader::instance()->addRoot(ROOT_PATH)->init();

function on_phpserver_request_finish()
{
    if(class_exists('\Db\Connection', false) && is_callable(array(\Db\Connection::instance(), 'closeAll')))
    {
        \Db\Connection::instance()->closeAll();
    }

    if(class_exists('\Redis\RedisMultiCache', false))
    {
        \Redis\RedisMultiCache::close();
    }

    if(class_exists('\Redis\RedisMultiStorage', false))
    {
        \Redis\RedisMultiStorage::close();
    }
}

