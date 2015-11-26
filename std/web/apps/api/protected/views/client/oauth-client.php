  <?php $this->renderPartial('sub-header', array('active'=>$active)); ?>

  <form class="post" enctype="multipart/form-data" method="post" id="user-checker">
    <?php
      $this->widget('system.myshow.widgets.MyshowGridView', array(
        'id'      => 'cities',
        'htmlOptions'     => array(),
        'itemsCssClass'   => 'bordered-table zebra-striped',
        'setDataIdOnTr'   => true,
        'dataProvider'    => $data,
        'selectableRows'  => 2,
        'pager'           => array('header'=>'','cssFile'=>''),
        'template'        => '{summary} {items}{pager}',
        'columns'         => array(
          array(
            'name'    => 'data',
            'header'  => 'API详细内容',
            'class'   => 'CheckerColumn',
            'columns' => array(
              array(
                'name'    => 'code',
                'header'  => 'API KEY',
                'value'   => '$data->_id',
                'type'    => 'show'
              ),
              array(
                'name'    => 'data',
                'header'  => 'API SECRET',
                'value'   => '$data->client_secret',
                'type'    => 'show'
              ),
              array(
                'name'    => 'data',
                'header'  => '创建时间',
                'value'   => 'Date("Y-m-d H:s:i", $data->created)',
                'type'    => 'show'
              ),
              array(
                'name'    => 'data',
                'header'  => '是否有效',
                'value'   => '$data->active',
                'htmlOptions' => array(
                  'options' => OauthClients::$actives
                ),
                'type'    => 'select'
              ),
              array(
                'name'    => 'data',
                'header'  => '注释',
                'htmlOptions' => array(
                  'tips'  => '最大字数100'
                ),
                'value'   => '$data->desc',
                'type'    => 'textarea'
              ),
              array(
                'name'    => 'data',
                'type'    => 'checkerbutton'
              ) 
            )
          )
        )
      ));
    ?>
    <div class="actions">
      <input type="submit" class="btn primary" value="提交审核结果">
    </div>
  </form>

