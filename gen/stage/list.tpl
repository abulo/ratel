@$define "backstage/{{.View}}/list_user.html"$@
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="renderer" content="webkit">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <title></title>
    <link rel="stylesheet" href="/resource/plugin/layui/css/layui.css">
</head>
<style>
    .layui-form-select .layui-input {
        padding-right: 0px;
        cursor: pointer;
    }
    .layui-form-item .layui-inline {
        margin-bottom: 5px;
        margin-right: 5px;
    }
    .layui-btn+.layui-btn {
        margin-left: 0px;
        margin-bottom: 0px;
    }
    .layui-btn-xs {
        height: 20px;
        line-height: 20px;
        padding: 0 5px;
        font-size: 12px;
    }
    .select{
        width: 80px;
        margin-right:0px;
    }
    .magin-top0{
        margin: 0 0;
        margin-top: -18px;
    }
</style>

<body>
    <form class="layui-form" action="@$urlfor "admin_{{.Table}}"$@" method="get">
        <div class="layui-form-item">
            <div class="layui-inline">
                <input type="text" name="keyword" value="" autocomplete="off" class="layui-input"
                    placeholder="关键字:名称/拼音/ID">
            </div>
            <div class="layui-inline select">
                <select name="status">
                    <option value="">状态</option>
                    <option value="active">激活</option>
                    <option value="inactive">待激</option>
                </select>
            </div>
            <div class="layui-inline">
                <button type="submit" class="layui-btn" lay-submit>筛选</button>
                @$if inArray "admin_{{.Table}}_add" .allow$@
                <a class="layui-btn" href="@$urlfor "admin_{{.Table}}_add"$@">添加</a>
                @$end$@
            </div>
        </div>
    </form>
    <table class="layui-table">
        <thead>
            <tr>
                <th>Id</th>
                {{range .Column}}<th>{{CamelStr .ColumnName}}</th>
                {{end}}
                <th>操作</th>
            </tr>
        </thead>
        <tbody class="content_style">
            @$range .list$@
            <tr>
                <td>@$.Id$@</td>
                {{range .Column}}<td>@$.{{CamelStr .ColumnName}}$@</td>
                {{end}}
                <td>
                    @$if inMultiArray  $.allow "admin_{{$.Table}}_show" "admin_{{$.Table}}_edit" "admin_{{$.Table}}_delete"$@
                    <div class="layui-dropdown">
                        <a  href="javascript:;" class="layui-btn layui-btn-xs layui-btn-normal ">操作</a>
                        <ul>
                            @$if inArray "admin_{{$.Table}}_edit" $.allow$@
                            <li><a href="@$urlfor "admin_{{$.Table}}_edit" ":{{$.Pri}}" .Id$@">编辑</a></li>
                            @$end$@
                            @$if inArray "admin_{{$.Table}}_show" $.allow$@
                            <li><a href="@$urlfor "admin_{{$.Table}}_show" ":{{$.Pri}}" .Id$@">查看</a></li>
                            @$end$@
                            @$if inArray "admin_{{$.Table}}_delete" $.allow$@
                            <li><a href="@$urlfor "admin_{{$.Table}}_delete" ":{{$.Pri}}" .Id$@" data-method="post" data-confirm="确定删除吗？">删除</a></li>
                            @$end$@
                        </ul>
                    </div>
                    @$end$@
                </td>
            </tr>
            @$end$@
        </tbody>
        <tfoot>
            <tr>
                <td colspan="{{.ListTotal}}">
                    <div>@$str2html .page$@</div>
                </td>
            </tr>
        </tfoot>
    </table>
</body>
<script src="/resource/plugin/layui/layui.js"></script>
<script>
    layui.use(["buildform","checkbox","tablelist", "form", "page", "upload","dropdown","delete"], function () {
        var form = layui.form,
            upload = layui.upload,
            $ = layui.jquery;
    });
</script>

</html>
@$end$@