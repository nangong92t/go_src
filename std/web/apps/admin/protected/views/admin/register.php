<?php
/* @var $this SiteController */
/* @var $model LoginForm */
/* @var $form CActiveForm  */

?>

<div class="mod login">

<h1></h1>

<div class="form">
<?php $form=$this->beginWidget('CActiveForm', array(
	'id'=>'form-register',
	'enableClientValidation'=>true,
	'clientOptions'=>array(
		'validateOnSubmit'=>true,
	),
)); ?>

  <fieldset>
    <legend>欢迎注册</legend>
    <p class="note"><span class="required">*</span>号为必填信息.</p>

    <div class="control-group">
      <?php echo $form->labelEx($model,'email'); ?>
      <?php echo $form->textField($model,'email'); ?>
      <?php echo $form->error($model,'email'); ?>
    </div>

	<div class="control-group">
      <?php echo $form->labelEx($model,'password'); ?>
      <?php echo $form->passwordField($model,'password'); ?>
      <?php echo $form->error($model,'password'); ?>
    </div>
    
    <div class="control-group">
      <?php echo $form->labelEx($model,'repassword'); ?>
      <?php echo $form->passwordField($model,'repassword'); ?>
      <?php echo $form->error($model,'repassword'); ?>
    </div>

    <div class="control-group">
      <?php echo $form->labelEx($model,'nickname'); ?>
      <?php echo $form->textField($model,'nickname'); ?>
      <?php echo $form->error($model,'nickname'); ?>
      <p class="hint">
        
      </p>
    </div>

    <div class="control-group buttons">
      <?php echo CHtml::submitButton('Register'); ?>
    </div>
  </fieldset>
<?php $this->endWidget(); ?>
</div><!-- form -->
