<style type="text/css">
			.tlup{
				margin-top: 5px; 
				padding: 10px;
				border: solid 1px #ccc;
				text-align: left;
			}
			.shuom{
				font-size: 16px;
				color: #ff004b;
				margin-top: 5px;
			}
			.txtlist{
				margin-top: 30px;
			}
			.txtlist span{
				margin-left: 20px;
			}
			.txtlist a{
				color: blue;
			}
			.txtlist ol li{
				margin-top: 5px;
				list-style: disc;
				list-style-type: decimal;
			}
		</style>
		<div>
			<div class="shuom">上传后会在下方列表显示并支持下载，审核通过后会在作者所在区域显示。请谨慎下载网友上传的不明格式的文件，防止中毒。</div>
			<form class="tlup">
			  <div class="layui-upload">
			  <label class="layui-btn layui-btn-normal" for="xFile"><i class="layui-icon">&#xe67c;</i>选择文件（5M以内 1KB以上）</label>
			  <input type="file" id="xFile" style="position:absolute;clip:rect(0 0 0 0);">
			  <span id="sfile" class="layui-inline layui-upload-choose"></span>
			 <!--  <button id="delfile" type="button" class="layui-btn"  onclick="onDelFile()">删除</button> -->
			  <button id="upf" class="layui-btn" style="float: right;" type="button" onclick="onUpload()">开始上传</button>
			  </div></form>
		</div>
		<div class="txtlist">
			<ul>
				{{.}}
			</ul>
		</div>
		<script type="text/javascript">
			$('#xFile').change(function(em){
			    $('#delfile').show();
			    $('#sfile').text(em.currentTarget.files[0].name);
			    path = $('#xFile').val();
			    return false;
			});

			function onUpload(){

		        var files = $('#xFile').prop('files');//获取到文件列表

		        if (files.length == 0) {
		            layui.use('layer', function(){
			          layui.layer.alert("请选择文件！", {icon: 5});
			        }); 
		            return;
		        } else {
		        	$('#upf').attr("disabled",true);
		        	$('#upf').text("正在上传...");
		        	var index = layer.load(1, {shade: [0.1,'#bbb'],time:0});
		        	var filename = $('#sfile').text();
		            var reader = new FileReader();
		            reader.readAsText(files[0], "gb2312");
		            reader.onload = function (evt) {
		                var fileString = evt.target.result;

		                $.ajax({
						    url:"/frame/uploadtxt?name="+filename,
						    data: fileString,
						    cache:false,
						    contentType:false,
						    processData:false,
						    type:'POST',
						    success:function(data){
						    	layer.close(index);
						    	if (data == "") {
			                		layui.use('layer', function(){
							          layui.layer.alert("上传成功！", {icon: 6});
							        });  
			                	}else{
			                		layui.use('layer', function(){
							          layui.layer.alert(data, {icon: 5});
							        }); 
			                	}
			                	loadPage('/frame/upload');
						    }
						});
		            }
		        }
			};
		</script>