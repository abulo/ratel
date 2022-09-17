@$define "backstage/{{.View}}/edit_{{.Table}}.html"$@
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="renderer" content="webkit">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <title></title>
    <link rel="stylesheet" href="/resource/plugin/layui/css/layui.css">
    <style type="text/css">
        .layui-form-label {
            width: 100px;
        }
        .layui-input-block {
            margin-left: 130px;
            min-height: 36px;
        }
        .attachment {
            float: left;
            margin-right: 20px;
        }
        .attachment>.layui-btn {
            width: 102px;
        }
        .layui-upload-list {
            width: 100px;
            height: 100px;
            border: 1px solid #4444;
            vertical-align: middle;
            line-height: 100px;
        }
        .layui-upload-list>img {
            max-height: 100px;
            max-width: 100px;
            vertical-align: middle;
            line-height: 100px;
        }
    </style>
</head>

<body>
    <fieldset class="layui-elem-field layui-field-title" style="margin-top: 20px;">
        <legend>{{.Title}}编辑</legend>
    </fieldset>
    <form class="layui-form" action="@$urlfor "admin_{{.Table}}_update" ":{{.Pri}}" .{{.Pri}}$@" method="post" enctype="multipart/form-data">
        {{range .Column}}<div class="layui-form-item">
            <label class="layui-form-label">{{Helper .ColumnName}}</label>
            <div class="layui-input-self">
                <input type="text" name="{{Helper .ColumnName}}" autocomplete="off" placeholder="请输入 {{Helper .ColumnName}}" class="layui-input"
                    lay-verify="required" value="@$.{{CamelStr .ColumnName}}$@">
            </div>
        </div>
        {{end}}
        <div class="layui-form-item">
            <div class="layui-input-block">
                <button class="layui-btn" lay-submit lay-filter="form">提交</button>
                <button type="reset" class="layui-btn layui-btn-primary">重置</button>
                <a href="@$.backUrl$@" class="layui-btn layui-btn-normal">返回</a>
            </div>
        </div>
    </form>
</body>
<script src="/resource/plugin/layui/layui.js"></script>
<script>
layui.use(["form","upload","laydate","layer"],function () {
    var form=layui.form
        ,layer=layui.layer
        ,upload=layui.upload
        ,laydate=layui.laydate
        ,$=layui.jquery;
    upload.render({
        elem: ".upload", //绑定元素
        url: @$urlfor "admin_upload" ":type" "images"$@, //上传接口
        done: function (res,index,upload) {
            //上传完毕回调
            if(res.path==""||res.url=="") {
                layer.msg("上传失败");
                return;
            }
            var item=this.item;
            item.siblings(".layui-upload-list").find(".layui-upload-img").attr("src",res.url)
            item.siblings(".layui-upload-list").find(".upload_input").attr("value",res.path)
            layer.msg("上传成功");
        },
        error: function () {
            //请求异常回调
            layer.msg("上传失败");
        }
    });
    form.on("submit(form)",function (data) {
        $(".layui-upload-file").remove();
    });
    form.verify({
        otherReq: function (value,item) {
            var $=layui.$;
            var verifyName=$(item).attr("name")
                ,verifyType=$(item).attr("type")
                ,formElem=$(item).parents(".layui-form")//获取当前所在的form元素，如果存在的话
                ,verifyElem=formElem.find("input[name="+verifyName+"]")//获取需要校验的元素
                ,isTrue=verifyElem.is(":checked")//是否命中校验
                ,focusElem=verifyElem.next().find("i.layui-icon");//焦点元素
            if(!isTrue||!value) {
                //定位焦点
                focusElem.css(verifyType=="radio"? {"color": "#FF5722"}:{"border-color": "#FF5722"});
                //对非输入框设置焦点
                focusElem.first().attr("tabIndex","1").css("outline","0").blur(function () {
                    focusElem.css(verifyType=="radio"? {"color": ""}:{"border-color": ""});
                }).focus();
                return "必填项不能为空";
            }
        }
    });
})
</script>

</html>
@$end$@