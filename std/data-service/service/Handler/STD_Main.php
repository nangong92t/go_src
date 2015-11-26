<?php
/**
 * 基础数据主接口
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-09-12
 */

namespace Handler;

/**
 * 主接口.
 */
class STD_Main extends Base
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
