<?php
/**
 * 资源常量类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-08-03
 */

namespace CConst;

/**
 * 资源常量.
 */
class Resource
{

    /**
     * 住宿点酒店类型常量.
     */
    const HOTEL_TYPE_JIUDIAN = 1;

    /**
     * 住宿点宾馆类型常量.
     */
    const HOTEL_TYPE_BINGUAN = 2;

    /**
     * 住宿点旅店类型常量.
     */
    const HOTEL_TYPE_LVDIAN = 3;

    /**
     * 住宿点民放类型常量..
     */
    const HOTEL_TYPE_MINFANG = 4;

    /**
     * 住宿类型常量.
     */
    public static $HotelTypes    = array(
        self::HOTEL_TYPE_JIUDIAN    => '酒店',
        self::HOTEL_TYPE_BINGUAN    => '宾馆',
        self::HOTEL_TYPE_LVDIAN     => '旅店',
        self::HOTEL_TYPE_MINFANG    => '民房'
    );

    /**
     * 房间的套房类型.
     */
    const ROOM_TYPE_TAOFANG =   1;

    /**
     * 房间的单人间类型.
     */
    const ROOM_TYPE_DANREN  =   2;

    /**
     * 房间的双人间类型.
     */
    const ROOM_TYPE_SHUANGREN   = 3;

    /**
     * 房间的三人间类型.
     */
    const ROOM_TYPE_SANREN  = 4;

    /**
     * 房间的四人间类型.
     */
    const ROOM_TYPE_SIREN   = 5;

    /**
     * 房间的五人间类型.
     */
    const ROOM_TYPE_WUREN   = 6;

    /**
     * 房间的六人间类型.
     */
    const ROOM_TYPE_LIUREN  = 7;

    /**
     * 房间的七人间类型.
     */
    const ROOM_TYPE_QIREN   = 8;

    public static $RoomTypes    = array (
        self::ROOM_TYPE_TAOFANG     => '套房',
        self::ROOM_TYPE_DANREN      => '单人间',
        self::ROOM_TYPE_SHUANGREN   => '双人间',
        self::ROOM_TYPE_SANREN      => '三人间',
        self::ROOM_TYPE_SIREN       => '四人间',
        self::ROOM_TYPE_WUREN       => '五人间',
        self::ROOM_TYPE_LIUREN      => '六人间',
        self::ROOM_TYPE_QIREN       => '七人间'
    );
}
