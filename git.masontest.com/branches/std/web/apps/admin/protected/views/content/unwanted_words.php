<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>Content</li>
    <li class="active">Unwanted words</li>
</ol>
<p>
    <span style="float: right">Time: <?php echo $current_time ?></span>
</p>
<p style="clear: both">
    <input type="checkbox" id="select_all"> Select all
    <input id="add-button" type="button" value="Add">
    <input id="delete-button" type="button" value="Delete">

    <span style="float: right">
        <?php echo $current_from ?>-<?php echo $current_to ?> of <?php echo $total ?>
        Show
        <select>
            <option>1-50</option>
        </select>
        <input type="text" placeholder="Search">
    </span>
</p>
<table class="table table-striped table-hover">
    <thead>
        <tr>
            <th></th>
            <th>Unwanted words</th>
            <th>used times</th>
            <th>Added by</th>
            <th>Added time</th>
        </tr>
    </thead>
    <tbody id="unwanted_words-container">
        <?php foreach ($list as $one) : ?>
        <tr id="unwanted_words-<?php echo $one['word_id']; ?>">
            <td><?php echo CHtml::checkBox('unwanted_words_id[]', false, array('value' => $one['word_id'])); ?></td>
            <td><?php echo $one['word']; ?></td>
            <td>1</td>
            <td><?php echo Yii::app()->user->getName(); ?></td>
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
