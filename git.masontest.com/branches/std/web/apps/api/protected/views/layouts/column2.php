<?php 
    Yii::app()->clientScript->registerCssFile(Yii::app()->clientScript->getCoreScriptUrl().'/css/admin/dashboard.css');

    $colorClasses   = array('red', 'orange', 'yellow', 'green', 'sky-blue', 'blue', 'grey');
    $topMenus       = ProjectParams::$menus;

    $curControll    = Yii::app()->controller->id;
    $curAction      = $this->getAction()->getId();

    $curCA          = Utils::getRootCA($curControll, $curAction, $topMenus);
    $curSubMenu     = isset($topMenus[ $curCA[0] ]['subMenus']) ? $topMenus[ $curCA[0] ]['subMenus'] : array();

    if (!isset($curCA[1])) { $curCA[1] = $curControll . '/' . $curAction; }
?>

<?php $this->beginContent('//layouts/main'); ?>
<div role="navigation" class="navbar navbar-inverse navbar-fixed-top">
    <div class="container-fluid">
        <div class="navbar-header">
            <button data-target=".navbar-collapse" data-toggle="collapse" class="navbar-toggle" type="button"> <span class="sr-only">Toggle navigation</span> <span class="icon-bar"></span> <span class="icon-bar"></span> <span class="icon-bar"></span> </button>
            <a href="#" class="navbar-brand">自由行后台管理</a>
        </div>
        <div class="navbar-collapse collapse">
            <ul class="nav navbar-nav navbar-right">
                <?php foreach($topMenus as $ca => $detail) : ?>
                <li<?php if ($curCA[0] === $ca) : ?> class="active"<?php endif; ?>>
                    <a href="<?php echo url($ca); ?>"><?php echo $detail['name']; ?></a>
                </li>
                <?php endforeach; ?>
            </ul>
            <form class="navbar-form navbar-right">
                <input type="text" placeholder="Search..." class="form-control input-sm">
            </form>
        </div>
    </div>
</div>
<div class="container-fluid">
    <div class="row">
        <div class="col-sm-3 col-md-2 sidebar">
            <?php if (empty($curSubMenu)) : ?>
            <b>暂无子菜单</b>
            <?php else : ?>
            <ul class="nav nav-sidebar">
                <?php $i=0; $colors = count($colorClasses); foreach ($curSubMenu as $ca => $detail) : ?>
                <li<?php if ($curCA[1] === $ca) : ?> class="active"<?php endif; ?>>
                    <a class="<?php echo $colorClasses[ $i++ % $colors ]; ?>-background" href="<?php echo url($ca); ?>"><?php echo $detail['name']; ?></a>
                </li>
                <?php endforeach; ?>
            </ul>
            <?php endif; ?>
            
            <div class="left-bar-bottom">
                <span class="glyphicon glyphicon-user"></span>
                <span class="glyphicon glyphicon-bell"></span>
            </div>
        </div>
        <div class="col-sm-9 col-sm-offset-3 col-md-10 col-md-offset-2 main">
            <?php echo $content; ?>
        </div>
    </div>
</div>
	<?php
        if (0) {
		$this->beginWidget('zii.widgets.CPortlet', array(
			'title'=>'Operations',
		));
		$this->widget('zii.widgets.CMenu', array(
			'items'=>$this->menu,
			'htmlOptions'=>array('class'=>'operations'),
		));
		$this->endWidget();
        }
	?>
<?php $this->endContent(); ?>
