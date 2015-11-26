<?php 
    Yii::app()->clientScript->registerCssFile(Yii::app()->clientScript->getCoreScriptUrl().'/css/admin/dashboard.css');

    $colorClasses   = array('grey');
    $topMenus       = ProjectParams::$menus;

    $curControll    = Yii::app()->controller->id;
    $curAction      = $this->getAction()->getId();

    $curCA          = Utils::getRootCA($curControll, $curAction, $topMenus);

    if (!isset($curCA[0])) { $curCA[0] = 'admin/main'; }
    if (!isset($curCA[1])) { $curCA[1] = $curControll . '/' . $curAction; }
    
?>

<?php $this->beginContent('//layouts/main'); ?>
<div role="navigation" class="navbar navbar-inverse navbar-fixed-top">
    <div class="container-fluid">
        <div class="navbar-header">
            <form class="navbar-form navbar-right">
                <div class="form-group has-success has-feedback">
                    <input type="text" placeholder="Search..." class="form-control input-sm">
                    <span class="glyphicon glyphicon-search form-control-feedback"></span>
                </div>
            </form>
        </div>
        <div class="navbar-collapse collapse">
            <ul class="nav navbar-nav navbar-right">
            </ul>
        </div>
    </div>
</div>
<div class="container-fluid">
    <div class="row">
        <div class="col-sm-3 col-md-2 sidebar">
            <ul class="nav nav-sidebar">
                <?php $i=0; $colors = count($colorClasses); foreach ($topMenus as $ca => $detail) : ?>
                <li class="<?php echo $colorClasses[ $i++ % $colors ]; ?>-background <?php if ($curCA[0] === $ca) : ?> active<?php endif; ?>">
                    <a class="main-menu-a" href="<?php echo url($ca); ?>"><?php echo $detail['name']; ?></a>
                    <?php 
                        $curSubMenu     = isset($topMenus[ $ca ]['subMenus']) ? $topMenus[ $ca ]['subMenus'] : array();
                    ?>
                    <?php if (!empty($curSubMenu)) : ?>
                    <ul class="ly-dropdown-menu clearfix">
                        <?php foreach ($curSubMenu as $ca2 => $detail2) : ?>
                        <li<?php if ($curCA[1] === $ca2) : ?> class="active"<?php endif; ?>>
                            <a href="<?php echo url($ca2); ?>"><?php echo $detail2['name']; ?></a>
                        </li>
                        <?php endforeach; ?>
                    </ul>
                    <?php endif; ?>
                </li>
                <?php endforeach; ?>
            </ul>
            
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
