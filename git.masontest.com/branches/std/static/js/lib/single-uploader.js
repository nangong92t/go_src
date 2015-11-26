/*
 * single-uploader js lib
 *
 * v0.0.1-dev | for ly project
 * build by Sam(samsheng133@gmail.com) on 2014-10-09
 * MIT Licensed
 */
define(function(require, exports, module) {
    exports.singleuploader    = function(options) {
        var container   = $(options['container']);

        container.append('<div id="thumbnails" class="hide"></div>');

        var next    = function() {};
        if (typeof options['afterHandle'] === 'undefined') {
            alert('please input the afterHandle.');
            return false;
        } else {
            next    = options['afterHandle'];
        }

        //同步加载样式
        require('swfupload/swfupload.js');
        require('swfupload/imghandlers.js');

        var swfuploadUrl = staticUri+"/js/swfupload/";
console.log(apiUrl+"/upload/file");

        var swfu = new SWFUpload({
            // Backend Settings
            upload_url: apiUrl+"/upload/file",
            post_params: {
                          "token":token,
                          "uploadType":"swfupload",
                          "fileInputName":"Filedata"
                        },

            // File Upload Settings
            use_query_string: false,
            requeue_on_error: false,
            file_size_limit : "2 MB",   // 2MB
            file_types : "*.jpg; *.png; *.gif; *.jpeg",
            file_types_description : "JPG Images",
            file_upload_limit : "0",

            // Event Handler Settings - these functions as defined in Handlers.js
            //  The handlers are not part of SWFUpload but are part of my website and control how
            //  my website reacts to the SWFUpload events.
            file_queue_error_handler : fileQueueError,
            file_dialog_complete_handler : fileDialogComplete,
            upload_progress_handler : uploadProgress,
            upload_error_handler : uploadError,
            upload_success_handler : function(file, serverData) { uploadSuccess(file, serverData, this, next, function() {}); },
            upload_complete_handler : uploadComplete,

            // Button Settings
            button_placeholder_id : "spanButtonPlaceholder",
            button_image_url : swfuploadUrl+"SmallSpyGlassWithTransperancy_17x18.png",
            button_width: 45,
            button_height: 18,
            button_text : '<span class="button" type="button">Add</span>',
            button_text_style : '.button { font-family: inherit; font-size: 14px; }',
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
