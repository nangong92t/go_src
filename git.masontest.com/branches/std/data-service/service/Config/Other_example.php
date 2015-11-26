<?php
/**
 * 通用配置.
 * 
 * @author yongf <yongf@jumei.com>
 */

namespace Config;

/**
 *通用配置文件.
 */
class Config
{
    /* CDN的URL路径，需要和jumei_web的保持一致 */
    public $CDNBaseURL = array(
        'Product' => 'http://p0.jmstatic.com/product/',
        'IcProduct' => 'http://p0.jmstatic.com/ic/product/',
        'DealProduct' => 'http://p0.jmstatic.com/deal_product/',
        'DealContent' => 'http://p0.jmstatic.com/deal_content/',
        'Mobile' => 'http://p0.jmstatic.com/mobile/',
        'Brand' => 'http://p0.jmstatic.com/brand/logo',
        'JiandingBrand' => 'http://p0.jmstatic.com/ic/brand/'
    )

}
