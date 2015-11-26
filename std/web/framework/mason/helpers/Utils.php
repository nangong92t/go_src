<?php
/**
 * 综合工具类
 *
 */

final class Utils
{
    static $fileExtension = array('gif','jpg','jpeg','png');

    /**
     * 访问远程数据
     */
    public static function fetchRemoteResult($component, $action, $params) {
        if (!($ip = N::cache()->get('server_ip'))) {
            $ip = N::db()->queryScalar('config', 'value',
                array('varname = :varname', array(':varname' => 'server_domain')));
            #$ip = 'test.anoah.com'; # for test
            N::cache()->set('server_ip', $ip);
        }

        $internetUrl = 'http://' . $ip . '/api/?q=json/' . str_replace('.', '/', $component) . '/' . $action;

        $result = self::httpPost($internetUrl, array('info' => json_encode($params)));
        return json_decode($result);
    }

    /**
     *
     */
    public static function httpPost($url, array $data) {
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_POST, 1);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
        $params = curl_exec($ch);
        curl_close($ch);

        return $params;
    }

    /**
     * 建立文件夹
     *
     * @param string $aimUrl
     * @return viod
     **/
    public static function createDir($aimUrl) {
        $aimUrl = str_replace('', '/', $aimUrl);
        $aimDir = '';
        $arr = explode('/', $aimUrl);
        foreach ($arr as $str) {
            $aimDir .= $str . '/';
            if (!file_exists($aimDir)) {
                mkdir($aimDir);
            }
        }
    }

    /**
     * 建立文件
     *
     * @param string $aimUrl
     * @param boolean $overWrite 该参数控制是否覆盖原文件
     * @return boolean
     **/
    public static function createFile($aimUrl, $overWrite = false) {
        if (file_exists($aimUrl) && $overWrite == false) {
            return false;
        } elseif (file_exists($aimUrl) && $overWrite == true) {
            self::unlinkFile($aimUrl);
        }
        $aimDir = dirname($aimUrl);
        self::createDir($aimDir);
        touch($aimUrl);
        return true;
    }

    /**
     * 移动文件夹
     *
     * @param string $oldDir
     * @param string $aimDir
     * @param boolean $overWrite 该参数控制是否覆盖原文件
     * @return boolean
     **/
    public static function moveDir($oldDir, $aimDir, $overWrite = false) {
        $aimDir = str_replace('', '/', $aimDir);
        $aimDir = substr($aimDir, -1) == '/' ? $aimDir : $aimDir . '/';
        $oldDir = str_replace('', '/', $oldDir);
        $oldDir = substr($oldDir, -1) == '/' ? $oldDir : $oldDir . '/';
        if (!is_dir($oldDir)) {
            return false;
        }
        if (!file_exists($aimDir)) {
            self::createDir($aimDir);
        }
        @$dirHandle = opendir($oldDir);
        if (!$dirHandle) {
            return false;
        }
        while(false !== ($file = readdir($dirHandle))) {
            if ($file == '.' || $file == '..') {
                continue;
            }
            if (!is_dir($oldDir.$file)) {
                self::moveFile($oldDir . $file, $aimDir . $file, $overWrite);
            } else {
                self::moveDir($oldDir . $file, $aimDir . $file, $overWrite);
            }
        }
        closedir($dirHandle);
        return rmdir($oldDir);
    }

    /**
     * 移动文件
     *
     * @param string $fileUrl
     * @param string $aimUrl
     * @param boolean $overWrite 该参数控制是否覆盖原文件
     * @return boolean
     **/
    public static function moveFile($fileUrl, $aimUrl, $overWrite = false) {
        if (!file_exists($fileUrl)) {
            return false;
        }
        if (file_exists($aimUrl) && $overWrite = false) {
            return false;
        } elseif (file_exists($aimUrl) && $overWrite = true) {
            self::unlinkFile($aimUrl);
        }
        $aimDir = dirname($aimUrl);
        self::createDir($aimDir);
        rename($fileUrl, $aimUrl);
        return true;
    }

    /**
     * 删除文件夹
     *
     * @param string $aimDir
     * @return boolean
     **/
    public static function unlinkDir($aimDir) {
        $aimDir = str_replace('', '/', $aimDir);
        $aimDir = substr($aimDir, -1) == '/' ? $aimDir : $aimDir.'/';
        if (!is_dir($aimDir)) {
            return false;
        }
        $dirHandle = opendir($aimDir);
        while(false !== ($file = readdir($dirHandle))) {
            if ($file == '.' || $file == '..') {
                continue;
            }
            if (!is_dir($aimDir.$file)) {
                self::unlinkFile($aimDir . $file);
            } else {
                self::unlinkDir($aimDir . $file);
            }
        }
        closedir($dirHandle);
        return rmdir($aimDir);
    }

    /**
     * 删除文件
     *
     * @param string $aimUrl
     * @return boolean
     **/
    public static function unlinkFile($aimUrl) {
        if (file_exists($aimUrl)) {
            unlink($aimUrl);
            return true;
        } else {
            return false;
        }
    }

    /**
     * 复制文件夹
     *
     * @param string $oldDir
     * @param string $aimDir
     * @param boolean $overWrite 该参数控制是否覆盖原文件
     * @return boolean
     **/
    public static function copyDir($oldDir, $aimDir, $overWrite = false) {
        $aimDir = str_replace('', '/', $aimDir);
        $aimDir = substr($aimDir, -1) == '/' ? $aimDir : $aimDir.'/';
        $oldDir = str_replace('', '/', $oldDir);
        $oldDir = substr($oldDir, -1) == '/' ? $oldDir : $oldDir.'/';
        if (!is_dir($oldDir)) {
            return false;
        }
        if (!file_exists($aimDir)) {
            self::createDir($aimDir);
        }
        $dirHandle = opendir($oldDir);
        while(false !== ($file = readdir($dirHandle))) {
            if ($file == '.' || $file == '..') {
                continue;
            }
            if (!is_dir($oldDir . $file)) {
                self::copyFile($oldDir . $file, $aimDir . $file, $overWrite);
            } else {
                self::copyDir($oldDir . $file, $aimDir . $file, $overWrite);
            }
        }
        return closedir($dirHandle);
    }

    /**
     * 复制文件
     *
     * @param string $fileUrl
     * @param string $aimUrl
     * @param boolean $overWrite 该参数控制是否覆盖原文件
     * @return boolean
     **/
    public static function copyFile($fileUrl, $aimUrl, $overWrite = false) {
        if (!file_exists($fileUrl)) {
            return false;
        }
        if (file_exists($aimUrl) && $overWrite == false) {
            return false;
        } elseif (file_exists($aimUrl) && $overWrite == true) {
            self::unlinkFile($aimUrl);
        }
        $aimDir = dirname($aimUrl);
        self::createDir($aimDir);
        copy($fileUrl, $aimUrl);
        return true;
    }


    /***
     * 上传文件模块
     *
     *
     */
    public static function getUrlFile($url, $bType)
    {
        if ($content  = file_get_contents($url))
        {
            $fileNode = explode('.', $url);
            $extension= $fileNode[ count($fileNode)-1 ];
            if (!in_array($extension, self::$fileExtension)) $extension = 'jpg';

            $fileNode = explode('/', $url);
            $fileName = $fileNode[ count($fileNode)-1 ];

            $fileMd5  = md5($content);
            $c        = new EMongoCriteria;
            $c->hash  = $fileMd5;
            $isExists = Attach::model()->find($c);
            if ($isExists) return $isExists->_id;

            $fileNameMd5  = md5($fileName.time());
            $fileRePath   = substr($fileNameMd5, 0, 1) . '/' . substr($fileNameMd5, 1, 3) . '/' . substr($fileNameMd5, 4, 5);
            $fileNameMd5  = substr($fileNameMd5, 9, strlen($fileNameMd5));

            $static = Yii::app()->params['static'];
            $uploadPath = $static['uploadPath'];

            $filePath = $uploadPath . '/' . $fileRePath;
            $newFile  = $filePath.'/'.$fileNameMd5.'.'.$extension;

            self::createDir($filePath);
            if ($size=file_put_contents($newFile, $content))
            {
                $newAttach  = new Attach();
                $newAttach->uploadType  = $bType;
                $newAttach->type        = mime_content_type($newFile);
                $newAttach->creator     = Yii::app()->user->getId();
                $newAttach->name        = $fileName;
                $newAttach->size        = $size;
                $newAttach->extension   = $extension;
                $newAttach->hash        = $fileMd5;
                $newAttach->isdel       = 0;
                $newAttach->savepath    = $fileRePath;
                $newAttach->savename    = $fileNameMd5.'.'.$extension;
                $newAttach->action      = 0;

                $newAttach->save();
                return $newAttach->_id;
            }
        }

        return 0;
    }


    /**
     * 获取当前url根菜单.
     *
     * @param string  $ctrName       控制器名称.
     * @param string  $actName       动作名称.
     * @param array   $menus         项目配置的所有菜单数组.
     *
     * @return array.
     */
    public static function getRootCA($ctrName, $actName, array $menus)
    {
        $curCA  = $ctrName . '/' . $actName;
        $isFind = 0;

        $getSubMenu = function($curMenus) use (&$getSubMenu, $curCA, &$isFind) {
            $result = array();
            if (isset($curMenus[$curCA])) { return array($curCA); }

            foreach ($curMenus as $ca => $detail) {
                if (isset($detail['subMenus']) && isset($detail['subMenus'][$curCA])) {
                    $result[]   = $ca;
                    $isFind     = 1;
                } else if (isset($detail['subMenus'])) {
                    /*
                    $subRes     = $getSubMenu($detail['subMenus']);
                    $result[]   = $ca;
                    if (isset($subRes[0])) $result[]   = $subRes[0];
                    */
                }
                if ($result) break;
            }
            return $result;
        };

        foreach ($menus as $k => $one) {
            $result = $getSubMenu(array($k => $one));
            if ($isFind) break;
        }

        return $result;
    }

}

