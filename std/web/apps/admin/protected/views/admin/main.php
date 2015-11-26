<ol class="breadcrumb">
    <li class="active">Home</li>
</ol>

<style>
.home-header-right { line-height:55px;text-align: right}
</style>

<div class="row">
    <div class="col-md-8">
        <h3>Welcome <?php echo Yii::app()->user->getName(); ?>!</h3>
    </div>
    <div class="col-md-4 home-header-right">
        Today: 
        <select name="selector" class="">
            <option value="1">test 1</option>
        </select>         
    </div>
</div>

<table class="table table-striped">
    <thead>
        <tr>
            <th>#</th>
            <th>Users</th>
            <th>Topics</th>
            <th>Comments</th>
            <th>Ads clicks</th>
            <th>Download</th>
        </tr>
    </thead>
    <tbody>
        
        <tr>
            <td class="header">New</td>
            <?php foreach ($columns as $one) : ?>
            <td><?php echo $detail[$one]['new']; ?></td>
            <?php endforeach; ?>
            <td>0</td>
            <td>0</td>
        </tr>
        <tr>
            <td class="header">Reported</td>
            <?php foreach ($columns as $one) : ?>
            <td><?php echo $detail[$one]['reported']; ?></td>
            <?php endforeach; ?>
            <td>0</td>
            <td>0</td>
        </tr>
        <tr>
            <td class="header">Deleted</td>
            <?php foreach ($columns as $one) : ?>
            <td><?php echo isset($detail[$one]['blocked']) ? $detail[$one]['blocked'] : $detail[$one]['deleted']; ?></td>
            <?php endforeach; ?>
            <td>0</td>
            <td>0</td>
        </tr>
    </tbody>
</table>

<div id="change-status-map"></div>
