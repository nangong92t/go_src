<?php

class Controller_Query extends Controller_OptoolBase
{

    public function action_Test()
    {
        header('Content-Type: text/plain; charset=utf-8');

        // {{ 示例代码:
        \PHPClient\Text::instance('User_Info')
        $userInfo = RpcClient_User_Info::instance();

        $result = $userInfo->byUid(5100);
        var_dump($result);

        $result = $userInfo->getInfoByUid(1373);
        var_dump($result);

        // }}
    }
}
