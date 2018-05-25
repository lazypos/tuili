<!DOCTYPE html>
<html>
<head>
	<meta name="msvalidate.01" content="1E6D64FB2D22FAE27C63B1C77A421B95" />
	<meta name="google-site-verification" content="jzt3kSkGg8ocnf9knMFaaa34dT0OXRQvFAzns3008oE" />
	<meta name="baidu-site-verification" content="GigBIxXxNP" />
	<title>推理屋-推理小说分享平台</title>
	<meta http-equiv="keywords" content="推理屋,推理屋官网,推理小说分享">
	<meta name="keywords" content="推理屋,推理屋官网,推理小说分享">
	<meta name="description" content="推理屋,推理屋官网,推理小说分享">
	<link rel="stylesheet" href="../../layui/css/layui.css">
	<style type="text/css">     
		.logo{
			margin-right: 80px;
		}
		/*.login{
			margin-left: 80px;
		}*/
	</style>
</head>
<body id="推理屋" style="text-align: center; width:700px; display: block; margin: auto;margin-top: 0px;">
	<div>
		<div>
			<ul class="layui-nav layui-bg-green">
			  <li class="layui-nav-item logo"><b>推理屋-推理小说分享平台</b></li>
			  <li class="layui-nav-item layui-this"><a href="javascript:;" onclick="loadPage('/frame/home')">首页</a></li>
			  <li class="layui-nav-item"><a href="javascript:;" onclick="loadPage('/frame/download')">小说下载</a></li>
			  <li class="layui-nav-item"><a href="javascript:;" onclick="loadPage('/frame/upload')">小说分享</a></li>
			  <!-- <li class="layui-nav-item"><a href="javascript:;" onclick="loadPage('/frame/veriy')">待审核</a></li> -->
			  <li class="layui-nav-item"><a href="javascript:;" onclick="loadPage('/frame/support')">留言区</a></li>
			  <!-- <li class="layui-nav-item login"><a href="javascript:;">注册</a></li>
			  <li class="layui-nav-item"><a href="javascript:;">登陆</a></li> -->
			</ul>
		</div>

		<div class="content"></div>
	</div>

<script src="/layui/layui.js"></script>
<script src="/layui/jquery-1.7.2.min.js"></script>
<script>
	layui.use('layer', function(){});
  layui.use('element', function(){
    var element = layui.element;
  });
	layui.use('layer', function(){});

	function loadPage(path){
		$(".content").load(path + '?d='+ Date.parse(new Date()));
	}

	function loadPage1(path){
		$(".content").load(path + '&d='+ Date.parse(new Date()));
	}

  $(function(){
    loadPage('/frame/home');
  });
</script>
<div style="padding-top: 20px; color: #4b4b4b; font-size: 14px;">www.tuiliwu.cn 推理屋-推理小说</div>
</body>
</html>