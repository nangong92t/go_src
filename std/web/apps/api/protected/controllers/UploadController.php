<?php
/**
 * 文件上传类.
 *
 * @author Aven<929313990@qq.com>
 * @date 2014-06-22
 */
class UploadController extends Controller
{

    /**
     * 上传类型 form 普通的表单上传  swfupload flash上传
     * @var array
     */
    protected $_uploadType   = array("form","swfupload");
    /**
     * 上传的文件格式
     * @var array
     */
    protected $_fileType     = array("vehicle_image","vehicle_zip");
    /**
     * 上传文件类型的格式要求规则
     *  ext 后缀,maxsize 文件上传大小限制,host 保存的主机,newLoc 挂着本机的可写路径
     * @var array
     */
    protected $_fileTypeRule = array(
        "image"     =>  array(
            "ext"       => array("jpg","png","gif"),
            "maxsize"   => 2097152,
            "rpcServer" => "basic",
            "newLoc"    => "../../../static/"
        ),
        "zip"       =>  array(
            "ext"       => array("zip"),
            "maxsize"   => 2097152,
            "rpcServer" => "basic",
            "newLoc"    => "../../../static/"
        ),
    );

    /**
     * 输出跨域设置
     *
     * @return void.
     */
    public function actionAllowCross()
    {
        header('Content-type: text/xml');
        $info   = '<?xml version="1.0"?>
<!DOCTYPE cross-domain-policy SYSTEM "http://www.macromedia.com/xml/dtds/cross-domain-policy.dtd">
<cross-domain-policy>
<allow-access-from domain="std.com" />
<allow-access-from domain="*.std.com" />
<allow-access-from domain="std.com:800" />
<allow-access-from domain="*.std.com:800" />
<allow-access-from domain="masontest.com" />
<allow-access-from domain="*.masontest.com" />
</cross-domain-policy>';
        echo $info;
    }

    /**
     * File Action: Post 上传文件.
     *
     * @return $this->response. 
     */
    public function actionFile()
    {
         try {
           //安全校验
            $token  = "";
            if (isset($_POST['token'])) {
                $token   = $_POST['token'];
            }

            $creator    = Yii::app()->user->getId();

           //获得文件基本属性
           $file = $this->_getFileAtrr();
           //上传类型判断
           //文件上传
        
           //$this->_moveTmpFileToNewLoc($fileType,$file);
            $data   = file_get_contents($file['tmp_name']);

           //RPCserver调用
           $rpcServer = RpcClient_STD_Topic::Instance();
           $result = $rpcServer->addBackground($token, $file['name'], $file['ext'], base64_encode($data));

           $file = $result['data'];

           $this->response['status'] = 200;
           $this->response['mesg'] = "upload success!";
           $this->response['body'] =  $file;

         } catch (Exception $e) {
             $this->response['status'] = 400;
             $this->response['mesg']   = $e->getMessage();
         }
         
         echo json_encode($this->response);
    }

    /**
     * 上传附件move
     * @param  String $fileType 
     * @param  array $file     
     * @return void   
     * @throws Exception If [上传失败]     
     */
    protected function _moveTmpFileToNewLoc($fileType,$file)
    {
        $newLoc = $this->_fileTypeRule[$fileType]['newLoc'].$file['savePath']."/".$file['newName'].".".$file['ext'];
        if(!@move_uploaded_file($file['tmp_name'],$newLoc))
        {
          throw new Exception("UploadErr: upload file flail !");
        }
    }

    /**
     * 检查文件上传格式是否正确
     * @param  string  $fileType
     * @param  string  $ext  
     * @return void
     * @throws Exception If [后缀非法]
     */
    protected function _checkFileExt($fileType, $ext)
    {

        if (!in_array($ext,$this->_fileTypeRule[$fileType]['ext']) )
        {
            throw new Exception("UploadErr:不合法的上传文件格式 !");
        }
    }

    /**
     * 检查文件上传格式是否正确
     * @param  string  $fileType
     * @param  string  $ext  
     * @return void
     * @throws Exception If [超过上传大小]
     */
    protected function _checkFileSize($filesize)
    {
        $maxsize = 2097152;
        if($maxsize < $filesize)
        {
          $fileK = round($maxsize/1024);
          throw new Exception("UploadErr:上传文件大小超过 最大限制{$fileK} k!");
        }

    }

    /**
     * 获得组装后附件信息
     * @return array
     * @throws Exception If [not file with upload]
     */
    
    protected function _getFileAtrr()
    {
        
        $fileInputName   = $_POST['fileInputName'];
        // $savaInfo = $this->_getSaveInfo($fileType,$uploadType);
        $fileTmpDir = $_FILES[$fileInputName]['tmp_name'];
        if(!$_FILES[$fileInputName]['tmp_name']) throw new Exception("UploadErr:没有上传文件 !");
        $file  = array();
        
        if(is_file($fileTmpDir)){
            $file['name']        = $_FILES[$fileInputName]['name'];
            $file['tmp_name']    = $_FILES[$fileInputName]['tmp_name'];
            $pFileS              = explode(".", $file['name']);
            $file['ext']         = $pFileS[ count($pFileS) - 1];
        }
        return $file;
    }

    /**
     * 获得新上传附件目录
     * @param  String $fileType   
     * @param  String $uploadType 
     * @return array
     */
    protected function _getSaveInfo($fileType,$uploadType)
    {
        $newLoc = $this->_fileTypeRule[$fileType]['newLoc'];
        $host   = str_replace('http://', '', Yii::app()->clientScript->getCoreScriptUrl());
        $subLoc = 'uploads/' . date("Y/m/d");
        $newLoc = $newLoc.$subLoc;
        $this->_MkFolder($newLoc);
        return array("saveHost"=>$host,"savePath"=>$subLoc);
    }

    /**
     * mkdir new folder
     * @param  String $folder 
     * @return void
     */
    protected function _MkFolder($folder)
    {

      if(!is_readable($folder))
        {
            $this->_MkFolder(dirname($folder) );
            if(!is_file($folder)) @mkdir($folder,0777);
        }

    }

}
