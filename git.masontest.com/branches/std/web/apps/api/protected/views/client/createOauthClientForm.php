<?php

return array(
  'title'   => '创建新API',
  'enctype' => 'multipart/form-data',
  'class'   => 'create-form', 
  'elements'=>array(
    'desc'=>array(
      'type'=>'textarea',
      'maxlength'=>100,
      'value' => 'test',
      'hint'=>'最大字数100'
    ),
    'active'  => array(
      'type'=> 'dropdownlist',
      'items' => OauthClients::$actives,
      'prompt'=>'请选择:',
    )
  ),

  'buttons'=>array(
    'save'=>array(
      'class'=> 'btn btn-submit primary',
      'type'  =>'submit',
      'label'=>'创建',
    ),
  ),
);
