<?php
/**
 * PHPClient的配置文件.
 * 
 * @author yongf <yongf@jumei.com>
 */

namespace Config;

/**
 * PHPClient的配置文件.
 */
class PHPClient
{
    public $rpc_secret_key = '769af463a39f077a0340a189e9c1ec28';
    public $monitor_log_dir = '/tmp/monitor-log';
    public $Vehicle= array(
        'uri' => 'tcp://127.0.0.1:2001',
        'user' => 'Optool',
        'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
    );

    public $User = array(
        'uri' => 'tcp://127.0.0.1:2002',
        'user' => 'Optool',
        'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
    );

    public $ResourceService = array(
        'uri' => 'tcp://127.0.0.1:2003',
        'user' => 'Optool',
        'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
    );

    public $BasicService    = array(
        'uri' => 'tcp://127.0.0.1:2004',
        'user' => 'Optool',
        'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
    );

    public $OrderService    = array( 
        'uri' => 'tcp://127.0.0.1:2005',
        'user' => 'Optool',
        'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
    );

}
