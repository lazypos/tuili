<style type="text/css">
		.layui-colla-title{
			cursor: default;
		}
		.layui-collapse a{
			text-decoration: underline;
			cursor: pointer;
			margin-left: 10px;
			margin-right: 10px;
			display: inline-block;
		}
	</style>
	<div id="推理屋" style="margin-top: 5px; color: #ff0045; font-size: 16px;">共收录推理小说{{.TOTAL}}部 累计下载{{.DownCounts}}次</div>
<div class="layui-collapse" style="margin:auto; margin-top: 5px;">
		  <div class="layui-colla-item">
		    <h2 class="layui-colla-title">欧美地区作家</h2>
		    <div class="layui-colla-content layui-show" style="display: inline-block;">
		    	{{.OM}}
		    </div>
		  </div>
		  <div class="layui-colla-item">
		    <h2 class="layui-colla-title">日本及非欧美作家</h2>
		    <div class="layui-colla-content layui-show">
		    	{{.RB}}
		    </div>
		  </div>
		  <div class="layui-colla-item">
		    <h2 class="layui-colla-title">华语作家</h2>
		    <div class="layui-colla-content layui-show">
		    	{{.ZG}}
		    </div>
		  </div>
		   <div class="layui-colla-item">
		    <h2 class="layui-colla-title">原创及多作者合集</h2>
		    <div class="layui-colla-content layui-show">
		    	{{.YC}}
		    </div>
		  </div>
		</div>