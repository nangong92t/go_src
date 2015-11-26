<?php
Yii::app()->clientScript->registerCssFile(Yii::app()->clientScript->getCoreScriptUrl().'/css/signin.css');
?>
<form role="form" action="post" action="/admin/login" class="form-signin">
    <h2 class="form-signin-heading">请登陆</h2>
    <input type="email" autofocus="" required="" placeholder="用户名" class="form-control">
    <input type="password" required="" placeholder="密码" class="form-control">
    <label class="checkbox">
        <input type="checkbox" value="remember-me"> 记住我
    </label>
    <button type="submit" class="btn btn-lg btn-primary btn-block">登陆</button>
</form>
