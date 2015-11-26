框架说明(过渡期)
===========

* ### 这是一个过渡期的临时框架。 在新框架及集成和发布方案确定后代码库结构还需进行一定的调整。各类库需分离单独建代码库。
### 此过渡期框架应当在总体思路上符合将来新框架的设计， 并且在今后迁移至最新框架时，不需要对业务代码进行重构或大的调整。

* ### 集成方式：
	
	1. vendor目录在项目根目录下：
	
		在启动脚本中添加如下代码：
		
			// ROOT_PATH 为项目根目录, 根据实际情况进行定义
			define('ROOT_PATH', realpath('../').DIRECTORY_SEPARATOR); 
			// 加载Bootstrap
			require ROOT_PATH.'Vendor/Bootstrap/Autoloader.php';
			// 初始化autoloader
			\Bootstrap\Autoloader::instance()->init();
			
	1. vendor目录不在项目根目录下：
	
	   在启动脚本中添加如下代码:
    		
			// ROOT_PATH 为项目根目录, 根据实际情况进行定义
			define('ROOT_PATH', realpath('../').DIRECTORY_SEPARATOR); 
			// 加载Bootstrap
			require ROOT_PATH.'../../Vendor/Bootstrap/Autoloader.php';
			// 初始化autoloader
			\Bootstrap\Autoloader::instance()->addRoot(ROOT_PATH)->init();
			
			
* ### Config等定义：

	为了避免在项目代码中进行繁琐的配置代码注册。请在项目根目录下创建配置目录Config, 所有的配置文件有以类文件来定义，以便实现自动按需加载,类文件名即Vendor的名字，如： Memcache.php、PHPServer.php。当然,业务也可根据自己的需要来添加任意合适的配置类。 所有的配置类均在命名空间Config下。
	
* ### 命名约定：
   
   为了规范和提高运行效率，所有的类及其文件请按照命名空来组织。	

