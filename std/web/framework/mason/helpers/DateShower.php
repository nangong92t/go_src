<?php
/**
 * 日期时间
 * @package nlsp.helper
 */
class DateShower 
{
    public static function formatWeek($timestamp = '', $pre = '星期')
    {
        if (!$timestamp) $timestamp = time();
        
        switch (date('N'))
        {
            case '1': $str .= '一'; break;
            case '2': $str .= '二'; break;
            case '3': $str .= '三'; break;
            case '4': $str .= '四'; break;
            case '5': $str .= '五'; break;
            case '6': $str .= '六'; break;
            case '7': $str .= '日'; break;
        }
        
        return $pre . $str;
    }
    
	/**
     * timestamp转换成显示时间格式
     * @param $timestamp
     * @return unknown_type
     */
    public static function tTimeFormat($timestamp) {
        $curTime = time();
        $space = $curTime - $timestamp;
        //1分钟
        if ($space < 60) {
            $string = "刚刚";
            return $string;
        } elseif ($space < 3600) { //一小时前
            $string = floor($space / 60) . "分钟前";
            return $string;
        }
        $curtimeArray = getdate($curTime);
        $timeArray = getDate($timestamp);
        if ($curtimeArray['year'] == $timeArray['year']) {
            if ($curtimeArray['yday'] == $timeArray['yday']) {
                $format = "%H:%M";
                $string = strftime($format, $timestamp);
                return "今天 {$string}";
            } elseif (($curtimeArray['yday'] - 1) == $timeArray['yday']) {
                $format = "%H:%M";
                $string = strftime($format, $timestamp);
                return "昨天 {$string}";
            } else {
                $string = sprintf("%d月%d日 %02d:%02d", $timeArray['mon'], $timeArray['mday'], $timeArray['hours'], $timeArray['minutes']);
                return $string;
            }
        }
        $string = sprintf("%d年%d月%d日 %02d:%02d", $timeArray['year'], $timeArray['mon'], $timeArray['mday'], $timeArray['hours'], $timeArray['minutes']);
        return $string;
    }
}
