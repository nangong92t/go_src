<div class="row show-grid">
    <div class="inner-left-container">
        <?php $this->renderPartial('_menu'); ?>

        <div class="top-button-nav">
            <button id="create-new-vehicle" type="button" class="btn btn-primary btn-xs">创建新车辆</button>
        </div>
        <div class="clearfix"></div>
        <div class="btn-group search-filter-buttons">
            <button class="btn btn-default btn-xs active" type="button">全部车型</button>
            <button class="btn btn-default btn-xs" type="button">大型</button>
            <button class="btn btn-default btn-xs" type="button">中型</button>
            <button class="btn btn-default btn-xs" type="button">微型</button>
            <button class="btn btn-default btn-xs" type="button">越野</button>
            <button class="btn btn-default btn-xs" type="button">商务</button>
        </div>

        <table class="table table-striped table-hover">
            <thead>
                <tr>
                    <th class="cursor">#<span class="caret"></span></th>
                    <th class="cursor">编号<span class="caret"></span></th>
                    <th class="cursor">
                        车牌号
                        <span class="caret"></span>
                    </th>
                    <th class="cursor">类型<span class="caret"></span></th>
                    <th class="cursor">舒适度<span class="caret"></span></th>
                    <th class="cursor">座位数<span class="caret"></span></th>
                    <th class="cursor">所有者<span class="caret"></span></th>
                    <th>管理</th>
                </tr>
            </thead>
            <tbody id="vehicle-list">
                <tr data-vehicleid="1">
                    <td>1</td>
                    <td>Mark</td>
                    <td>Otto</td>
                    <td>@mdo</td>
                    <td>1</td>
                    <td>Mark</td>
                    <td>Otto</td>
                    <td>@mdo</td>
                </tr>
                <tr data-vehicleid="2">
                    <td>2</td>
                    <td>Jacob</td>
                    <td>Thornton</td>
                    <td>@fat</td>
                    <td>1</td>
                    <td>Mark</td>
                    <td>Otto</td>
                    <td>@mdo</td>
                </tr>
            </tbody>
        </table>
    </div>

    <div id="item-detail-container" class="col-md-4"></div>
</div>
