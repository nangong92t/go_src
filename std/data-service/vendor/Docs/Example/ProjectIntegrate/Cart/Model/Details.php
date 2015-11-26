<?php
namespace Model;

/**
 *
 * @author Su Chaos<suchaoabc@163.com>
 * @uses \Memcache\Pool
 */
class Details extends CartBase{

    /**
     *
     * @return self
     */
    public static function instance()
    {
        return parent::instance();
    }
    public function get($cartId)
    {
        return array();
    }

    public function updateItemCount($cartId, $productId, $count)
    {

        return true;
    }

    public function testCache()
    {
        $mc = \Memcache\Pool::instance();
        $cacheData = $mc->get('tt');
        $mc->set('tt', 123, 3);
        var_dump($cacheData);
    }

    public function testDbConnection()
    {
        $data = \Db\Connection::instance()->read('jumei')
                ->select('order_id, quantity, total_price, uid, timestamp')
                ->from('tuanmei_orders')
                ->where(array('order_id'=>246807))
                ->queryAll();
        var_dump($data);
        return $data;
    }

    public function testLog()
    {
       return \Log\Handler::instance('inventoryLog')
        ->log(array('CartTest', 'sku889911', '1363417830', 'TestService', array('id'=>1, 'user'=>'dev'), array('time'=>1363417830, 'ok'=>'yes'), 1, 33));
    }

    public function testRedis()
    {
        $redis =  \Redis\RedisMultiCache::getInstance('fav');
        $redis->set('a', 23);
        var_dump($redis->get('a'));
    }
}