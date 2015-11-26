<?php 
    $curAction  = $this->getAction()->getId();
?>

<ol class="breadcrumb">
    <li><a href="<?php echo url('admin/main'); ?>">Home</a></li>
    <li>User</li>
    <li class="active"><?php echo ucfirst($curAction); ?></li>
</ol>

<ul class="nav nav-tabs" role="tablist">
    <?php foreach ($subTypes as $type => $lable) : ?>
    <li role="presentation"<?php echo $curType == $type ? 'class="active"' : ''; ?>><a href="<?php echo url('user/'.$curAction.'/utype/'.$type); ?>"><?php echo $lable; ?></a></li>
    <?php endforeach; ?>
</ul>
<p>
    <span style="float: right">Today: <?php echo date("Y-m-d"); ?></span>
</p>
<p>
    <input id="select_all" type="checkbox"> Select all
    <input id="block-button" type="button" value="Block">
    <input id="block-delete-button" type="button" value="Block and delete all posts">
</p>

<table class="table table-striped table-hover">
    <tbody>
        <?php foreach ($list as $one) : ?>
        <tr>
            <td width="30%">
                Username: <?php echo $one['username']; ?><br />
                IP: none<br />
                Status: Active<br />
            </td>
            <td>
                <div class="user-detail-container">
                    <ul style="list-style-type:none" class="list-group">
                        <li style="width: 50%; float: left">gender: <?php echo $gender[ $one['gender'] ];?></li>
                        <li style="width: 50%; float: left">Email: none</li>
                        <li style="width: 50%; float: left">Birthday: none</li>
                        <li style="width: 50%; float: left">Age: <?php echo $one['age']; ?></li>
                        <li style="width: 50%; float: left">Living with: none</li>
                        <li style="width: 50%; float: left">Last Login date: none</li>
                        <li style="width: 50%; float: left">UDID: none</li>
                        <li style="width: 50%; float: left">System: none</li>
                        <li style="width: 50%; float: left">Chosen labels: <?php echo $one['lllness']; ?></li>
                    </ul>
                </div>
            </td>
        </tr>
        <?php endforeach; ?>
        <!--
        <tr>
            <td colspan="2">Activities: Topics: xxx Comments: xx Following: xx Followers: xxx <span style="color: red">Blocked by:</span> xx</td>
        </tr>-->
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

