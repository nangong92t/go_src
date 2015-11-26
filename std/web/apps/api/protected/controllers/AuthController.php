<?php
/**
 * Api认证控制器.
 *
 * @author TonyXu<tonycbcd@gmail.com>
 * @date 2014-06-17
 */
class AuthController extends Controller
{

    /**
     * This is the default 'error' action that is invoked
     * when an action is not explicitly requested by users.
     */
    public function actionError()
    {
        $this->render();
    }

    /**
     * Auth Action: GET, 主认证方法.
     *
     * @return mixed. 
     */
    public function actionMain()
    {
        $auth = function() {
            header('WWW-Authenticate: Basic realm="LYAPI"');
            header('HTTP/1.0 401 Unauthorized');
            exit;
        };

        if (!isset($_SERVER['PHP_AUTH_USER'])) {
            $auth();
            exit();
        } else {
            $aUser          = trim($_SERVER['PHP_AUTH_USER']);
            $aPass          = trim($_SERVER['PHP_AUTH_PW']);
            $curIp          = $_SERVER['REMOTE_ADDR'];
            $curUserAgent   = $_SERVER['HTTP_USER_AGENT'];

            if (!$aUser) { $auth(); exit; }

            $OauthRpc       = RpcClient_STD_Oauth::Instance();
            $OauthClient    = $OauthRpc->getClient($aUser);

            if (!$OauthClient || !isset($OauthClient['data']['secret_key']) || $OauthClient['data']['secret_key'] !== $aPass) {
                $auth();
                exit;
            }

            $created        = time();
            $secretKey      = $OauthClient['data']['secret_key'];
            $accessToken    = $this->buildToken(array($curIp, $curUserAgent, $secretKey, $created));
            $refreshToken   = $this->buildToken(array($curIp, $curUserAgent, $secretKey, $created), 'refresh');

            $OauthRpc->setAccessToken($OauthClient['data']['oauth_clients_id'], $accessToken, $refreshToken, $curIp, $created);
            
            $this->render(array('access_token' => $accessToken), 200, 'ok');
        }
    }

}
