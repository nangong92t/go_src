<?php
/**
 * \Component\Util
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-09-12
 */

namespace Component;

/**
 * Class \Component\Util.
 */
class Util
{
    
    /**
     * 截取html文本摘要.
     * 
     * @param string  $content    内容.
     * @param integer $charsCount 截取个数.
     * 
     * @return string
     */
    public static function getAbstractInputFromHtml($content, $charsCount)
    {
        $content = html_entity_decode(strip_tags(stripslashes($content)), ENT_QUOTES, 'UTF-8');
        $content = preg_replace('/[\x00-\x20]+/', ' ', $content); // WARNING: /\s+/ in php 5.2.9 is buggy, which will replace 0xA0 to 0x20
        $content = str_replace('　', ' ', $content);
        $content = preg_replace("/(^\s+)|(\s+$)/us", '', $content);
        $content = trim($content);
        $content = self::truncate($content, $charsCount);
        return addslashes(trim($content));
    }
    
    /**
     * 字符串截取.
     * 
     * @param string  $string      字符串.
     * @param integer $length      截取长度.
     * @param string  $etc         后缀显示.
     * @param boolean $break_words 是否短词.
     * @param boolean $middle      是否中间截取.
     * 
     * @return string
     */
    public static function truncate($string, $length = 80, $etc = '...', $break_words = false, $middle = false)
    {
        if ($length == 0) {
            return '';
        }
        $string = trim(strip_tags($string));

        if (is_callable('mb_strlen')) {

            if (mb_detect_encoding($string, 'UTF-8, ISO-8859-1') === 'UTF-8') {

                // $string has utf-8 encoding
                if (mb_strlen($string, 'UTF-8') > $length) {

                    if (!$break_words && !$middle) {
                        $string = preg_replace('/\s+?$/u', '', mb_substr($string, 0, $length + 1, 'UTF-8'));
                    }

                    if (!$middle) {

                        return mb_substr($string, 0, $length, 'UTF-8') . $etc;
                    } else {
                        return mb_substr($string, 0, $length / 2, 'UTF-8') . $etc . mb_substr($string, - $length / 2, 'UTF-8');
                    }
                } else {
                    return $string;
                }
            }
        }
        // $string has no utf-8 encoding
        if (strlen($string) > $length) {
            $length -= min($length, strlen($etc));
            if (!$break_words && !$middle) {
                $string = preg_replace('/\s+?(\S+)?$/', '', substr($string, 0, $length + 1));
            }
            if (!$middle) {
                return substr($string, 0, $length) . $etc;
            } else {
                return substr($string, 0, $length / 2) . $etc . substr($string, - $length / 2);
            }
        } else {
            return $string;
        }
    }
    
    /**
     * 得到数组中一个键的值.
     * 
     * @param array  $arr   数组.
     * @param string $field Key.
     * 
     * @return array
     */
    public static function getArrayFieldValue(array $arr, $field = 'id')
    {
        if (!empty($arr) && is_array($arr)) {
            $keys = array_map(
                function ($a) use ($field) {
                    return isset($a[$field]) ? $a[$field] : null;
                }
                , $arr
            );
            return $keys;
        }
        return array();
    }

}
