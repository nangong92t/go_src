<h2><?php echo $this->pageTitle; ?></h2>
  <ul class="tabs">
      <li<?php if ($active==1): ?> class="active"<?php endif; ?>><a href="<?php echo url('client/oauthClient', array('active'=>1)); ?>">有效API列表</a></li>
      <li<?php if ($active==0): ?> class="active"<?php endif; ?>><a href="<?php echo url('client/oauthClient', array('active'=>0)); ?>">无效API列表</a></li>
      <li<?php if ($active<0) : ?> class="active"<?php endif; ?>><a href="<?php echo url('client/createOauthClient'); ?>">创建新API</a></li>
  </ul>
  <form method="get" action="<?php echo url('client/oauthClient'); ?>">
      <div class="clearfix">
          <input type="hidden" name="r" value="client/oauthClient" />
          <input type="hidden" name="active" value="<?php echo $active; ?>" />
          <input type="text" class="large" placeholder="请输入API KEY" name="keyword" id="keyword">
          <input type="submit" class="btn" value="搜索">
      </div>
  </form>


