<?php
/**
 * 用户相关主接口
 *
 * @author TonyXu<tonycbcd@gmail.com> 
 * @date 2014-10-20
 */

namespace Handler;

/**
 * 用户主接口.
 */
class User_Main extends Base
{

    /**
     * 通过用户Id获取用户信息.
     *
     * @param integer $uid Uid.
     *
     * @return array
     */
    public function getUserById($uid)
    {
        $this->response['data'] = "Hello Tony, id:".$uid;
        return $this->response;
    }

}
