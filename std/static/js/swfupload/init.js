var swfu;

window.onload = function() {
    var settings = {

        // Backend Settings
        upload_url: uploadUri, // Relative to the SWF file (or you can use absolute paths)
        post_params: {"sessionId" : sessionId, "type": uploadType},

        // File Upload Settings
        file_size_limit : "30720",   // 30MB
        file_types : "*.jpg;*.gif;*.png;",
        file_types_description : "All Files",
        file_upload_limit   :   "50",  //最多上传30张
        file_queue_limit    :   "50",  //最多选择30张

        // Event Handler Settings (all my handlers are in the Handler.js file)
        file_dialog_start_handler :     fileDialogStart,
        file_queued_handler :           fileQueued,
        file_queue_error_handler :      fileQueueError,
        file_dialog_complete_handler :  fileDialogComplete,
        upload_start_handler :          uploadStart,
        upload_progress_handler :       uploadProgress,
        upload_error_handler :          uploadError,
        upload_success_handler :        uploadSuccess,
        upload_complete_handler :       uploadComplete,
        queue_complete_handler :        queueComplete,

        // Button Settings
        button_image_url : staticUri + "/js/swfupload/add_photo.png",  // Relative to the SWF file
        button_placeholder_id : "flash_upload_select_picture",
        button_width:   74,
        button_height:  25,

        // Flash Settings
        flash_url : staticUri + "/js/swfupload/swfupload.swf",

        custom_settings : {
            progressTarget : "flash_upload_progress",
            cancelButtonId : "btnCancel1"
        },
        debug: false
    };

    swfu = new SWFUpload(settings);
    $('input.primary').bind('click', function() {
      start_upload();
    });
};



//开始上传
function start_upload(){
    swfu.startUpload();
    if (swfu.getStats().files_queued <= 0) {
        $("#muti_edit_photos").submit();
    } else {
        $('#btnUpload').attr('disabled',true).removeClass('btn5').val("Uploading ...");
    }
}

//单图上传回调函数,返回上传完成文件的信息
function ts_upload_success(serverData){
    var data    =   eval('('+serverData+')');
    if(data.status==true){
        return true;
    }else{
        return true;
    }
}

//当文件队列有文件时
function enableUploadButton(file){
    $('#btnUpload').attr('disabled',false).addClass('btn5').val("Start Upload");
}

//全部上传完成
function queueComplete(numFilesUploaded) {
    $("#up_num").val(numFilesUploaded);
    $("#muti_edit_photos").submit();
}

