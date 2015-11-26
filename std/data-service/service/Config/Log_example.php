<?php
/**
 *日志的配置文件.
 *
 * @author rongyix<rongyix@jumei.com>
 */

namespace Config;

/**
 *日志的配置文件.
 */
class Log
{
    // 文件日志的根目录.请确认php进程对此目录可写
    public $FILE_LOG_ROOT = '/var/log/koubei/';

    // 数据库日志配置
    public $db = array(
        'logger' => 'file',
    );

    // Sphinc 日志配置.
    public $sphinx  = array(
        'logger'    => 'file',
    );

}
