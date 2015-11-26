<?php
/**
 * 用户常量类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-08-16
 */

namespace CConst;

/**
 * 用户常量.
 */
class User
{

    /**
     * 男性.
     */
    const MALE = 1;

    /**
     * 女性.
     */
    const FEMALE = 2;

    public static $Gender    = array(
        self::MALE      => '男性',
        self::FEMALE    => '女性'
    );

}
