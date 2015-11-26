<?php
/**
 * SphinxClient的配置文件.
 * 
 * @author XuRongYi <rongyix@jumei.com>
 */

namespace Config;

if (!defined('SPHINX_AGENT_ENABLE')) {
    define('SPHINX_AGENT_ENABLE', true);
}

/**
 * SphinxClient的配置文件.
 */
class SphinxClient
{
    /**
     * Configs of Redis.
     * @var array
     */
    public $default = array(
        'host' => '192.168.25.11',
        'port' => 9312,
        'agents' => array(
                        'reports' => 'agent_reports',
                        'deal_comments' => array('index' => 'agent_deal_comments', 'sharding' => true, 'sharding_mod' => 12),
                        'products' => 'agent_products',
                        'product_functions' => 'agent_product_functions',
                        'keywords' => 'sphinx_keywords',
        ),
        'options' => array(
                        'ConnectTimeout' => 1,
                        'MaxQueryTime' => 2000,
        ),
    );
}
