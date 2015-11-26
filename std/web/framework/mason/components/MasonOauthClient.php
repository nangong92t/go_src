<?php
/**
 * mason Oauth client, 通过此对象生成前端访问Api时所需的Access Token.
 *
 * @auth TonyXu<tonycbcd@gmail.com>
 * @date 2014-08-22
 */
class MasonOauthClient
{
    
    /**
     * Web端固定APP Key.
     */
    private static $appKey  = 'ec6af911428eb106fa5c7cdb16e7ac310325be1106a98a1b211b79e3';

    /**
     * Web端固定SECRET KEY.
     */
    private static $secretKey   = '6b0fa72598335529830cdfcb6ac3cacd9dac0d90ac7a26f389d84034';

    /**
     * 创建Token规则.
     *
     * @param array  $data   创建必要的数据.
     * @param string $offset 偏移量.
     *
     * @return token string.
     */
    private static  function buildToken(array $data = array(), $offset='') {
        return md5(implode(',', $data) . $offset);
    }

    /**
     * 与Api端一致的access token生成器.
     *
     * @return string Access Token.
     */
    public static function buildAccessToken()
    {
        $curIp          = $_SERVER['REMOTE_ADDR'];
        $curUserAgent   = $_SERVER['HTTP_USER_AGENT'];
        $created        = time();
        $secretKey      = self::$secretKey;
        $expiredIn      = 86400;
        $clientId       = 1;

        $OauthRpc       = RpcClient_User_Oauth::Instance();
        $accessToken    = self::buildToken(array($curIp, $curUserAgent, $secretKey, $created));
        $refreshToken   = self::buildToken(array($curIp, $curUserAgent, $secretKey, $created), 'refresh');

        $OauthRpc->setAccessToken($clientId, $accessToken, $refreshToken, $curIp, $created, array(), $expiredIn);

        return array($accessToken, $expiredIn);
    }

}
