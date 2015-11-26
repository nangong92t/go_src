<?php
/**
 * Album Helper帮助类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-07-06
 */

namespace Helper;

/**
 * Album Helper.
 */
class Album
{

    /**
     * 通过Attach模型数据获取图片地址.
     *
     * @param array  $attach Attach数组源.
     *
     * @return string.
     */
    public static function getPictureUrl(array $attach = array())
    {
        $result = 'http://' . $attach['savehost'] . '/' . 
            $attach['savepath'] . '/' . $attach['savename'] . '.' . $attach['extenstion'];
        return $result;
    }

}
