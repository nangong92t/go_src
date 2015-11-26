<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>Content</li>
    <li class="active">Comment</li>
</ol>


<p>
    <input id="select_all" type="checkbox"> Select all
    <input id="delete-button" type="button" value="Delete">
    <input id="delete-ban-button" type="button" value="Delete & Ban">
    <input id="hide-button" type="button" value="Hide">

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
            <th>Comment</th>
            <th>Unwanted words</th>
            <th>Reported times</th>
            <th>Comment by</th>
            <th>Time</th>
        </tr>
    </thead>
    <tbody>
        <?php foreach ($comments as $one) :?>
        <tr id="comment-<?php echo $one['comment_id']; ?>">
            <td><input type="checkbox" name="comment_id[]" value="<?php echo $one['comment_id']; ?>"</td>
            <td><?php echo $one['content']; ?>
            <td>0</td>
            <td>0</td>
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
