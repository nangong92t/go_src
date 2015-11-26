<?php

/**
 * UserIdentity represents the data needed to identity a user.
 * It contains the authentication method that checks if the provided
 * data can identity the user.
 * @package nlsp.component
 */
class UserIdentity extends CUserIdentity
{
    const ERROR_ROLE_NO_ACCESS   = 99;
    protected $id;

    protected $accessRoleId;

	/**
	 * Constructor.
     *
	 * @param string  $username     username.
	 * @param string  $password     password.
     * @param integer $accessFoleId Role Id.
	 */
	public function __construct($username,$password, $accessRoleId = 0)
	{
		$this->username=$username;
		$this->password=$password;
        $this->accessRoleId = $accessRoleId;
	}

    /**
     * Authenticates a user.
     * The example implementation makes sure if the username and password
     * are both 'demo'.
     * In practical applications, this should be changed to authenticate
     * against some persistent user identity storage (e.g. database).
     * @return boolean whether authentication succeeds.
     */
    public function authenticate()
    {
        $user = RpcClient_STD_Account::instance()->login($this->username, $this->password);

        if (!$user && !isset($user['data'])) {
            $this->errorCode = self::ERROR_USERNAME_INVALID;
        } else if ($user['error']) {
            $this->errorCode = self::ERROR_PASSWORD_INVALID;
        } else if ($this->accessRoleId && (int)$user['data']['user']['user_role_id'] !== $this->accessRoleId) {
            $this->errorCode = self::ERROR_ROLE_NO_ACCESS;
        } else {
            $this->id   = $user['data']['user']['user_id'];
            $this->setState('user',  $user['data']['user']);
            $this->setState('token', $user['data']['session']);
            $this->errorCode = self::ERROR_NONE;
        }

        return $this->errorCode;
    }

    public function getId()
    {
        return $this->id;
    }

    public function getUser()
    {
        return $this->getState('user');
    }

    public function getNickName()
    {
        return $this->getState('nickName') ? $this->getState('nickName') : $this->_nickName;
    }
}
