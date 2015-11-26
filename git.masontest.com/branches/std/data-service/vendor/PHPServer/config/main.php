<?php

return array(

    'workers' => array(
        // LY Vehicle 车辆管理功能模块
        'STD' => array(
            'protocol'              => 'tcp',    // [必填]tcp udp
            'port'                  => 9001,     // [必填]监听的端口
            'child_count'           => 10,       // [必填]worker进程数 注意:每个进程大概占用30M内存，总worker数量不要超过800个
            'ip'                    => "0.0.0.0",// [选填]绑定的ip，注意：生产环境要绑定内网ip 不配置默认是0.0.0.0
            'recv_timeout'          => 1000,     // [选填]从客户端接收数据的超时时间          不配置默认1000毫秒
            'process_timeout'       => 30000,    // [选填]业务逻辑处理超时时间               不配置默认30000毫秒
            'send_timeout'          => 1000,     // [选填]发送数据到客户端超时时间            不配置默认1000毫秒
            'persistent_connection' => true,     // [选填]是否是长连接                      不配置默认是短链接（短连接每次请求后服务器主动断开）
            'max_requests'          => 1000,     // [选填]进程接收多少请求后退出              不配置默认是0，不退出
            'worker_class'          => 'JmTextWorker',// worker使用的类
            'bootstrap'             => '../../../service/init.php', // 进程初始化时调用一次，可以在这里做些全局的事情，例如设置autoload
        ),

        // 统计接口调用结果 只开一个进程 已经配置好，不用设置
        /*
        'StatisticWorker' => array(
            'protocol'              => 'udp',
            'port'                  => 2207,
            'child_count'           => 1,
        ),

        // 查询接口调用结果 只开一个进程 已经配置好，不用再配置
        'StatisticService' => array(
            'protocol'              => 'tcp',
            'port'                  => 20202,
            'child_count'           => 1,
        ),
         */

        // 查询接口调用结果 只开一个进程 已经配置好，不用再配置
        'StatisticGlobal' => array(
            'protocol'              => 'tcp',
            'port'                  => 20203,
            'child_count'           => 1,
        ),

        // 查询接口调用结果 只开一个进程 已经配置好，不用再配置
        'StatisticProvider' => array(
            'protocol'              => 'tcp',
            'port'                  => 20204,
            'child_count'           => 1,
        ),

        // 监控server框架的worker 只开一个进程 framework里面需要配置成线上参数
        /*
        'Monitor' => array(
            'protocol'              => 'tcp',
            'port'                  => 20305,
            'child_count'           => 1,
            'framework'             => array(
                 'phone'   => '15551251335,15551251335',      // 告警电话
                 'url'     => 'http://sms.jumeicd.com/send',  // 发送短信调用的url
                 'param'   => array(                          // 发送短信用到的参数
                     'channel' => 'monternet',
                     'key'     => 'tester_123456',
                     'task'    => 'test',
                 ),
                 'min_success_rate' => 98,                    // 框架层面成功率小于这个值时触发告警
                 'max_worker_normal_exit_count' => 1000,      // worker进程退出（退出码为0）次数大于这个值时触发告警
                 'max_worker_unexpect_exit_count' => 10,      // worker进程异常退出（退出码不为0）次数大于这个值时触发告警
             )
         ), */

        // [开发环境用，生产环境可以去掉该项]耗时任务处理，发送告警短信 邮件，监控master进程是否退出,开发环境监控文件更改等
        'TaskWorker' => array(
            'protocol'              => 'udp',
            'port'                  => 10203,
            'child_count'           => 1,
        ),

        // [开发环境用，生产环境可以去掉该项]rpc web测试工具
        'TestClientWorker' => array(
            'protocol'              => 'tcp',
            'port'                  => 30303,
            'child_count'           => 1,
        ),

	/*
        // [开发环境用，生产环境可以去掉该项]thrift rpc web测试工具
        'TestThriftClientWorker' => array(
            'protocol'              => 'tcp',
            'port'                  => 30304,
            'child_count'           => 1,
        ),

        // Thrift Worker
        'ThriftWorker' => array(                                          // 注意：键名固定为服务名
            'protocol'              => 'tcp',                             // 固定tcp
            'port'                  => 9090,                              // 每组服务一个端口
            'child_count'           => 1,                                 // 启动多少个进程提供服务
            'persistent_connection' => true,                              // thrift默认使用长链接
            'provider'              => '../../../Provider',                  // 这里是thrift生成文件所放目录,可以是绝对路径
            'handler'               => '../../../Handler',                   // 这里是对thrift生成的Provider里的接口的实现
            'bootstrap'             => '../../../init.php',             // 进程启动时会载入这个文件，里面可以做一些autoload等初始化工作
            'monitor_log_path'      => '/var/logs/jm/monitor',
        ),*/
    ),

    'ENV'          => 'dev', // dev or production
    'worker_user'  => '', //运行worker的用户,正式环境应该用低权限用户运行worker进程

    // 数据签名用私匙
    'rpc_secret_key'    => 'ab1f8e61026a7456289c550cb0cf77cda44302b4',

);
