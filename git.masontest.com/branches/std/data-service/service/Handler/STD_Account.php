<?php
/**
 * 账号数据主接口
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-10-22
 */

namespace Handler;

/**
 * 账号接口.
 */
class STD_Account extends Base
{
    /**
     * The registing new user api.
     *
     * @param string $username User name.
     * @param string $password User password.
     *
     * @return response array(newuser)
     */
    public function register($username, $password)
    {
        $model  = \Model\User::Instance();

        if ($error = $model->validate(array('username'=>$username, 'password'=>$password)) != "") { 
            $this->response['error']    = $error;
            return $this->response;
        }

        $newUser    = array(
            'username'  => $username,
            'password'  => $password,
            'created'   => time(),
            'user_role_id'  => 0,       //  the default user role is normal.
        );

        $userId = $model->addNew($newUser);

        $newUser['user_id'] = $userId; 

        $this->response['data'] = $newUser;
        return $this->response;
    }

    /**
     * To login account.
     *
     * @param string $username User name.
     * @param string $password User password.
     *
     * @return response array(session, user array).
     */
    public function login($username, $password)
    {
        $model  = \Model\User::Instance();

        if ($error = $model->validate(array('username'=>$username, 'password'=>$password)) != "") { 
            $this->response['error']    = $error;
            return $this->response;
        }

        $user   = $model->find(array('username' => $username));
        if (!$user) {
            $this->response['error']    = "Sorry, no this user";
            return $this->response;
        }

        if ($user['password'] != $password) {
            $this->response['error']    = "Sorry, username or password is wrong.";
            return $this->response;
        } 

        $response['user']               = $user;
        $response['session']            = 'testtestse';     // 这里会修改， 以后数据层不再负责session 会话.

        $this->response['data']         = $response;

        return $this->response;
    }

}
