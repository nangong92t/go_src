<?php
/**
 * LY Module Level
 *
 * @author Tony<tonycbcd@gmail.com>
 * @date 2014-07-07
 */

namespace Module;

/**
 * Abstruct model class.
 */
class ModuleBase
{
    const SERIALIZE_BY_JSON = 'json';
    const SERIALIZE_BY_PHP = 'php';
    const SERIALIZE_BY_MSGPACK = 'msgpack';

    /**
     * default cache time.
     */
    public $ttl = 1800;
 
    /**
     * default object cache
     */
    protected static $instances_ = array();

    /**
     * To get a self single object.
     *
     * @return static
     */
    public static function Instance()
    {
        $class = get_called_class();
        if (!isset(static::$instances_[$class])) {
            static::$instances_[$class] = new $class;
        }
        $curObj = static::$instances_[$class];
        if (method_exists($curObj, 'init')) {
            $curObj->init();
        }
        return $curObj;
    }
    
    /**
     * 缓存键生成方法，需要获得缓存的键的时候，请使用本方法，无毒副作用.
     * 
     * @param string $preFix 前缀key.
     * @param string $key    Key.
     * 
     * @return string
     */
    protected function generateKey($preFix, $key)
    {
        return $preFix . $key;
    }

    /**
     * 序列化PHP的数据结构(如果是序列化对象,请不要使用json序列).
     * 
     * @param array  $data 需要序列化的数据.
     * @param string $type 序列化方式.
     *
     * @return string
     */
    protected function serializeData(array $data = array(), $type = self::SERIALIZE_BY_MSGPACK)
    {
        if (!$data) {
            $data = array();
        }
        $dataString = '';
        if ($type == self::SERIALIZE_BY_JSON) {
            $dataString = json_encode($data);
        } elseif ($type == self::SERIALIZE_BY_PHP) {
            $dataString = serialize($data);
        } else {
            $dataString = msgpack_pack($data);
        }

        return $dataString;
    }
    
    /**
     * 反序列化数据(如果是序列话对象,请不要使用json序列).
     * 
     * @param string $string 需要反序列化的字符串.
     * @param string $type   序列化方式.
     *
     * @return array
     */
    protected function unSerializeData($string = '', $type = self::SERIALIZE_BY_MSGPACK)
    {
        if (!$string) {
            return array();
        }
        $result = array();
        if ($type == self::SERIALIZE_BY_JSON) {
            if (!$result = @json_decode($string, true)) {
                $result = array();
            }
        } elseif ($type == self::SERIALIZE_BY_PHP) {
            if (!$result = @unserialize($string)) {
                $result = array();
            }
        } else {
            if (!$result = @msgpack_unpack($string)) {
                $result = array();
            }
        }

        return $result;
    }

}
