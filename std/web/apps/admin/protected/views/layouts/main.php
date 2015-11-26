<?php
    $userId = Yii::app()->user->getId();
 ?>
<!DOCTYPE html>
<html>
  <head>
    <title><?php echo CHtml::encode($this->pageTitle); ?></title>
    <meta charset="utf-8">
    <script type="text/javascript" src="<?php echo Yii::app()->clientScript->getCoreScriptUrl(); ?>/js/lib/jquery-1.7.2.min.js"></script>
    <script type="text/javascript" src="<?php echo Yii::app()->clientScript->getCoreScriptUrl(); ?>/js/sea.js?v2.1.1"></script>
    <script>var staticUri='<?php echo Yii::app()->clientScript->getCoreScriptUrl(); ?>',token='<?php echo Yii::app()->user->getState("token"); ?>',customApp=[], apiUrl='http://<?php echo API_DOMAIN; ?>', userId = <?php echo $userId ? $userId : 0; ?>, initInfo='<?php echo $this->info; ?>';</script>
    <?php if ($this->jsApp) : ?>
      <?php if (!is_array($this->jsApp)) { $this->jsApp = array($this->jsApp, 'index', array()); } ?>
    <script>customApp = ['<?php echo $this->jsApp[0]; ?>', '<?php echo $this->jsApp[1]; ?>', <?php echo json_encode($this->jsApp[2]); ?>];</script>
    <?php endif; ?>
    <script src="<?php echo staticurl(); ?>/js/init.js"></script>
    <link rel="stylesheet" href="<?php echo staticurl() ?>/css/admin/style.css">
  </head>
  <body class="admin">
    <?php echo $content; ?>
  </body>
</html>

