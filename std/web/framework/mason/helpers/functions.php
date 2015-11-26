<?php
/**
 * Yii框架专用函数 对Yii框架的系统组件调用方法做缩写
 * @package nlsp.helper
 */

function request()
{
  return Yii::app()->Request;
}

function auth()
{
  return Yii::app()->authManager;
}

function db()
{
  return Yii::app()->db;
}

function cs()
{
  return Yii::app()->getClientScript();
}


function user()
{
  return Yii::app()->getUser();
}


function url($route,$params=array(),$ampersand='&'){
  return Yii::app()->createAbsoluteUrl($route,$params,'http',$ampersand);
}


function param($name)
{
  return Yii::app()->params[$name];
}

function a2p($alias)
{
  return YiiBase::getPathOfAlias($alias);
}

function baseurl(){
  return Yii::app()->request->baseUrl;
}

function staticurl() {
  return Yii::app()->clientScript->getCoreScriptUrl();
}

function isInvalidArray($array){
  $isInvalid = true;
  if(is_array($array)){
    foreach($array as $k=>$v){
      if(is_array($v)){
        if(!isInvalidArray($v)){
          $isInvalid = false;
        }
      }
      elseif(!empty($v)){
        $isInvalid = false;
        break;

      }
    }
  }
  return $isInvalid;
}
