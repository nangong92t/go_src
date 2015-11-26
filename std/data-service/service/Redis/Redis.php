<?php
/**
 * RedisBase.
 * 
 * @author hongfengd <hongfengd@jumei.com>
 */

namespace Redis;

/**
 * Redis.
 */
class Redis
{
    
    /**
     * 得到Cache 级别Redis.
     * 
     * @return RedisCache
     */
    public static function getCache()
    {
        return RedisMultiCache::getInstance('default');
    }

    /**
     * 得到Storage 级别Redis.
     * 
     * @return RedisStorage
     */
    public static function getStorage()
    {
        return RedisMultiStorage::getInstance('storage');
    }

}
