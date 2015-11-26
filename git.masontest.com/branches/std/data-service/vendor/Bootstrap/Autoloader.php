<?php
/**
 * Autoloader.
 *
 * @author Su Chao<suchaoabc@163.com>
 */

namespace Bootstrap;

/**
 * Jumei框架体系中的类库自动加载类.
 */
class Autoloader{
    protected static $sysRoot = array();
    protected static $instance;
    protected function __construct()
    {
        static::$sysRoot[] = realpath(__DIR__.'/../..').DIRECTORY_SEPARATOR.'vendor'.DIRECTORY_SEPARATOR;
    }

    /**
     *
     * @return self
     */
    public static function instance()
    {
        if(!static::$instance)
        {
            static::$instance = new static;
        }
        return static::$instance;
    }

    /**
     * 添加根目录. 默将使用Autoloader目录所在的上级目录为根目录。
     *
     * @param string $path
     * @return self
     */
    public function addRoot($path)
    {
        static::$sysRoot[] = $path;
        return $this;
    }

    /**
     * 按命名空间自动加载相应的类.
     * @param string $name 命名空间及类名
     */
    public function loadByNamespace($name)
    {
        $classPath = str_replace('\\', DIRECTORY_SEPARATOR ,$name);

        foreach(static::$sysRoot as $k => $root)
        {
            $classFile = $root.$classPath.'.php';
            if(!is_file($classFile))
            {
                if($k === 0)
                {
                    $classFile = $root.'Vendor'.DIRECTORY_SEPARATOR.$classPath.'.php';
                    if(is_file($classFile))
                    {
                        require($classFile);
                        return true;
                    }
                }
            }
            else
            {
                require($classFile);
                return true;
            }
        }

        return false;
    }

    /**
     *
     * @return \Bootstrap\Autoloader
     */
    public function init()
    {
        spl_autoload_register(array($this, 'loadByNamespace'));

        \PHPClient\JMTextRpcClient::config((array)new \Config\PHPClient());
        return $this;
    }
}
