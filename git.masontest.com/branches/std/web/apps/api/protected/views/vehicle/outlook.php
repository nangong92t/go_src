:<div class="row show-grid">
    <div class="inner-left-container">
        <?php $this->renderPartial('_menu'); ?>

        <div class="top-button-nav">
            <button id="create-new" type="button" class="btn btn-primary btn-xs">创建新外形摸板</button>
        </div>
        <div class="clearfix"></div>

        <table class="table table-striped table-hover">
            <thead>
                <tr>
                    <th class="cursor">#</th>
                    <th class="cursor">名称</th>
                    <th class="cursor">创建时间</th>
                    <th class="cursor">创建者</th>
                    <th>管理</th>
                </tr>
            </thead>
            <tbody id="vehicle-list">
                <tr data-vehicletypeid="1">
                    <td>1</td>
                    <td>模型1</td>
                    <td>2014-06-20</td>
                    <td>sss</td>
                    <td>
                        <a href="#">编辑</a>
                        <a href="#">删除</a>
                    </td>
                </tr>
                <tr data-vehicletypeid="2">
                    <td>2</td>
                    <td>模型2</td>
                    <td>2014-06-20</td>
                    <td>sss</td>
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
            什么是车辆外形摸板？
        </h5>
        <p>
            什么， 什么什么
        </p> 
    </div>
</div>
