<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>Content</li>
    <li class="active">Labels</li>
</ol>

<p>
    <span style="float: right">Time: <?php echo $current_time ?></span>
</p>
<p style="clear: both">
    <input id="select_all" type="checkbox"> Select all
    <input id="add-button" type="button" value="Add">
    <input id="delete-button" type="button" value="Delete">
    <select>
        <option>Show New</option>
    </select>

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
            <th>Labels</th>
            <th>Status</th>
            <th>Used times</th>
            <th>Input by</th>
            <th>Input time</th>
        </tr>
    </thead>
    <tbody id="label-container">
        <?php foreach ($list as $one) : ?> 
        <tr id="label-<?php echo $one['label_id']; ?>">
            <td><input name="label_id[]" type="checkbox" value="<?php echo $one['label_id']; ?>" /></td>
            <td><?php echo $one['name']; ?></td>
            <td>New</td>
            <td>0</td>
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

