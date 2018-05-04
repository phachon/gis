<?php
/**
 * 客户端 PHP 上传示例
 * @author: phachon@163.com
 */

// 上传地址
$url = "http://127.0.0.1:8087/image/upload";
// 预分配的 appname 
$appname = "test";
// 预分配 appname 对应的 appKey 
$appKey = "ad%4a*a&ada@#ada";
// 加密算法
$token = md5($appname.$appKey);
// 设置 header 头
$headers = ['Appname: test', 'Token: '.$token,];
// 文件绝对路径
$file = realpath(__DIR__.'/../image/test.jpg');

echo $file;
// 发送数据
$data = array(
	'upload' => new CURLFile(realpath($file))
);

// 开始 curl 上传
$ch = curl_init();
curl_setopt($ch, CURLOPT_SAFE_UPLOAD, true);
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
curl_setopt($ch, CURLOPT_POST, TRUE);
curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
$response = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

var_dump($response);
