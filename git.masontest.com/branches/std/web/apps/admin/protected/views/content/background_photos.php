<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>Content</li>
    <li class="active">Background photos</li>
</ol>

<p>
    <span style="float: right">Time: <?php echo $current_time ?></span>
</p>
<p style="clear: both">
    <input id="select_all" type="checkbox"> Select all
    <span id="spanButtonPlaceholder"></span>
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
    </span>
</p>
<table class="table table-striped table-hover">
    <thead>
        <tr>
            <td></td>
            <td>Photos</td>
            <td>Status</td>
            <td>Used times</td>
            <td>Used ratio</td>
            <td>Input by</td>
            <td>Input time</td>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($list as $one) : ?>
        <tr id="background-<?php echo $one['attach_id']; ?>">
            <td><input name="background_id[]" type="checkbox" value="<?php echo $one['attach_id']; ?>"></td>
            <td><img src="<?php echo staticurl() . '/uploads/' . $one['savepath'].'/'.$one['savename'].$one['extention']; ?>" width="200" /></td>
            <td>New</td>
            <td>none</td>
            <td>none</td>
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

