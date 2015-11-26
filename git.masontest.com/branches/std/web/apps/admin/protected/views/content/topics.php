<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>Content</li>
    <li class="active">Topic</li>
</ol>

<?php $form = $this->beginWidget('CActiveForm', array('action' => $this->createURL('content/topic_op'))); ?>
<p>
    <?php echo CHtml::checkBox('select_all', false, array('id' => 'select_all'))?>
    <?php echo CHtml::label('Select all', '', array('style' => 'font-weight: normal')) ?>
    <?php echo CHtml::button('', array('value' => 'Delete', 'name' => '', 'id'=>'delete-button')); ?>
    <?php echo CHtml::button('', array('value' => 'Delete & Ban', 'name' => '', 'id'=>'delete-ban-button')); ?>
    <?php echo CHtml::button('', array('value' => 'Hide', 'name' => '', 'id'=>'hide-button')); ?>
    <span style="float: right">
        <?php echo $current_from ?>-<?php echo $current_to ?> of <?php echo $total ?>
        Show
        <select>
            <option>1-50</option>
        </select>
    </span>
</p>
<table class="table table-striped table-hover">
    <thead>
        <tr>
            <th></th>
            <th>Topics</th>
            <th>Unwanted words</th>
            <th>Reported times</th>
            <th>Author</th>
            <th>Time</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($topics as $one) : ?>
        <tr id="topic-<?php echo $one['topic_id']; ?>">
            <td><?php echo CHtml::checkBox('topic_id[]', false, array('value' => $one['topic_id'])); ?></td>
            <td><?php echo $one['content']; ?></td>
            <td>0</td>
            <td>1</td>
            <td><?php echo $one['author']; ?></td>
            <td><?php echo date('Y-m-d', $one['created']); ?></td>
        </tr>
        <?php endforeach; ?>
    </tbody>
</table>

<?php
    $this->widget('CLinkPager',array(
            'header'            => '',
            'firstPageLabel'    => '首页',
            'lastPageLabel'     => '末页',
            'prevPageLabel'     => '上一页',
            'nextPageLabel'     => '下一页',
            'pages'             => $pages,
            'maxButtonCount'    => 5,
        )
    );
?>
<?php $this->endWidget(); ?>
