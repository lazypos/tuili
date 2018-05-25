<style type="text/css">
			.zcwtj{
			float: right;
		}
		.liuyan{
			padding: 5px;
		}
		.liuyanbt{
			text-align: right;
			padding: 5px;
			border-bottom: solid 1px #ccc;
		}
		.liuyanbt span{
			margin-left: 10px;
		}
		</style>
		<div style="margin:auto;">
			<div style="margin-top: 5px; margin-bottom: 20px;">
				<div class="layui-collapse">
				  <div class="layui-colla-item">
				    <h2 class="layui-colla-title">最新反馈留言</h2>
				  </div>
				   <div>
				</div>
					<div style="text-align: left;">
						{{.}}
				   </div>

			</div>

			<div>
				<form class="layui-form layui-form-pane" action="">
			 	<div class="layui-form-item layui-form-text">
	          <label class="layui-form-label">留言  (感谢您的意见、建议、反馈、留言)</label>
	          <div class="layui-input-block">
	            <textarea placeholder="请输入内容" class="layui-textarea"></textarea>
	          </div>
	        </div></form>
	        <button class="layui-btn zcwtj" onclick="onsupport()">提交</button>
			</div>
		</div>
		<script type="text/javascript">
			function onsupport(){
				var msg = $('.layui-textarea').val()
				if (msg != "") {
					layui.use('layer', function(){
				      	time: 0
					    layui.layer.alert('确定要提交吗？', {
						    icon: 6
						    ,btn: ['是','否']
						    ,yes:function(index){
						    	layer.close(index);
						    	url = "/frame/supportup?txt="+msg+"&d="+Date.parse(new Date());
								$.get(url, function(data, status){ 
									layui.use('layer', function(){
							          layui.layer.alert("留言成功！", {icon: 6});
							        }); 
									loadPage('/frame/support');
						        });
							}	
						});
					});
				}
			}
		</script>