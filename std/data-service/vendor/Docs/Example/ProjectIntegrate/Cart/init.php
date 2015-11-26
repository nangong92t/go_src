<?php
define('ROOT_PATH', __DIR__.DIRECTORY_SEPARATOR);
require ROOT_PATH.'../../../../../Vendor/Bootstrap/Autoloader.php';
\Bootstrap\Autoloader::instance()->addRoot(ROOT_PATH)->init();