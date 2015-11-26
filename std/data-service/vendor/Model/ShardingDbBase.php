<?php
/**
 * Base classes of models which base on database.
 *
 * @author Huang HaoJie <haojieh@jumei.com>
 */

namespace Model;

/**
 * Base classes of models which base on database.
 */
class ShardingDbBase extends DbBase
{

    public static $db;

    /**
     * Get a instance of DbShardingConnection of the specified connecton name.
     *
     * @param string $rule Sharding Rule is an instance of \Core\Lib\OrderShardingRule.
     *
     * @return \Core\Lib\DbShardingConnection
     */
    public function getDbSharding($rule = null)
    {
        if ( ! static::$db instanceof \Db\ShardingConnection) {
            static::$db = \Db\ShardingConnection::instance();
        }
        if ($rule != null) {
            static::$db->setRule($rule);
        }

        return static::$db;
    }

}
