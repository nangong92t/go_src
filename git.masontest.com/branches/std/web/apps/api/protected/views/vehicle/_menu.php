<?php

    $topMenus       = ProjectParams::$menus;
    $menus          = $topMenus['resource/index']['subMenus']['vehicle/index']['subMenus'];

    $curControll    = Yii::app()->controller->id;
    $curAction      = $this->getAction()->getId();
    $curCA          = $curControll . '/' . $curAction;
?>

<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">首页</a></li>
    <li><a href="<?php echo url('resource/index'); ?>">资源管理</a></li>
    <li class="active">车辆管理</li>
</ol>

<ul class="nav nav-tabs">
    <?php foreach ($menus as $ca => $one) : ?>
    <li<?php if ($curCA === $ca) : ?> class="active"<?php endif ?>>
        <a href="<?php echo url($ca); ?>"><?php echo $one['name']; ?></a>
    </li>
    <?php endforeach; ?>
</ul>


