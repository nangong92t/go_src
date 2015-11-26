<?php
/**
 * Mobile认证控制器.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-07-22
 */
class MobileController extends Controller
{
    /**
     * 重载统一数据输出格式.
     */
    protected $response = array(
        'status'        => 0,           // 0:表示逻辑错误，app不应执行下一步操作;1:表示请求成功，app可以继续执行操作
        'error_id'      => 0,           // 系统的统一错误号，未出现错误置0。注意error_id一定要细化，因为error_id在某些情况下可能影响到客户端行为。比如用户输入时用户名或者密码错误，无论是密码还是用户名info都返回该信息，但是error_id需要注明到底是用户名错误还是密码错误，因为用户名错误的时候会返回到用户名录入的页面，但是密码错误页面就仍然停留在密码输入页。error_info和error_id一一对应.
        'error_info'    => '',          // 出错说明，未出现错误留空.
        'info'          => '',          // 请求的结果文本，例如用户登陆成功时info为登录成功。该字段主要是辅助客户端显示请求结果.
        'data'          => array()      // 请求结果的值.
    );

    /**
     *  重载Render方法.
     *
     *  @param mixed   $data      输出数据.
     *  @param integer $status    如$response.status中说明.
     *  @param integer $errorId   如$response.error_id中说明.
     *  @param string  $errorInfo 如$response.error_info中说明.
     *  @param string  $info      如$response.info中说明.
     *
     *  @return void.
     */
    public function render($data = null, $status = 1, $errorId = 0, $errorInfo = '', $info = '')
    {
        $this->response['data']     = $data;
        if ($status) { $this->response['status']        = $status; }
        if ($errorId) { $this->response['error_id']     = $errorId; }
        if ($errorInfo) { $this->response['error_info'] = $errorInfo; }
        if ($info) { $this->response['info']            = $info; }

        $httpStatus     = $status == 1 ? 200 : ($this->response['error_id'] > 400 ? $this->response['error_id'] : 403);
        $status_header  = 'HTTP/1.1 ' . $httpStatus . ' ' . $this->_getStatusCodeMessage($httpStatus);

        header($status_header);
//        header('Content-type: text/json');
        header('Access-Control-Allow-Origin: *');
        header('Access-Control-Allow-Method: GET,POST,PUT,DELETE,OPTIONS');
        echo CJSON::encode($this->response);
        Yii::app()->end();
    }

    /**
     * This is the default 'error' action that is invoked
     * when an action is not explicitly requested by users.
     */
    public function actionError()
    {
        $defaultErrorCode   = 401;
        $this->render('', 0, $defaultErrorCode, $this->_getErrorIdMessage($defaultErrorCode));
    }

    /**
     * Mobile端的接口的认证与用户主持.
     *
     * @return void.
     */
    public function actionMain()
    {
        if (isset($_POST['authorize'])) {
            $this->register();
        } elseif (isset($_POST['rewrite'])) {
            $this->login();
        }

        $this->actionError();
    }

    /**
     * Mobile端获取用户信息，当rewrite=1时, 并附带授权.
     */
    public function login()
    {
        $req        = Yii::app()->request;

        $user       = $req->getParam('user');       // 用户名
        $password   = $req->getParam('password');   // 密码hash结果
        $rewrite    = $req->getParam('rewrite');    // 是否重新授权:0表示不重新授权，仅作为登录使用，1表示重新授权.

        $UserRpc    = RpcClient_User_Main::Instance();
        $userInfo   = $UserRpc->getUserByName($user);

        if (!$userInfo['data'] || !isset($userInfo['data']['password'])) {
            $this->render('', 1, 2, $this->_getErrorIdMessage(2));
        } elseif ($userInfo['data']['password'] != $password) {
            $this->render('', 1, 3, $this->_getErrorIdMessage(3));
        } 

        $result     = $userInfo['data'];
        if (1 || $rewrite) {
            // 如果rewrite为1,则无论之前是否授权过，都重新授权，并且删除之前的token.
            //
            $accessData = $this->authorize($userInfo['data']);
            $result['access_token'] = $accessData['accessToken'];
            $result['expire_time']  = $accessData['expired'];
        } else {
            // 在rewrite为0时，如果用户已经授权过，并且未过期，则直接返回上一次授权结果，并且返回上一次调用时间和区域；如果未授权则授权;

            // TODO: 
            $OauthRpc       = RpcClient_User_Oauth::Instance();
            $accessToken    = $OauthRpc->getAccessTokenByCondition($condition);   

        }

        // TODO: note login info.

        $this->render($result);
    }

    /**
     * Mobile端用户注册，当authorize=1时,并附带授权.
     *
     * @return void.
     */
    public function register()
    {
        $req        = Yii::app()->request;

        $user       = $req->getParam('user');       // 用户名
        $password   = $req->getParam('password');   // 密码hash结果
        $email      = $req->getParam('email');      // 电子邮箱
        $mobile     = $req->getParam('mobile');     // 手机号码
        $authorize  = $req->getParam('authorize');  // 是否同时授权:1授权，0只注册

        $result     = array();

        $UserRpc    = RpcClient_User_Main::Instance();
        $response   = $UserRpc->addNewUser($user, $password, $email, $mobile);
        if ($response['error']) {
            $this->render('', 0, 1, $response['error']); 
        }
        $result     = $response['data'];

        if ($authorize) {
            $accessData = $this->authorize($response['data']);
            $result['access_token'] = $accessData['accessToken'];
            $result['expire_time']  = $accessData['expired'];
        }

        $this->render($result);
    }

    /**
     * 接口授权.
     *       
     * @param array $userInfo 当前授权用户信息.
     *
     * @return string Access token. 
     */
    public function authorize(array $userInfo)
    {
        $mobileClient   = Yii::app()->params['mobileClient'];

        $OauthRpc       = RpcClient_User_Oauth::Instance();

        $curIp          = $_SERVER['REMOTE_ADDR'];
        $curUserAgent   = isset($_SERVER['HTTP_USER_AGENT']) ? $_SERVER['HTTP_USER_AGENT'] : 'mobile';

        $created        = time();
        $secretKey      = $mobileClient['SecretKey'];
        $accessToken    = $this->buildToken(array($curIp, $curUserAgent, $secretKey, $created));
        $refreshToken   = $this->buildToken(array($curIp, $curUserAgent, $secretKey, $created), 'refresh');

        $accessData     = $OauthRpc->setAccessToken($mobileClient['ClientId'], $accessToken, $refreshToken, $curIp, $created, $userInfo);
        $accessData     = $accessData['data'];
        $accessData['accessToken']  = $accessToken;

        return $accessData;
    }

    /**
     * 重载 Token 验证流程.
     *
     * @param string $token Token认证字符串.
     *
     * @return http.
     */
    public function verify($token)
    {
        $result = $this->Authentication($token);
        if ($result['status'] != 1) {
            $this->render('', $result['status'], $result['error_id'], $result['error_info']);
        }
    }

    /**
     * 重载 Token 验证.
     *
     * @param string $token Token认证字符串.
     *
     * @return $this->response.
     */
    public function Authentication($token)
    {
        if ($token) {
            $OauthRpc       = RpcClient_User_Oauth::Instance();
            $accessData     = $OauthRpc->getAccessToken($token);

            if (!$accessData) return $this->response;

            $accessData     = $accessData['data'];

            // 检测过期与否
            if (($accessData['created'] + $accessData['expires_in']) < time()) {
                $this->response['status']       = 0;
                $this->response['error_id']     = 403;
                $this->response['error_info']   = $this->_getErrorIdMessage(403);
                return $this->response;
            }

            $this->response['status']       = 1;
            $this->response['error_id']     = 0;
            $this->response['error_info']   = '';
        }

        return $this->response;
    }

    /**
     * 根据error id获取错误信息.
     *
     * @param integer $errorId 错误Id.
     *
     * @return string error info.
     */
    protected function _getErrorIdMessage($errorId)
    {
        $errors = Array(
            2   => '无此用户', 
            3   => '用户名或密码有误', 
            401 => '未授权',
            403 => '已过期',
        );

        return (isset($errors[$errorId])) ? $errors[$errorId] : '';
    }

    /**
     * Mobile Api 获取数据接口.
     *
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     * @param string $params 传入参数，JSON格式字符串.
     *
     * @return mixed. 
     */
    public function actionProxy($rpc, $action)
    {
        $req    = Yii::app()->request;
        $token  = $req->getParam('access_token');
        $params = $_POST;
        unset($params['access_token']);

        /** 为车载接口暂时打开 token **/ 
        $specialAction = array(
            'setVehicleTravelingTrack',
            'getVehicleRunningList',
            'getVehicleRunningGps'
        );

        if (!in_array($action, $specialAction)) {
            $this->verify($token);
        }

        $this->RpcProxy($rpc, $action, $params);
    }

}
