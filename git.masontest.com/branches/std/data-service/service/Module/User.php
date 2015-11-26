<?php
/**
 * Class Module\User
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-10-22
 */

namespace Module;

/**
 * User.
 */
class User extends ModuleBase
{
    
    /**
     * 获取用户列表..
     * 
     * @param array $userId 用户Id.
     * 
     * @return array
     */
    public function getUsers($userId) 
    {
        $model      = \Model\User::Instance();
        $user       = $model->getAlbumWithCoverByIds($userId);

        return $user;
    }

}
