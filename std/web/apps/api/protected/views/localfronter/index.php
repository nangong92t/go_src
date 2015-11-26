<div class="row show-grid">
    <div class="inner-left-container">
        <ol class="breadcrumb">
            <li><a href="<?php echo url('admin/main'); ?>">首页</a></li>
            <li><a href="<?php echo url('resource/index'); ?>">资源管理</a></li>
            <li class="active">地接商管理</li>
        </ol>

        <ul class="nav nav-tabs">
            <li class="active"><a href="#">地接商列表</a></li>
        </ul>

        <div class="top-button-nav">
            <button id="create-new" type="button" class="btn btn-primary btn-xs">创建新地接商</button>
        </div>
        <div class="clearfix"></div>
        <div class="btn-group search-filter-buttons">
            <button class="btn btn-default btn-xs active" type="button">过滤条件</button>
            <button class="btn btn-default btn-xs" type="button">条件1</button>
            <button class="btn btn-default btn-xs" type="button">条件2</button>
            <button class="btn btn-default btn-xs" type="button">条件3</button>
            <button class="btn btn-default btn-xs" type="button">条件4</button>
        </div>

        <table class="table table-striped table-hover">
            <thead>
                <tr>
                    <th class="cursor">#<span class="caret"></span></th>
                    <th class="cursor">名称<span class="caret"></span></th>
                    <th class="cursor">景区<span class="caret"></span></th>
                    <th class="cursor">规模<span class="caret"></span></th>
                    <th class="cursor">联系人<span class="caret"></span></th>
                    <th class="cursor">入股时间<span class="caret"></span></th>
                    <th class="cursor">资源<span class="caret"></span></th>
                    <th class="cursor">当前接团数<span class="caret"></span></th>
                    <th class="cursor">历史接团数<span class="caret"></span></th>
                    <th class="cursor">好评数<span class="caret"></span></th>
                    <th>管理</th>
                </tr>
            </thead>
            <tbody id="localfronter-list">
                <tr data-localfronterid="1">
                    <td>1</td>
                    <td class="cursor localfronter-name">乐山乐风</td>
                    <td class="cursor attraction-name">乐山大佛</td>
                    <td>小于5k</td>
                    <td>百合</td>
                    <td>2014-05-18</td>
                    <td class="cursor localfronter-resource">10</td>
                    <td class="cursor current-group">1</td>
                    <td class="cursor history-group">2/3</td>
                    <td class="cursor user-rank">-23</td>
                    <td>
                        <a href="#">编辑</a>
                        <a href="#">删除</a>
                    </td>
                </tr>
                <tr data-localfronterid="2">
                    <td>2</td>
                    <td class="cursor localfronter-name">乐山乐风</td>
                    <td class="cursor attraction-name">乐山大佛</td>
                    <td>小于5k</td>
                    <td>百合</td>
                    <td>2014-05-18</td>
                    <td class="cursor localfronter-resource">10</td>
                    <td class="cursor current-group">1</td>
                    <td class="cursor history-group">2/3</td>
                    <td class="cursor user-rank">-23</td>
                    <td>
                        <a href="#">编辑</a>
                        <a href="#">删除</a>
                    </td>
                </tr>
            </tbody>
        </table>
    </div>

    <div id="item-detail-container" class="col-md-4">
        <h5 class="detail-title">
            地界商管理
        </h5>
        <p>
            这里放置相关说明.
        </p> 
    </div>
</div>
