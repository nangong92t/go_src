<?php
/**
 * 用户数据模型.
 * 
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-09-12
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

}
