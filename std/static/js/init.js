seajs.config({
  paths: {
    'T': 'templates',
    'css': '../css',
  },
  alias: {
    'app': 'app',
    'tpl': 'template',
    'underscore': 'lib/underscore-min',
    'backbone': 'lib/backbone-min',
    'utils': 'lib/utils',
    'store': 'lib/store',
    'md5': 'lib/md5',
    'moment': 'lib/moment.min',
    'ui': 'lib/jquery-ui.min',
    'bootstrap': 'lib/bootstrap.min',
    'easyui': 'lib/easyui/jquery.easyui.min',
    'template': 'lib/seajs-text-debug',
    'timepicker': 'lib/jquery-ui-timepicker-addon',
    'editor': 'app/editor',
    'baiduMapLib': 'lib/baidu-map-lib',
    'datetimepicker': 'lib/bootstrap-datetimepicker.min'
  },
  preload: ['ui', 'moment', 'md5', 'store', 'bootstrap', 'lib/utils', 'template', 'underscore','datetimepicker'],
  base: staticUri + '/js/',
  charset: 'utf-8'
});

var apiRoot = '/api/v1/';

if (typeof customApp !== 'undefined' || !_.isEmpty(customApp)) { 
    jQuery(function(){
        seajs.use('backbone', function() {
            seajs.use(customApp[0], function(loader) {
                if (!_.isNull(loader)) {
                    window.Main = new loader[ customApp[1] ](customApp[2]);
                }
            });
        });
        $(window).resize(function(){
            utils.changeMapSize();
        });
    });
}
