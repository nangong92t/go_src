<?php

$toPost     = function($url, $params=array()) {
    $ch = curl_init();//初始化curl
    curl_setopt($ch,CURLOPT_URL, $url);//抓取指定网页
    //curl_setopt($ch, CURLOPT_HEADER, 0);//设置header
    //curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);//要求结果为字符串且输出到屏幕上
    curl_setopt($ch, CURLOPT_POST, 1);//post提交方式
    curl_setopt($ch, CURLOPT_POSTFIELDS, $params);
    $data = curl_exec($ch);//运行curl

    if ($error = curl_error($ch) ) {
        die($error);
    }
    curl_close($ch);
    
    return $data;
};

// $url = "http://api.ly.com/upload/file";
// $params = array (
//     'token'  => '21e7da9b0efc87c48b8acc87f2332296',
//     'attachId'  => '1'
//     #'imgFile'   => '@/tmp/1385022923940.jpg' 
// );

//$url    = 'http://koubei.jumeicd.com/Ajax/UploadReportImage';
$url    = 'http://api.ly.com/rest/21e7da9b0efc87c48b8acc87f2332296/Vehicle_Type/addType';
$params = array (
    'data[name]'  => '大型车'
    #'imgFile'   => '@/tmp/1385022923940.jpg' 
);

/**
$url    = 'http://tony.koubei.jumeicd.com/Ajax/AddReport';
$params = array (
    'user_id'   => 1497,
    'time'      => 1395565072,
    'verify_code'   => '2f2458d6476a95f8d513eedfc3c36061',

    'product_id'    => 10992,
    'title'     => 'just test asdf asdf asdf asdf asdf asd asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asd ',
    'content'   => 'just test for conten tasd fasdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf asdf ',
    //'rating'    => array('930'=>4, '931'=>2, '932'=>1, '933'=>5),      // ping fen
    'deal_hash_id'  => 'bj120515p10991',         // product hash id
    'order_id'  =>  234,        // dingdan
    'img_list'  => ''           
);

**/


print_r($toPost($url, $params)."\n\n");//输出结果

