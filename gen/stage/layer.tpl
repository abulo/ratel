@$define "backstage/{{.View}}/layer_{{.Table}}.html"$@
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="renderer" content="webkit">
    <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
    <title></title>
    <link rel="stylesheet" href="/static/plugin/layui/css/layui.css">
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
    <form class="layui-form" action="@$urlfor "admin_{{.Table}}_layer"$@" method="get">
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
            </div>
        </div>
    </form>
    <table class="layui-table layui-form magin-top0">
        <thead>
            <tr>
               @$if eq .input "checkbox"$@
                <th><input type="checkbox" name="" title="Id" lay-skin="primary" lay-filter="all"></th>
               @$else$@
                <th>Id</th>
                {{range .Column}}<th>{{CamelStr .ColumnName}}</th>
                {{end}}
               @$end$@
            </tr>
        </thead>
        <tbody class="content_style">
            @$range $index, $elem := .list $@
            <tr>
                @$if eq $.input "checkbox"$@
                <td><input type="checkbox" name="title" value="@$.Id$@" title="@$.Id$@" data-json="@$marshalHtml $elem$@" lay-skin="primary"></td>
                @$else$@
                <td><input type="radio" name="title" value="@$.Id$@" title="@$.Id$@" data-json="@$marshalHtml $elem$@" lay-skin="primary"></td>
                @$end$@
                {{range .Column}}<tr>{{CamelStr .ColumnName}}</tr>
                {{end}}
            </tr>
            @$end $@
        </tbody>
        <tfoot>
            <tr>
                <td><button class="layui-btn layui-btn-sm" id="bind">确定</button></td>
                <td colspan="{{.LayerTotal}}">
                    <div>@$str2html .page$@</div>
                </td>
            </tr>
        </tfoot>
    </table>
</body>
<script src="/static/plugin/layui/layui.js"></script>
<script>
    layui.use(["buildform","checkbox","tablelist", "form", "page", "upload","dropdown","delete"], function () {
        var form = layui.form,
            $ = layui.jquery;
        $("#bind").on("click",function(){
            var data = [];
            $("input[name='title']").each(function (i, v) {
                if ($(v).prop("checked")) {
                    data.push(JSON.parse($(v).attr("data-json")))
                }
            });
            parent.@$js .callback $@(data);
        });
    });
</script>
</html>
@$end $@