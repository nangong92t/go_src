<?php
/**
 * Rest主控制器类.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-10-05
 */
class RestController extends Controller
{

    /**
     * Rest Action: GET, 通过Get方式获取数据.
     *
     * @param string $token  Token认证字符串.
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     * @param string $params 传入参数，JSON格式字符串.
     *
     * @return mixed. 
     */
    public function actionView($token, $rpc, $action, $params)
    {
        // $this->verify($token);
        $this->RpcProxy($rpc, $action, $params, $token);
    }

    /**
     * Rest Action: GET NO TOKEN, 通过Get方式获取数据.
     *
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     * @param string $params 传入参数，JSON格式字符串.
     *
     * @return mixed. 
     */
    public function actionViewNoToken($rpc, $action, $params)
    {
        $this->RpcProxy($rpc, $action, $params);
    }

    /**
     * Rest Action: CREATE, 通过POST方式获取数据, 这里会封装获取流数据的方法.
     *
     * @param string $token  Token认证字符串.
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     *
     * @return mixed. 
     */
    public function actionCreate($token, $rpc, $action)
    {
        // $this->verify($token);
        $params     = array($token);
        foreach ($_POST as $one) { $params[]    = $one; }

        $this->RpcProxy($rpc, $action, $params);
    }

    /**
     * Rest Action: UPDATE, 通过PUT方式获取数据.
     *
     * @param string $token  Token认证字符串.
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     *
     * @return mixed. 
     */
    public function actionUpdate($token, $rpc, $action)
    {
        $this->verify($token);
        $params     = $_POST;
        $this->RpcProxy($rpc, $action, $params);
    }

    /**
     * Rest Action: DELETE, 通过DELETE方式获取数据.
     *
     * @param string $token  Token认证字符串.
     * @param string $rpc    RPC主控制器.
     * @param string $action RPC动作.
     * @param string $params 传入参数，JSON格式字符串.
     *
     * @return mixed. 
     */
    public function actionDelete($token, $rpc, $action, $params)
    {
        $this->verify($token);
        $this->RpcProxy($rpc, $action, $params);
    }

}
