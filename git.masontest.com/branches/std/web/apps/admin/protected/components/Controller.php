<?php
/**
 * Controller is the customized base controller class.
 * All controller classes for this application should extend from this base class.
 */
class Controller extends MasonController
{
    /**
     * 默认:页面布局.
     */
    public $layout  = '//layouts/column2'; 

    /**
     * 是否必须登陆.
     */
    protected $requiredLogin    = true;

    /**
     * 首执行函数.
     */
    public function init()
    {
        parent::init();
    }

    /**
     * 在Action执行后hooker，判定那些Action必须登陆.
     *   获取当前action的方法： $result=$action->getid();
     */
    public function beforeAction($action)
    {
        // 执行是否登陆.
        $curActionId = $action->id;

        if ($this->requiredLogin && $curActionId != 'login' && $curActionId !== 'logout') {
            // api toke
            $token  = $this->getToken();
            
            if (Yii::app()->user->isGuest || !$token) {
                $this->redirect( url('admin/login') );
            }
        }

        return true;
    }
}
