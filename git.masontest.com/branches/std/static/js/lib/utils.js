/**
 * 公共工具js函数库
 */

var utils = {};

utils.setCookie  = function(name,value,days)
{
    var Days = days ? days : 0;
    var exp  = new Date();
    if (Days) exp.setTime(exp.getTime() + Days*24*60*60*1000);
    document.cookie = name + "="+ escape (value) + ";path=/" + (Days ? ";expires=" + exp.toGMTString() : '');
}

utils.getCookie  = function(name)
{
    var arr = document.cookie.match(new RegExp("(^| )"+name+"=([^;]*)(;|$)"));
    if(arr != null) return unescape(arr[2]); return null;
}

utils.delCookie  = function(name) {
    var exp = new Date();
    exp.setTime(exp.getTime() - 1);
    var cval = this.getCookie(name);
    if (cval != null) document.cookie= name + "=" + cval + ";path=/;expires=" + exp.toGMTString();
}

utils.checkIsLogin = function() {
    if (!this.getCookie('Token')) {
        window.MyBrowser.Routers['index'].navigate('!/login', {trigger: true});
        return false;
    }
    return true;
}

/*
 * 创建一个层
 */
utils.createDiv = function(id)
{
    if (!$('#'+ id).length)
    {
        $('body').append('<div id="'+ id +'"></div>');
    }

    return $('#'+ id);
};

/**
 * 类javascript alert信息提示框
 * type: error|succeed|question|more...
 * 样式定义在css/ui.css
 */
utils.alert = function(msg, type, callback)
{
    msg = '<div style="padding:40px;font-size:14px;">'+msg+'</div>';
    this.openFloatDialog(msg, {
        title: '操作提示',
        width: 400, 
        height: 'auto',
        modal: true,
        autoOpen: true,
        resizable: false,
        autoResize: true,
        "class": 'utils-confirm',
        position: {my:'center', at:'center'},
        buttons: {
            '关闭': function() {
                $(this).dialog('close');
                if (typeof type == 'function') type();
                else if (typeof callback == 'function') callback();
            }
        }
    }, callback);
};


/**
 * 确认对话框
 */
utils.confirm = function(msg, callback, options)
{
    msg = '<div style="padding:40px;font-size:14px;">'+msg+'</div>';
    this.openFloatDialog(msg, {
        title: '操作确认',
        width: options && options['width'] ? options['width'] : undefined,
        height: 'auto',
        modal: true,
        autoOpen: true,
        resizable: false,
        autoResize: true,
        "class": 'utils-confirm',
        position: {my:'center', at:'center'}
    }, callback);
};


/**
 * iframe
 */
utils.iframe = function(url, callback, options)
{
    !options  && (options = {});
    var iframe  = '<iframe frameborder=0 src="'+url+'" width="100%" height="99.9%"></iframe>';
    var idName  = options['idName'] ? options['idName'] : 'dialog-box'; 
    return this.createDiv(idName).html(iframe).dialog({
        title: options['title'] ? options['title'] : '',
        width: options['width'] ? options['width'] : 400,
        height: options['height'] ? options['height'] : undefined,
        modal: true,
        autoOpen: true,
        resizable: options['resizable'] ? options['resizable'] : false,
        autoResize: options['autoResize'] ? options['autoResize'] : true,
        position: {my:'center', at:'center'},
        btnFullscreen : options['btnFullscreen'] ? options['btnFullscreen'] : false,
        open: function() {
            if (typeof options['afterOpenCallback'] === 'function') {
                options['afterOpenCallback'](this);
            }
        },
        close : function(){typeof callback == 'function' && callback.apply(); $('#player-dialog-box').html(''); },
        buttons: {
        }
    });
};

utils.openFloatDialog  = function(template, options, next) {
    var levelId = typeof options['id'] !== 'undefined' ? options['id'] : 'dialog-box';
    var node    = this.createDiv(levelId);

    var windowW = $(window).width();
    var windowH = $(window).height();
    var w       = windowW * 0.7;
    var h       = windowH * 0.8;

    node.html(template);
    seajs.use('ui', function(){
        node.dialog({
            title: options['title'] ? options['title'] : '',
            width: options['width'] ? options['width'] : w,
            height: options['height'] ? options['height'] : h,
            cache: options['cache'] ? options['cache'] : false,
            modal: options['modal'] ? options['modal'] : true,
            autoOpen: options['autoOpen'] ? options['autoOpen'] : true,
            resizable: options['resizable'] ? options['resizable'] : false,
            autoResize: options['autoResize'] ? options['autoResize'] : false,
            position: options['position'] ? options['position'] : {my:'center', at:'center'},
            closeOnEscape: typeof options['closeOnEscape'] !== 'undefined' ? options['closeOnEscape'] : false,
            draggable: true,
            open: typeof options['open'] !== 'undefined' ? options['open'] : false,
            close: typeof options['close'] !== 'undefined' ? options['close'] : false
        });

        if (typeof options['onClose'] !== 'undefined') {
            $('.panel-tool-close').bind('click', function() {
                options['onClose']();
            });
        }
        utils.changeMapSize('loaded');
        next && next();
    });
};

utils.closeFloatDialog = function(options) {
    if (typeof options === 'undefined') { options = {}; }
    var levelId = typeof options['id'] !== 'undefined' ? options['id'] : 'dialog-box';

    $('#' + levelId).dialog('destroy');;
};

/**
 * 一般用于在操作完成后给用户一个提示信息
 * type: info|error|succeed(默认)|more...
 * 该层会自动关闭
 */
utils.wait = function(msg, type, callback)
{
    var wait = this.createDiv('dialog-wait');
    wait.html(msg);

    if (typeof type == 'undefined')
    {
        wait.attr('class', '');
    }
    else
    {
        if (typeof type == 'function') callback = type;
        else if (typeof type == 'string')
        {
            wait.attr('class', type);
        }
    }

    var left = $(document).scrollLeft() + ($(window).width() / 2 - tips.width() / 2);
    var top = $(document).scrollTop() + ($(window).height() / 2 - tips.height() / 2);
    wait.css({'top': top, 'left': left}).show();
    return wait;
    if (typeof callback == 'function') callback();
};

utils.wait_close = function(wait, callback){
    wait.fadeOut('slow');
    if (typeof callback == 'function') callback();
};


/**
 * 异操作延时等待后，给用户一个提示信息
 * type: info|error|succeed(默认)|more...
 * 该层会自动关闭
 */
utils.tips = function(msg, type, callback)
{
    var tips = this.createDiv('dialog-tips');
    tips.html(msg);

    if (typeof type == 'undefined')
    {
        tips.attr('class', '');
    }
    else
    {
        if (typeof type == 'function') callback = type;
        else if (typeof type == 'string')
        {
            tips.attr('class', type);
        }
    }

    var left = $(document).scrollLeft() + ($(window).width() / 2 - tips.width() / 2);
    var top = $(document).scrollTop() + ($(window).height() / 2 - tips.height() / 2);
    tips.css({'top': top, 'left': left}).show();

    var timerID = setTimeout(function(){
        tips.fadeOut('slow');
        if (typeof callback == 'function') callback();
    }, 3000);

    tips.on('click', function(){
        tips.fadeOut('slow');
        clearTimeout(timerID);
        timerID = null;
        if (typeof callback == 'function') callback();
    });
};

/**
 * 右击js
 * 基于jquery， jquery.easyui.min.js
 * @param event 事件
 * @param array options 
 * @param string id div 
 *
 */
utils.contextMenu = function(e) {
    var left=0, top=0;
    if (e.type && e.type == 'contextmenu') {
        e.preventDefault();
        left = e.pageX;
        top = e.pageY;
    }else if (e.latLng) { //地图
        e.stop();
        left = e.pixel.x+20;
        top = e.pixel.y+117;
    }

    var divId = arguments[2] || 'menu-list-div-id';
    var items = $('#'+divId).children('.menu-item');
    for(var i=0; i<items.length; i++) {
        $('#'+divId).menu('removeItem', items[i]);
    }

    var menuList = arguments[1] || [];
    for (var i=0; i<menuList.length; i++) {
        $('#'+divId).menu('appendItem', menuList[i]);
    }

    $('#'+divId).menu('show', {
        left: left,
        top: top 
    });
    return false; 
};

/**
 * 右击js
 * 基于jquery， jquery.easyui.min.js
 * @param element $('#id') || id
 * @param array div [{text:'text', title:'title', alt:'alt'}, {}]
 *
 */
utils.rightClick = function(element, menuList) {

   if(typeof element == 'string') {
        var obj = $('#'+element);    
    } else {
        var obj = element;
    }

    obj.bind('contextmenu', function(e) {
        return utils.contextMenu(e, menuList);
    });
};


utils.createDiv('dialog-box').ajaxStart(function(){
        //var left = $(document).scrollLeft() + ($(window).width() / 2 - $(this).width() / 2);
        //var top = $(document).scrollTop() + ($(window).height() / 2 - $(this).height() / 2);
        //$(this).css({'top': top, 'left': left}).show();
    }).ajaxSuccess(function(){
    }).ajaxStop(function(){
    }).ajaxError(function(event, XMLHttpRequest, ajaxOptions, thrownError){
        if (XMLHttpRequest.readyState == 0)
        {
            //uc.alert('网络异常，连接超时。');
        }
        else if (thrownError == 'Forbidden')
        {
        }
    });

utils.elementSize = function(element) {
    if (element) {
        if(typeof element == 'string') {
            var obj = $('#' + element);
        } else {
            var obj = element;
        }
    } else {
        var obj = $(window);
    }
    return ({width: obj.width(), height: obj.height()});
};

utils.changeMapSize = function() {
    var winWidth        = $(window).width();
    var winHeight       = $(window).height();

    var leftMenuWidth   = $('.sidebar').width();

    var height = winHeight - 160;
    var width  = winWidth - leftMenuWidth - 60;

    $('#map-canvas').css('width', width + 'px').css('height', height + 'px');
};

utils.initLeftBar   = function() {
    var afterClick  = function(e) {
        var curObj  = $(e.currentTarget);
        var curLiObj    = curObj.parent();
        var subMenus    = $('.ly-dropdown-menu', curLiObj);
        if (subMenus.length === 0) {
            return true;
        }
        
        if (curLiObj.hasClass('active')) {
            curLiObj.removeClass('active');   
        } else {
            curLiObj.addClass('active');
        }
        return false; 
    };

    $('a.main-menu-a').bind('click', function(e) {
        return afterClick(e);
    });
};

utils.showInitInfo = function() {
    if (initInfo !== '') {
        this.alert(initInfo);
    }
};

utils.initDatetimePicker    = function() {
    var picker  = $('.datetimepicker');
    picker.datepicker({
        ampm: false,//上午下午是否显示  
        timeFormat: 'hh:mm',//时间模式  
        stepHour: 1,//拖动时间时的间隔  
        stepMinute: 5,//拖动分钟时的间隔
        showHour: true,//是否显示小时，默认是true
        showMinute: true,//是否显示分钟，默认是true
        buttonImage: '/img/commons/dees.png', 
        dateFormat: 'yy-mm-dd'
    });
    //picker.timepicker();
};

utils.initialize    = function() {
    this.initLeftBar();
    this.showInitInfo();
    //this.initDatetimePicker();
};

$(function() {
    utils.initialize();
});
