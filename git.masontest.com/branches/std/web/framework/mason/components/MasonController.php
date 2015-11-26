<?php
/**
 * Controller is the customized base controller class.
 * All controller classes for this application should extend from this base class.
 */
class MasonController extends CController
{

    /**
     * Cookie中提示消息的key.
     */
    const COOKIE_MESSAGE_KEY = '_message';

    /**
     * 是否默认就获取Access Token.
     */
    public $isGetAccessToken   = true;

	/**
	 * @var string the default layout for the controller view. Defaults to '//layouts/column1',
	 * meaning using a single column layout. See 'protected/views/layouts/column1.php'.
	 */
	public $layout='//layouts/column1';

	/**
	 * @var array context menu items. This property will be assigned to {@link CMenu::items}.
	 */
	public $menu=array();

	/**
	 * @var array the breadcrumbs of the current page. The value of this property will
	 * be assigned to {@link CBreadcrumbs::links}. Please refer to {@link CBreadcrumbs::links}
	 * for more details on how to specify this property.
	 */
	public $breadcrumbs=array();

    /**
     * ly project custom js app proto.
     */
    public $jsApp=array();

    /**
     * ly project custom info.
     */
    public $info    = '';

    /**
     * 首执行函数.
     *
     * return void.
     */
    public function init()
    {
        date_default_timezone_set("PRC");

        /* init Rpc client */
        TextRpcClient::config(Yii::app()->params['rpc']);

        /* to check access token expired. */
        // if ($this->isGetAccessToken) {
        //    $this->checkAccessToken();
        // }
        // $this->setInfo();
    }

    /**
     * 带提示消息的跳转并退出.
     *
     * @param $url Url.
     * @param $msg 提示信息.
     *
     * @return void
     */
    protected function redirectMsg($url, $msg)
    {
        $key = self::COOKIE_MESSAGE_KEY;
        setcookie($key, $msg, null, '/');
        /*
        $and = strpos($url, '?') ? '&' : '?';
        $pattern = '/[&\?]' . self::URL_MESSAGE_KEY . '=.*&?/U';
        $url = preg_replace($pattern, '', $url);
        $url .= $and . self::URL_MESSAGE_KEY . '=' . urlencode($msg);
        */
        $this->redirect($url);
    }

    /**
     * 检测当前是否有Access token， 如果没有则生成.
     *
     * @return void.
     */
    protected function checkAccessToken()
    {
        $curTime    = time();
        $isExpired  = true;

        $yiiSession = Yii::app()->session;
        if ($yiiSession['token_expired'] && ($yiiSession['token_expired'] > $curTime)) {
            $isExpired  = false;
        }
        if ($this->isGetAccessToken && $isExpired) {
            $tokenData  = MasonOauthClient::buildAccessToken();
            Yii::app()->session['token']            = $tokenData[0];
            Yii::app()->session['token_expired']    = $curTime + $tokenData[1];
        }
    }

    /**
     * 获取当前后端数据层token.
     *
     * @return token string.
     */
    protected function getToken()
    {
        return Yii::app()->user->getState('token');
    }

    private function setInfo()
    {
        //$msg = !empty($_GET[self::URL_MESSAGE_KEY]) ? urldecode($_GET[self::URL_MESSAGE_KEY]) : '';
        $key = self::COOKIE_MESSAGE_KEY;
        $msg = '';
        if (!empty($_COOKIE[$key])) {
            $msg = $_COOKIE[$key];
            unset($_COOKIE[$key]);
            setcookie($key, null, -1, '/');
        }
        $this->info = $msg;
    }

 }
