<?php
/**
 * Class Module\Event
 *
 * @author hongfengd<hongfengd@jumei.com>
 */

namespace Module;

/**
 * Event.
 */
class Event extends ModuleBase
{
    
    /**
     * 用户最后查看时间.
     */
    const REDIS_LAST_VIEW_TIME_USER = 'event_user_last_view_time_';
    
    /**
     * 系统最后一条消息时间.
     */
    const REDIS_SYSTEM_LAST_TIME = 'event_sys_last_time_';
    
    /**
     * 系统最后一条消息(memcache中)时间.
     */
    const MEMCACHE_SYSTEM_LAST_TIME_SUFFIX = 'mc';
    
    /**
     * 回复消息总数.
     */
    const REDIS_HASH_COMMENTS_COUNT = 'koubei_comment_count_';
    
    /**
     * 系统消息总数.
     */
    const REDIS_SYSTEM_USER_COUNT = 'koubei_system_user_count';
    
    /**
     * 用户消息总数.
     */
    const REDIS_USER_EVENT_LIST = 'koubei_user_event_';
    
    /**
     * 得到口碑报告回复消息总数.
     * 
     * @param integer $uid Uid.
     * @param integer $pid 权限组Id.
     * 
     * @return integer
     */
    public function countEventReprotCommentByUid($uid, $pid)
    {
        $type = \Model\Event::EVENT_TYPE_ADD_REPORT_COMMENTS;
        return $this->getCountEventByTypeOfUid($type, $uid, $pid);
    }
    
    /**
     * 得到用户收到的系统消息总数.
     * 
     * @param integer $uid Uid.
     * @param integer $pid 权限组Id privilege_group.
     * 
     * @return integer
     */
    public function countEventSystemByUid($uid, $pid)
    {
        $type = \Model\Event::EVENT_TYPE_ADD_SYSTEM_USER;
        return $this->getCountEventByTypeOfUid($type, $uid, $pid);
    }
    
    /**
     * 通过消息类型统计消息条数.
     * 
     * @param integer $type 消息类型.
     * @param integer $uid  用户Id.
     * @param integer $pid  用户组privilege_group.
     * 
     * @return integer
     */
    private function getCountEventByTypeOfUid($type, $uid, $pid)
    {
        $redis = \Redis\Redis::getStorage();
        $memcache = \Memcache\Pool::instance();
        switch ($type) {
            case \Model\Event::EVENT_TYPE_ADD_SYSTEM_MSG:
            case \Model\Event::EVENT_TYPE_ADD_SYSTEM_USER:
            case \Model\Event::EVENT_TYPE_ADD_SYSTEM_GROUP:
                $userKey = $this->generateKey(self::REDIS_LAST_VIEW_TIME_USER, $uid);
                $lastViewTime = $redis->get($userKey);

                $timeKey = $this->generateKey(self::REDIS_SYSTEM_LAST_TIME, 'GLOBAL');
                $timeKeyMc = $this->generateKey($timeKey, self::MEMCACHE_SYSTEM_LAST_TIME_SUFFIX);
                if (!$lastSysTime = $memcache->get($timeKeyMc)) {
                    $lastSysTime = $redis->get($timeKey);
                }
                if ($lastSysTime <= $lastViewTime) {
                    $timeKey = $this->generateKey(self::REDIS_SYSTEM_LAST_TIME . 'GROUP_', $pid);
                    $timeKeyMc = $this->generateKey($timeKey, self::MEMCACHE_SYSTEM_LAST_TIME_SUFFIX);
                    if (!$lastGroupTime = $memcache->get($timeKeyMc)) {
                        $lastGroupTime = $redis->get($timeKey);
                    }

                    if ($lastGroupTime <= $lastViewTime) {
                        $timeKey = $this->generateKey(self::REDIS_SYSTEM_LAST_TIME . 'USER_', $uid);
                        $timeKeyMc = $this->generateKey($timeKey, self::MEMCACHE_SYSTEM_LAST_TIME_SUFFIX);
                        if (!$lastUserTime = $memcache->get($timeKeyMc)) {
                            $lastUserTime = $redis->get($timeKey);
                        }

                        if ($lastUserTime <= $lastViewTime) {
                            return 0;
                        }
                    }
                }
                return 1;
            case \Model\Event::EVENT_TYPE_ADD_REPORT_COMMENTS:
                $count = $redis->get($this->generateKey(self::REDIS_HASH_COMMENTS_COUNT, $uid));
                break;
            case \Model\Event::EVENT_TYPE_ADD_SYSTEM_USER:
                $count = $redis->get($this->generateKey(self::REDIS_SYSTEM_USER_COUNT, $uid));
                break;
            default:
                $cacheKey_type = $this->generateKey(self::REDIS_USER_EVENT_LIST, $type . "_" . $uid);
                $count = $redis->lSize($cacheKey_type);
                break;
        }
        return $count ? $count : 0;
    }
    
    /**
     * 通过消息类型分页获取用户的消息 返回消息类型对于的rid集合.
     * 
     * @param string  $type         消息类型(all：所有|reports：报告|comments：短评|buy：购买|reviews：评论).
     * @param integer $uid          用户Id.
     * @param integer $page         开始页.
     * @param integer $limit        每页条数.
     * @param integer &$countNumber 数据总数.
     * 
     * @return array
     */
    public function getUserEventsByType($type, $uid, $page, $limit, &$countNumber)
    {
        $typeEvents = array ();
        $cacheKey = '';
        switch ($type) {
            case 'all':
                $cacheKey = $this->generateKey(self::REDIS_USER_EVENT_LIST, $uid);
                break;
            case 'reports':
                $cacheKey = $this->generateKey(self::REDIS_USER_EVENT_LIST, \Model\Event::EVENT_TYPE_ADD_REPORT . '_' . $uid);
                break;
            case 'comments':
                $cacheKey = $this->generateKey(self::REDIS_USER_EVENT_LIST, \Model\Event::EVENT_TYPE_ADD_DEALCOMMENTS . '_' . $uid);
                break;
            case 'buy':
                $cacheKey = $this->generateKey(self::REDIS_USER_EVENT_LIST, \Model\Event::EVENT_TYPE_ADD_PRODUCT . '_' . $uid);
                break;
            case 'reviews':
                $cacheKey = $this->generateKey(self::REDIS_USER_EVENT_LIST, \Model\Event::EVENT_TYPE_ADD_REVIEWS . '_' . $uid);
                break;
        }
        $redis = \Redis\Redis::getStorage();
        $countNumber = $redis->lSize($cacheKey);
        if (!$countNumber) {
            $countNumber = 0;
            return $typeEvents;
        }
        $start = $limit * $page - $limit;
        $end = $limit * $page - 1;
        $eventIds = $redis->lRange($cacheKey, $start, $end);
        
        $events = array();

        foreach ($eventIds as $eventId) {
            $events[] = $this->getEventById($eventId);
        }

        foreach ($events as $event) {
            if (!isset($typeEvents[$event['event_type']])) {
                $typeEvents[$event['event_type']] = array();
            }
            $typeEvents[$event['event_type']][] = $event['rid'];
        }
        return $typeEvents;
    }

}
