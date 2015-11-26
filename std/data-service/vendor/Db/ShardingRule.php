<?php
/**
 * Class ShardingRule.
 *
 * @author Haojie Huang<haojieh@jumei.com>
 */

namespace Db;

/**
 * 拆表拆库规则类.
 */

abstract class ShardingRule
{

    const CFG_PREFIX = 'instance_';

    /**
     * 属性调用返回.
     *
     * @param string $name 属性名.
     * 
     * @return string
     */
    public function __get($name)
    {
        switch ($name) {
            case 'name':
                return $this->getCfgName();
                break;
        }
        return $this->getCfgName();
    }

    /**
     * 获取sharding后的表名.
     * 
     * @param string $table Sharding前的表名.
     * 
     * @return string.
     */
    abstract public function getTableName($table);

    /**
     * 获取sharding后的数据库名.
     * 
     * @param string $db Sharding前的数据库.
     * 
     * @return string.
     */
    abstract public function getDbName($db);

    /**
     * 获取sharding后的服务器配置项名/服务器名.
     *
     * @return string.
     */
    abstract public function getCfgName();
    
    /**
     * 获取更小的唯一标识名.
     * 
     * @return string.
     */
    public function getAtomName()
    {
        return $this->getTableName($this->getCfgName());
    }

}
