<?php

/**
 * 定义自身的WebUser方法
 *
 */
class MasonWebUser extends CWebUser
{

	public function login($identity,$duration=0)
	{
		$id=$identity->getId();
		$states=$identity->getPersistentStates();
		if($this->beforeLogin($id,$states,false))
		{
			$this->changeIdentity($id,$identity->getName(),$states);

			if($duration>0)
			{
				if($this->allowAutoLogin)
					$this->saveToCookie($duration);
				else
					throw new CException(Yii::t('yii','{class}.allowAutoLogin must be set true in order to use cookie-based authentication.',
						array('{class}'=>get_class($this))));
			}

			$this->afterLogin(false);

            $user  = $identity->getUser();
            $_SESSION['roles']   = $user['user_role_id'];
            $_SESSION['from']     = 'web';
		}
		return !$this->getIsGuest();
	}

	protected function beforeLogout()
	{
        $sid    = $this->getState('sid');
        RpcClient_STD_Account::instance()->logout(array('session'=>$sid));

		return true;
	}

}
