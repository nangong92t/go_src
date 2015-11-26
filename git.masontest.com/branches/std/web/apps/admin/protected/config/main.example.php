<?php

// uncomment the following to define a path alias
// Yii::setPathOfAlias('local','path/to/local-folder');

Yii::import('system.mason.helpers.functions', true);

// 定义的一些常量
require_once(realpath(dirname(__FILE__) . '/const.php'));

$COOKIE_SITE_DOMAIN = strstr($_SERVER['SERVER_NAME'],'.');
define('SITE_DOMAIN', 'www.std.com');
define('API_DOMAIN', 'api.std.com');
define('COOKIE_SITE_DOMAIN', $COOKIE_SITE_DOMAIN);
ini_set('session.name', 'PHPSESSID-DEV');
ini_set('session.cookie_domain', COOKIE_SITE_DOMAIN);

// This is the main Web application configuration. Any writable
// CWebApplication properties can be configured here.
return array(
	'basePath'=>dirname(__FILE__).DIRECTORY_SEPARATOR.'..',
    'defaultController'=>'Admin',
	'name'=>'STD Admin System',

	// preloading 'log' component
	'preload'=>array('log'),

	// autoloading model and component classes
	'import'=>array(
        'system.mason.components.*',
        'system.mason.helpers.*',
        'system.mason.widgets.*',
        'system.mason.widgets.grid.*',
		'application.models.*',
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
            'class' => 'MasonWebUser',
            'identityCookie' => array('domain' => COOKIE_SITE_DOMAIN),
		),
		// uncomment the following to enable URLs in path-format
		/*
		'urlManager'=>array(
			'urlFormat'=>'path',
			'rules'=>array(
				'<controller:\w+>/<id:\d+>'=>'<controller>/view',
				'<controller:\w+>/<action:\w+>/<id:\d+>'=>'<controller>/<action>',
				'<controller:\w+>/<action:\w+>'=>'<controller>/<action>',
			),
		),
		'db'=>array(
			'connectionString' => 'sqlite:'.dirname(__FILE__).'/../data/testdrive.db',
		),
		// uncomment the following to use a MySQL database
		'db'=>array(
			'connectionString' => 'mysql:host=localhost;dbname=std',
			'emulatePrepare' => true,
			'username' => 'std',
			'password' => 'std',
			'charset' => 'utf8',
		),
		*/
		'errorHandler'=>array(
			// use 'site/error' action to display errors
			'errorAction'=>'admin/error',
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
        'cookie' => array(
          'class' => 'StdHttpCookie',
          'domain' => COOKIE_SITE_DOMAIN,
          'path' => '/',
          //'expire' => 0,
          //'secure' => 0,
          //'httpOnly' => 0,
        ),
	),

	// application-level parameters that can be accessed
	// using Yii::app()->params['paramName']
	'params'=>array(
		// this is used in contact page
		'adminEmail'=>'webmaster@example.com',
        'rpc' => array(
            'rpc_secret_key' => '769af463a39f077a0340a189e9c1ec28',
            'STD' => array(
                'uri' => 'tcp://127.0.0.1:3001',
                'user' => 'Optool',
                'secret' => '{1BA09530-F9E6-478D-9965-7EB31A59537E}',
                //'compressor' => 'GZ',
            ),

        ),
	),
);
