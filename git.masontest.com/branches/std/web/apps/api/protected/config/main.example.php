<?php

// This is the main Web application configuration. Any writable
// CWebApplication properties can be configured here.
return array(
    'timeZone'=>'Asia/Shanghai',
    'basePath'=>dirname(__FILE__).DIRECTORY_SEPARATOR.'..',
    'name'=>'STD project Api proxy',

    // preloading 'log' component
    'preload'=>array('log'),

    // autoloading model and component classes
    'import'=>array(
        'system.mason.components.*',
        'system.mason.helpers.*',
        'application.components.*',
    ),

    'modules'=>array(
        // uncomment the following to enable the Gii tool
        /*
        'gii'=>array(
            'class'=>'system.gii.GiiModule',
            'password'=>'Enter Your Password Here',
            // If removed, Gii defaults to localhost only. Edit carefully to taste.
            'ipFilters'=>array('127.0.0.1','::1'),
        ),
        */
    ),

    // application components
    'components'=>array(
        'user'=>array(
            // enable cookie-based authentication
            'allowAutoLogin'=>true,
        ),
        // uncomment the following to enable URLs in path-format
        'urlManager'=>array(
            'urlFormat'=>'path',
            'rules'=>array(
                // AUTH routes
                array('auth/main', 'pattern'=>'auth', 'verb'=>'GET'),

                // UPLOADER routers
                array('upload/file', 'pattern'=>'upload/*', 'verb'=>'POST'),
                array('upload/allowCross', 'pattern'=>'crossdomain.xml', 'verb'=>'GET'),

                // REST routers
                array('rest/view', 'pattern'=>'v1/<rpc>/<action>/<params>', 'verb'=>'GET'),
                array('rest/view', 'pattern'=>'v1/<token>/<rpc>/<action>/<params>', 'verb'=>'GET'),
                array('rest/create', 'pattern'=>'v1/<token>/<rpc>/<action>', 'verb'=>'POST'),
                array('rest/update', 'pattern'=>'v1/<token>/<rpc>/<action>', 'verb'=>'PUT'),
                array('rest/delete', 'pattern'=>'v1/<token>/<rpc>/<action>/<params>', 'verb'=>'DELETE'),
            ),
        ),
        'errorHandler'=>array(
            // use 'site/error' action to display errors
            'errorAction'=>'auth/error',
        ),
        'log'=>array(
            'class'=>'CLogRouter',
            'routes'=>array(
                array(
                    'class'=>'CFileLogRoute',
                    'levels'=>'error, warning',
                ),
                // uncomment the following to show log messages on web pages
                /*
                array(
                    'class'=>'CWebLogRoute',
                ),
                */
            ),
        ),
        'clientScript' => array(
            'coreScriptUrl' => 'http://static.std.com',
            'scriptMap' => array(
                'init.js' => 'http://static.std.com/js/init.js'
            )
        ),
    ),

    // application-level parameters that can be accessed
    // using Yii::app()->params['paramName']
    'params'=>array(
        // this is used in contact page
        'adminEmail'=>'webmaster@example.com',
        'rpc' => array(
            'rpc_secret_key' => 'ab1f8e61026a7456289c550cb0cf77cda44302b4',
            'STD' => array(
                'uri' => 'tcp://127.0.0.1:9001',
                'user' => 'Optool',
                'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
            ),
        ),
    ),
);
