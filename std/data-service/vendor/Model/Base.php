<?php
/**
 * Base of all models.
 *
 * @author Su Chao<chaos@jumei.com>
 */

namespace Model;

/**
 * Abstract model,included commond methods for data access and manipulations for derived classes.
 * @uses Db\Connection
 */
abstract class Base{
    /**
     *
     * Instances of the derived classes.
     * @var array
     */
    protected static $instances = array();

    /**
     * Get instance of the derived class.
     *
     * @return \Model\Base
     */
    public static function instance()
    {
        $className = get_called_class();
        if (!isset(self::$instances[$className]))
        {
            self::$instances[$className] = new $className;
        }
        return self::$instances[$className];
    }
}
