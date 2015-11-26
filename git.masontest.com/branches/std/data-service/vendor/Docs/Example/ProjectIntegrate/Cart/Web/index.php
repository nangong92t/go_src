<?php
require __DIR__.'/../init.php';
\Model\Details::instance()->updateItemCount(3, 11, 3);
\Model\Details::instance()->testDbConnection();
\Model\Details::instance()->testLog();
\Model\Details::instance()->testRedis();

var_dump(\PHPClient\Text::inst('User')->setClass('Info')->getUserInfobyUid(5100));
var_dump(\PHPClient\Text::inst('User')->setClass('Address')->getListByUid(5100));
var_dump(\PHPClient\Text::inst('Example')->setClass('Example')->sayHello(5100));
var_dump(\PHPClient\Text::inst('Example')->setClass('Example\C1')->fn(5100));
echo "--------\n";
var_dump(RpcClient_User_Address::instance()->getListByUid(5100));