<?php

class AdminController extends Controller
{

    /**
     * This is the default 'index' action that is invoked
     * when an action is not explicitly requested by users.
     */
    public function actionIndex()
    {
        $this->requiredLogin    = false;

        /*
        $Rpc    = RpcClient_STD_Topic::Instance();
        $data   = $Rpc->getdetail(38);
        */
    
        if (!Yii::app()->user->isGuest)
            $this->redirect( url('admin/main') );
        else
            $this->redirect( url('admin/login') );

        $this->render('index');
    }

    public function actionMain()
    {
        $this->pageTitle  = '管理首页';

        $Rpc = RpcClient_STD_Main::instance();
        $data   = $Rpc->getUserById(1);
        die(var_dump($data));
        $token          = $this->getToken();
        $statDetail     = $Rpc->getAllStat($token, date("Y/m/d"));
        $params['detail']   = $statDetail['data']['detail'];
        $params['columns']  = array('users', 'topics', 'comments');

        $this->jsApp    = array('controller/home', 'Init', array('date' => $statDetail['data']['date']));

        $this->render('main', $params);
    }

    /**
     * This is the action to handle external exceptions.
     */
    public function actionError()
    {
        if($error=Yii::app()->errorHandler->error)
        {
            if(Yii::app()->request->isAjaxRequest)
                echo $error['message'];
            else
                $this->render('error', $error);
        }
    }

    /**
     *Displays the register page
     */
    public function actionRegister(){
        $model=new RegisterForm();
        // if it is ajax validation request
        if(isset($_POST['ajax']) && $_POST['ajax']==='register-form')
        {
            Yii::app()->end();
        }


        if(isset($_POST['RegisterForm']))
        {
            $user = new User();
            $user->attributes = $_POST['RegisterForm'];
            // validate user input and redirect to the previous page if valid
            // Yii::app()->user  means  MyshowWebUser
            if($user->register())
            {
                $this->redirect(Yii::app()->user->returnUrl);
            }	
            $model->addError('email', 'email已经被注册过了，请重新注册');
        }
        $this->pageTitle  = '注册';
        $this->breadcrumbs=array(
            $this->pageTitle
        );

        // display the login form
        $this->render('register',array('model'=>$model));
    }


    /**
     * Displays the login page
     */
    public function actionLogin()
    {
        $this->layout           = '//layouts/column1';
        $this->requiredLogin    = false;

        $model  = new LoginForm;

        // collect user input data
        if(isset($_POST['LoginForm']))
        {
            $model->attributes=$_POST['LoginForm'];
            // validate user input and redirect to the previous page if valid
            // Yii::app()->user  means  MyshowWebUser

            if($model->login())
            {
                $this->redirect(Yii::app()->user->returnUrl);
            }	
            $model->addError('password', '用户名和密码不匹配，请重新输入');

        }

        $this->pageTitle  = '请登录';
        $this->breadcrumbs=array(
            $this->pageTitle
        );

        // display the login form
        $this->render('login',array('model'=>$model));
    }

    /**
     * Logs out the current user and redirect to homepage.
     */
    public function actionLogout()
    {
        Yii::app()->user->logout();
        $this->redirect(Yii::app()->homeUrl);
    }
}
