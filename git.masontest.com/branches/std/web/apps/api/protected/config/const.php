<?php

class ProjectParams
{
    /**
     * 菜单栏
     */
    public static $menus  = array(
        'admin/main'        => array(
            'name'  => '首页'
        ),
        'resource/index'    => array(
            'name'  => '资源管理',
            'subMenus'  => array(
                'resource/index'        => array(
                    'name'  => '资源汇总'
                ),
                'vehicle/index'         => array(
                    'name'  => '车辆管理',
                    'subMenus'  => array(
                        'vehicle/index'     => array(
                            'name'  => '车辆列表',
                        ),
                        'vehicle/typelist'      => array(
                            'name'  => '车辆类型',
                        ),
                        'vehicle/outlook'   => array(
                            'name'  => '车辆外形摸板'
                        ),
                    )
                ),
                'localfronter/index'    => array(
                    'name'  => '地接商管理',
                ),
            )
        ),
        'user/index'    => array(
            'name'  => '用户管理'
        ),
        'role/index'    => array(
            'name'  => '角色与权限'
        ),
        'admin/logout'  => array(
            'name'  => '退出'
        )
    );

}


/**
 *  * 定义一些资源常量
 */
class ProfileItems
{
    /**
     * 性别
     */
    const MAN   = 1;

    const WOMAN = 2;

    public static $gender = array(
        self::MAN   => '男性',
        self::WOMAN => '女性'
    );

    /**
     * 艺人身高
     *
     */
    const MINHEIGHT = 50;
    const MAXHEIGHT = 250;

    /**
     * 血型
     *
     */
    public static $bloods = array(
        1 => 'A', 2 => 'B', 3 => 'AB', 4 => 'O'
    );


}
