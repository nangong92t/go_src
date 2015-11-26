<?php
/**
 * 用户数据模型.
 * 
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-10-22
 */

namespace Model;

/**
 * 用户数据模型.
 */
class User extends DbBase
{

    /**
     * Return model's table name.
     * 
     * @return string
     */
    public function getTableName()
    {
        return 'user';
    }

    /**
     * Returns the primary key of the associated database table.
     * 
     * @return string
     */
    protected function primaryKey()
    {
        return 'user_id';
    }

    /**
     * To validate the user input the data.
     *
     * @param array $data User data.
     *
     * @return string weathe wrong, add return the wrong infomation.
     */
    public function validate(array $data = array())
    {
        $error  = ''; 
        if (!$data['username']) {
            $error  = "Sorry, no username.";
        } elseif (!$data['password']) {
            $error  = "Sorry, no password.";
        }

        return $error;
    }

}
