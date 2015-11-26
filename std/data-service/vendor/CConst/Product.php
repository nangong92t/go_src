<?php
/**
 * 业务常量类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-07-31
 */

namespace CConst;

/**
 * 产品常量.
 */
class Product
{

    /**
     * 景点客运.
     */
    const TYPE_PASSENGER_TRANSPORT = 1;

    /**
     * 餐饮点套餐.
     */
    const TYPE_MEAL = 2;

    /**
     * 住宿点房间.
     */
    const TYPE_ROOM = 4;

    /**
     * 景点门票.
     */
    const TYPE_SCENIC_TICKET = 8;

    /**
     * 车辆.
     */
    const TYPE_VEHICLE = 16;

    /**
     * 座位.
     */
    const TYPE_VEHICLE_SEAT = 32;

    /**
     * 平台打包产品.
     */
    const TYPE_PLATFORM_PACKAGE_PRODUCT = 64;

    /**
     * 用户打包产品.
     */
    const TYPE_USER_PACKAGE_PRODUCT = 128;

    /**
     * 导游服务类型产品.
     */
    const TYPE_GUIDE = 256;

    /**
     * 纯服务类型产品.
     */
    const TYPE_SERVICE = 512;


    /**
     * 资源类型常量.
     */
    public static $Types    = array(
        self::TYPE_PASSENGER_TRANSPORT => '景点客运',
        self::TYPE_MEAL => '餐饮点套餐',
        self::TYPE_ROOM => '住宿点房间',
        self::TYPE_SCENIC_TICKET => '景点门票',
        self::TYPE_VEHICLE => '车辆',
        self::TYPE_VEHICLE_SEAT => '座位',
        self::TYPE_PLATFORM_PACKAGE_PRODUCT => '平台打包产品',
        self::TYPE_USER_PACKAGE_PRODUCT => '用户打包产品',
        self::TYPE_GUIDE => '导游',
        self::TYPE_SERVICE => '服务'
    );

}
