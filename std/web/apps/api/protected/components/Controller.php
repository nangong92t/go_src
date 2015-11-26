<?php
/**
 * Controller is the customized base controller class.
 * All controller classes for this application should extend from this base class.
 */
class Controller extends MasonController
{

    // 默认不生成 access token.
    public $isGetAccessToken   = false;

    /**
     * 统一数据输出格式.
     */
    protected $response = array(
        'status'  => 401,
        'mesg'  => '未授权',
        'body'  => array()
    );

    /**
     * Token 验证.
     *
     * @param string $token Token认证字符串.
     *
     * @return $this->response.
     */
    public function Authentication($token)
    {
        if ($token) {
            $OauthRpc       = RpcClient_STD_Oauth::Instance();
            $accessData     = $OauthRpc->getAccessToken($token);

            if (!$accessData) return $this->response;

            $accessData     = $accessData['data'];

            // 检测过期与否
            if (($accessData['created'] + $accessData['expires_in']) < time()) {
                $this->response['status']   = 403;
                $this->response['mesg']     = '已过期';
                return $this->response;
            }

            $this->response['status']   = 200;
            $this->response['mesg']     = '';
        }

        return $this->response;
    }

    /**
     * Token 验证流程.
     *
     * @param string $token Token认证字符串.
     *
     * @return http.
     */
    public function verify($token)
    {
        $result = $this->Authentication($token);
        if ($result['status'] != 200) {
            $this->render('', $result['status'], $result['mesg']);
        }
    }

    /**
     *  重载Render方法.
     *
     *  @param mixed $data 输出数据.
     */
    public function render($body = null, $status = 0, $mesg = '')
    {
        if ($body) { $this->response['body']     = $body; }
        if ($status) { $this->response['status']   = $status; }
        if ($mesg) { $this->response['mesg']     = $mesg; }

        $status_header = 'HTTP/1.1 ' . $this->response['status'] . ' ' . $this->_getStatusCodeMessage($this->response['status']);

        header($status_header);

        if ($this->response['mesg']) {
            $mesg   = $this->response['mesg'];
            foreach ($this->response as $k => $v) {
                unset($this->response[$k]);
            }
            $this->response['errmsg']   = $mesg;
            $this->response['errcode']  = 1;
        } else {            // 按照mobile端工程师要求做的.
            $this->response     = $this->response['body'];
        }
        echo CJSON::encode($this->response);
        Yii::app()->end();
    }

    /**
     * RPC核心代理方法.
     *
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     * @param mixed  $params 传入参数，JSON格式字符串.
     * @param string $token  用户Session token.
     *
     * @return mixed. 
     */
    public function RpcProxy($rpc, $action, $params, $token = null)
    {
        $className  = 'RpcClient_STD_' . $rpc;
        $data       = '';

        if (!is_array($params)) {
            $params     = json_decode($params);
            if ($token) array_unshift($params, $token);
        }

        try {
            $rpcObject  = $className::instance();
            $data       = call_user_func_array(array($rpcObject, $action), $params);
            if ($data['mesg'] != '') {
                $this->response['status']   = 403;
                $this->response['mesg']     = $data['mesg'];
            } else {
                if (isset($data['data'])) { $this->response['data']   = $data['data']; }
                $this->response['status']   = 200;
                $this->response['mesg']     = '';
                $this->response['body']     = $data['data'];
            }
        } catch (Exception $e) {
            $this->response['status']   = 400;
            $this->response['mesg']     = $e->getMessage();
        }

        $this->render();
    }

    /**
     * 获取status对应 http状态消息.
     *
     * @param integer $status Http状态码.
     *
     * @return string.
     */
    protected function _getStatusCodeMessage($status)
    {
        $codes = Array(
            200 => 'OK',
            400 => 'Bad Request',
            401 => 'Unauthorized',
            402 => 'Payment Required',
            403 => 'Forbidden',
            404 => 'Not Found',
            500 => 'Internal Server Error',
            501 => 'Not Implemented',
        );
        return (isset($codes[$status])) ? $codes[$status] : '';
    }

    /**
     * 统一token生成方法.
     *
     * @param array  $data   数据源.
     * @param string $offset 偏移字符串.
     *
     * @return string.
     */
    protected function buildToken($data = array(), $offset='')
    {
        return md5(implode(',', $data) . $offset);
    }

}
