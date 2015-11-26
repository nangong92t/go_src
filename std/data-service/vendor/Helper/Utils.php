<?php
/**
 * Utils Helper 工具类帮助类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-07-06
 */

namespace Helper;

/**
 * Utils Helper.
 */
class Utils
{

    /**
     * 通过Key在数组中找出对应Key Id Array.
     *
     * @param array  $dataSource 数组源.
     * @param string $fieldName  过滤字段名.
     *
     * @return array.
     */
    public static function getFieldIds(array $dataSource = array(), $fieldName)
    {
        $result = array();

        foreach ($dataSource as $one) {
            $result[]   = $one[$fieldName];
        }

        return $result;
    }

    /**
     * 根据字段Id，为母数组填加增加字段数据.
     *
     * @param array    $dataSource             源数组, 最好传入地址.
     * @param array    $child                  需填加数据数组.
     * @param string   $childFieldInDataSource 增加到源数组的字段名.
     * @param string   $childFieldName         需填加数据数组中的哪个字段.
     * @param string   $fieldName              填加数据数组中的哪个字段.
     * @param function $renderFun              渲染方法.
     *
     * @return array.
     */
    public static function addField(array &$dataSource = array(), array $child = array(), $childFieldInDataSource, $childFieldName, $fieldName, $renderFun = null)
    {
        $map    = array();

        // get map.
        foreach ($child as $one) {
            $map[ $one[$childFieldName] ] = $one;
        }

        // bind value
        foreach ($dataSource as $i => $one) {
            if (!isset($one[$fieldName]) || !isset($map[ $one[$fieldName] ])) { 
                $dataSource[$i][$childFieldInDataSource]    = '';
                continue;
            }
            $data   = $map[ $one[$fieldName] ];
            if ($renderFun) { $data   = $renderFun($data); }
            $dataSource[$i][$childFieldInDataSource]    = $data;
        }

        return $dataSource;
    }

    /**
     * 通过元素key获取元素value.
     *
     * @param array  $data  数组.
     * @param array  $field 元素key.
     *
     * @return string.
     */
    public static function getArrayElement(array $data = array(), $field)
    {
        return $data[$field];
    }

    /**
     * 获取元素value字符串.
     *
     * @param array  $data  数组.
     * @param array  $field 元素key.
     *
     * @return string.
     */
    public static function getArrayElementString(array $data = array(), $field)
    {
        $str = '';
        foreach ($data as $one) {
            if (!isset($one[$field])) { continue; }

            $str = $str . $one[$field] . ' ';
        }

        return $str;
    }
}

