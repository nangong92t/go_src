<?php
/**
 * 基于Yii Cookie的封装
 * 
 * <pre>
 * 创建一个COOKIE[简洁形式]: Yii::app()->cookie->set('cookie_name', 'cookie_value');
 * 创建一个COOKIE[自定义参数]: Yii::app()->cookie->set('cookie_name', 'cookie_value', array(
 *      'expire' => 0,
 *      'domain' => '',
 *      'path'   => '/',
 *      'secure'   => false,
 *      'httpOnly'   => false,
 * ));<br />
 * 获取一个COOKIE: Yii::app()->cookie->get('cookie_name');
 * 清除一个COOKIE: Yii::app()->cookie->clear('cookie_name');
 * 清除所有COOKIE: Yii::app()->cookie->clear();
 * </pre>
 * 
 * @author TonyXu
 * @package mason.component
 */
class MasonHttpCookie
{
	/**
	 * @var string domain of the cookie
	 */
	public $domain = '';
	/**
	 * @var integer the timestamp at which the cookie expires. This is the server timestamp. Defaults to 0, meaning "until the browser is closed".
	 */
	public $expire = 0;
	/**
	 * @var string the path on the server in which the cookie will be available on. The default is '/'.
	 */
	public $path = '/';
	/**
	 * @var boolean whether cookie should be sent via secure connection
	 */
	public $secure = false;
	/**
	 * @var boolean whether the cookie should be accessible only through the HTTP protocol.
	 * By setting this property to true, the cookie will not be accessible by scripting languages,
	 * such as JavaScript, which can effectly help to reduce identity theft through XSS attacks.
	 * Note, this property is only effective for PHP 5.2.0 or above.
	 */
	public $httpOnly = false;
    
    public function init()
    {
    }
    
    /**
     * 获取一个COOKIE的值
     * @param string $name COOKIE名
     * @return mixed 成功返回字符串,失败返回null
     */
    public function get($name)
    {
        return (isset(Yii::app()->request->cookies[$name])) ? Yii::app()->request->cookies[$name]->value : null;
    }
    
    /**
     * 创建一个COOKIE
     * @param string $name COOKIE名
     * @param string $value COOKIE值
     * @param array $params COOKIE自定义参数
     */
    public function set($name, $value, array $params = array())
    {
        $cookie = new CHttpCookie($name, $value);

        $cookie->expire = isset($params['expire']) ? $params['expire'] : $this->expire;
        if (0 < $cookie->expire) $cookie->expire = time() + $cookie->expire;
        $cookie->path = isset($params['path']) ? $params['path'] : $this->path;
        $cookie->domain = isset($params['domain']) ? $params['domain'] : $this->domain;
        $cookie->secure = isset($params['secure']) ? $params['secure'] : $this->secure;
        $cookie->httpOnly = isset($params['httpOnly']) ? $params['httpOnly'] : $this->httpOnly;
        
        Yii::app()->request->cookies[$name] = $cookie;
    }
    
    /**
     * 清除COOKIE
     * 如果没有给出$name参数,将清除所有
     * @param string $name 清除的COOKIE名
     */
    public function clear($name = '')
    {
        if ($name) unset(Yii::app()->request->cookies[$name]);
        else Yii::app()->request->cookies->clear();
    }
}
