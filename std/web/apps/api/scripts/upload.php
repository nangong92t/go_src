<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<title>文件上传</title>
</head>

<body>
<form action="http://api.ly.com/upload/file" method="post" enctype="multipart/form-data">
	<label for="file">Filename:</label>
	<input type="file" name="file" id="file" /> 
	<input type="input" name="attachId" value="1">
	<br />
	<input type="hidden" name="uploadType" value="form">
	<input type="hidden" name="token" value="21e7da9b0efc87c48b8acc87f2332296">
	<input type="hidden" name="fileType" value="vehicle_image">
	<input type="hidden" name="fileInputName" value="file">
	<input type="submit" name="submit" value="Submit" />
</form>
</body>
</html>
