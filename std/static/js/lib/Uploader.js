/*
 * Uploader js lib
 * 封装了图片上传业务逻辑。
 * v0.0.1-dev | for ly project
 * build by tony(tonycbcd@gmail.com) on 2014-06-11
 * MIT Licensed
 */
define(function(require, exports, module) {
    exports.Uploader    = function(options) {
        var next    = function() {};
        if (typeof options['afterHandle'] === 'undefined') {
            alert('必须设置afterHandle参数');
            return false;
        } else {
            next    = options['afterHandle'];
        }

        var afterDelHandle  = typeof options['afterDelHandle'] !== 'undefined' ? options['afterDelHandle'] : function(){};

        var apis    = {
            'getDefaultPictureInAlbum': 'BasicService_Album/getAlbumWithCoverByIds'
        };

        var deleteImage = function(one) {
            var curObj  = $(one.currentTarget);
            var curAttachId = curObj.data('attachid');
            curObj.prev().remove();
            curObj.remove(); 
            afterDelHandle(curAttachId);
        };

        var renderImages    = function(pictures, serverData) {
            if (!_.isArray(pictures)) {
                pictures    = [pictures];
            }

            var template    = '<img style="margin: 5px; width: 140px; opacity: 1;" src="<%= imgUrl %>" data-pinit="registered"><span data-attachid="<%= attachId %>" class="label del-album-picture-button label-danger" title="删除">X</span>';
            var container   = $('#thumbnails');
            for (var i=0; i<pictures.length; i++) {
                var one     = pictures[i];
                var imgUrl  = '';
                if (typeof one['savehost'] !== 'undefined') {
                    imgUrl  = 'http://' + one['savehost'] + '/' + one['savepath'] + '/' + one['savename'] + '.' + one['extenstion'];
                } else {
                    imgUrl  = one;
                }
                
                var attachId    = typeof serverData !== 'undefined' ? serverData['body']['attachId'] : one['attachId'];
                container.append( _.template(template, {imgUrl: imgUrl, attachId: attachId}) ); 
                next(one);
            }
            
            $('.del-album-picture-button').unbind().bind('click', function(e) {
                if (!confirm('确定删除？')) { return false; }
                deleteImage(e);
            }); 
        };


        var defaultAlbumId  = typeof options['albumId'] !== 'undefined' ? options['albumId'] : 0;
        if (defaultAlbumId) {
            // 加载现有相册中的图片.
            window.Store.getData(apis.getDefaultPictureInAlbum, [[defaultAlbumId]], function(res) {
                renderImages(res);
            });    
        }
        
        var fileType    = 'image';
        if (typeof options['fileType'] !== 'undefined') {
            fileType    = options['fileType'];
        }
        //同步加载样式
        require('swfupload/swfupload.js');
        require('swfupload/imghandlers.js');

        var swfuploadUrl = staticUri+"/js/swfupload/";
        var swfu = new SWFUpload({
            // Backend Settings
            upload_url: apiUrl+"/upload/file",
            post_params: {
                          "token":token,
                          "uploadType":"swfupload",
                          "fileType": fileType,
                          "fileInputName":"Filedata"
                        },

            // File Upload Settings
            file_size_limit : "2 MB",   // 2MB
            file_types : "*.jpg",
            file_types_description : "JPG Images",
            file_upload_limit : "0",

            // Event Handler Settings - these functions as defined in Handlers.js
            //  The handlers are not part of SWFUpload but are part of my website and control how
            //  my website reacts to the SWFUpload events.
            file_queue_error_handler : fileQueueError,
            file_dialog_complete_handler : fileDialogComplete,
            upload_progress_handler : uploadProgress,
            upload_error_handler : uploadError,
            upload_success_handler : function(file, serverData) { uploadSuccess(file, serverData, this, next, renderImages); },
            upload_complete_handler : uploadComplete,

            // Button Settings
            button_image_url : swfuploadUrl+"SmallSpyGlassWithTransperancy_17x18.png",
            button_placeholder_id : "spanButtonPlaceholder",
            button_width: 180,
            button_height: 18,
            button_text : '<span class="button">Select Images <span class="buttonSmall">(2 MB Max)</span></span>',
            button_text_style : '.button { font-family: Helvetica, Arial, sans-serif; font-size: 12pt; } .buttonSmall { font-size: 10pt; }',
            button_text_top_padding: 0,
            button_text_left_padding: 18,
            button_window_mode: SWFUpload.WINDOW_MODE.TRANSPARENT,
            button_cursor: SWFUpload.CURSOR.HAND,
            
            // Flash Settings
            flash_url : swfuploadUrl+"swfupload.swf",

            custom_settings : {
                upload_target : "divFileProgressContainer"
            },
            
            // Debug Settings
            debug: false
        }); 
    };
});
