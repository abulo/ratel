package stage

// Stage 模板
func Stage() string {

	tpl := `
package {{.Package}}
import (
	"net/http"
	"github.com/abulo/layout/dao"
	"github.com/abulo/layout/flash"
	"github.com/abulo/layout/module/{{.Module}}"

	"github.com/abulo/ratel/v3/gin"
	"github.com/abulo/ratel/v3/logger"
	"github.com/abulo/ratel/v3/pagination"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
)
// {{.Dao}}List {{.Title}}列表数据
func {{.Dao}}List(ctx *gin.Context) {
	//当前页面
	pageNum := cast.ToInt64(ctx.DefaultQuery("page", "1"))
	//每页多少数据
	perNum := cast.ToInt64(ctx.DefaultQuery("per_num", "15"))
	if pageNum < 1 {
		pageNum = 1
	}
	if perNum < 1 {
		perNum = 15
	}
	//封装页面跳转参数
	var pageURL []interface{}
	params := ctx.Params
	for _, entry := range params {
		pageURL = append(pageURL, ":"+entry.Key)
		pageURL = append(pageURL, entry.Value)
	}
	pageURL = append(pageURL, "page")
	pageURL = append(pageURL, ":num")
	pageURL = append(pageURL, "per_num")
	pageURL = append(pageURL, perNum)
	//数据库查询条件
	condition := make(map[string]interface{})
	offset := perNum * (pageNum - 1)
	condition["pageOffset"] = offset
	condition["pageSize"] = perNum
	//获取数据
	list, _ := {{.Module}}.{{.Dao}}List(ctx.Request.Context(), condition)
	//数据总记录数
	total, _ := {{.Module}}.{{.Dao}}Total(ctx.Request.Context(), condition)
	//分页组件
	objPage := pagination.NewPage(total, pageNum, perNum, gin.URLFor(ctx.GetRoute(), pageURL...))
	//保存当前页面地址
	currentUrl := gin.URLFor(ctx.GetRoute(), pageURL...)
	currentUrl = util.StrReplace(":num", cast.ToString(pageNum), currentUrl, -1)
	flash.PutUrl(ctx, currentUrl)
	//数据解析
	ctx.HTML(http.StatusOK, "backstage/{{.View}}/{{.Table}}_list.html", gin.H{
		"page":  objPage.HTML(),
		"list":  list,
		"allow": ctx.GetStringSlice("permission"),
	})
}
// {{.Dao}}Add {{.Title}}添加数据
func {{.Dao}}Add(ctx *gin.Context) {
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	ctx.HTML(http.StatusOK, "backstage/{{.View}}/{{.Table}}_add.html", gin.H{
		"allow": ctx.GetStringSlice("permission"),
		"backUrl": backUrl,
	})
}
// {{.Dao}}Create {{.Title}}创建数据
func {{.Dao}}Create(ctx *gin.Context) {
	//定义参数
	args := make([]string, 0)
	args = append(
		args,
		{{range .Column}}"{{Helper .ColumnName}}",
		{{end}}
	)
	redirect := make(map[string]interface{})
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	redirect["url"] = backUrl
	var res dao.{{.Dao}}
	//定义数据模型
	for _, arg := range args {
		if util.Empty(ctx.PostForm(arg)) {
			redirect["notice"] = "关键信息缺失"
			ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
				"redirect": redirect,
			})
			return
		}
	}
	{{range $index, $elem := .Column}}res.{{CamelStr $elem.ColumnName}} = {{Convert $elem $elem.ColumnName}}
	{{end}}
	redirect["notice"] = "操作成功"
	if _, err := {{.Module}}.{{.Dao}}Create(ctx.Request.Context(), res); err != nil {
		logger.Logger.Info(err)
		redirect["notice"] = "操作失败"
	}
	ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
		"redirect": redirect,
	})
}
// Edit{{.Dao}} {{.Title}}修改数据
func Edit{{.Dao}}(ctx *gin.Context) {
	{{.Pri}} := cast.ToInt64(ctx.Param("{{.Pri}}"))
	data, err := {{.Module}}.Show{{.Dao}}(ctx.Request.Context(), {{.Pri}})
	redirect := make(map[string]interface{})
	redirect["notice"] = "数据不存在"
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	redirect["url"] = backUrl
	if err != nil {
		logger.Logger.Info(err)
		ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
			"redirect": redirect,
		})
		return
	}
	//数据解析
	ctx.HTML(http.StatusOK, "backstage/{{.View}}/edit_{{.Table}}.html", gin.H{
		"{{.Pri}}":  {{.Pri}},
		"data":  data,
		"allow": ctx.GetStringSlice("permission"),
		"backUrl": backUrl,
	})
}
// Update{{.Dao}} {{.Title}}更新数据
func Update{{.Dao}}(ctx *gin.Context) {
	{{.Pri}} := cast.ToInt64(ctx.Param("{{.Pri}}"))
	_, err := {{.Module}}.Show{{.Dao}}(ctx.Request.Context(), {{.Pri}})
	redirect := make(map[string]interface{})
	redirect["notice"] = "数据不存在"
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	redirect["url"] = backUrl
	if err != nil {
		logger.Logger.Info(err)
		ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
			"redirect": redirect,
		})
		return
	}
	//定义参数
	args := make([]string, 0)
	args = append(
		args,
		{{range .Column}}"{{Helper .ColumnName}}",
		{{end}}
	)
	var res dao.{{.Dao}}
	//定义数据模型
	for _, arg := range args {
		if util.Empty(ctx.PostForm(arg)) {
			redirect["notice"] = "关键信息缺失"
			ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
				"redirect": redirect,
			})
			return
		}
	}
	{{range $index, $elem := .Column}}res.{{CamelStr $elem.ColumnName}} = {{Convert $elem $elem.ColumnName}}
	{{end}}
	res.UpdateAt = query.NewDateTime(util.Date("Y-m-d H:i:s", util.Now()))
	redirect["notice"] = "操作成功"
	if _, err := {{.Module}}.Update{{.Dao}}(ctx.Request.Context(), {{.Pri}}, res); err != nil {
		logger.Logger.Info(err)
		redirect["notice"] = "操作失败"
	}
	ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
		"redirect": redirect,
	})
}
// Delete{{.Dao}} {{.Title}}删除数据
func Delete{{.Dao}}(ctx *gin.Context) {
	{{.Pri}} := cast.ToInt64(ctx.Param("{{.Pri}}"))
	_, err := {{.Module}}.Show{{.Dao}}(ctx.Request.Context(), {{.Pri}})
	redirect := make(map[string]interface{})
	redirect["notice"] = "数据不存在"
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	redirect["url"] = backUrl
	if err != nil {
		logger.Logger.Info(err)
		ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
			"redirect": redirect,
		})
		return
	}
	redirect["notice"] = "操作成功"
	if _, err := {{.Module}}.Delete{{.Dao}}(ctx.Request.Context(), {{.Pri}}); err != nil {
		logger.Logger.Info(err)
		redirect["notice"] = "操作失败"
	}
	ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
		"redirect": redirect,
	})
}
// Show{{.Dao}} {{.Title}}查看数据
func Show{{.Dao}}(ctx *gin.Context) {
	{{.Pri}} := cast.ToInt64(ctx.Param("{{.Pri}}"))
	data, err := {{.Module}}.Show{{.Dao}}(ctx.Request.Context(), {{.Pri}})
	redirect := make(map[string]interface{})
	redirect["notice"] = "数据不存在"
	backUrl := flash.GetUrl(ctx)
	if !util.Empty(backUrl) {
		backUrl = gin.URLFor("admin_{{.Table}}")
	}
	redirect["url"] = backUrl
	if err != nil {
		logger.Logger.Info(err)
		ctx.HTML(http.StatusOK, "backstage/common/redirect.html", gin.H{
			"redirect": redirect,
		})
		return
	}
	//数据解析
	ctx.HTML(http.StatusOK, "backstage/{{.View}}/show_{{.Table}}.html", gin.H{
		"{{.Pri}}":  {{.Pri}},
		"data":  data,
		"allow": ctx.GetStringSlice("permission"),
		"backUrl": backUrl,
	})
}
// Layer{{.Dao}} {{.Title}}弹框
func Layer{{.Dao}}(ctx *gin.Context) {
	// 当前页面
	pageNum := cast.ToInt64(ctx.DefaultQuery("page", "1"))
	// 每页多少数据
	perNum := cast.ToInt64(ctx.DefaultQuery("per_num", "15"))
	//弹框参数
	callback := ctx.DefaultQuery("callback", "callback")
	input := ctx.DefaultQuery("input", "checkbox")

	if pageNum < 1 {
		pageNum = 1
	}
	if perNum < 1 {
		perNum = 15
	}
	// 封装页面跳转参数
	var pageURL []interface{}
	params := ctx.Params
	for _, entry := range params {
		pageURL = append(pageURL, ":"+entry.Key)
		pageURL = append(pageURL, entry.Value)
	}
	pageURL = append(pageURL, "page")
	pageURL = append(pageURL, ":num")
	pageURL = append(pageURL, "per_num")
	pageURL = append(pageURL, perNum)
	// 数据库查询条件
	condition := make(map[string]interface{})
	offset := perNum * (pageNum - 1)
	condition["pageOffset"] = offset
	condition["pageSize"] = perNum
	// 获取数据
	list, _ := {{.Module}}.List{{.Dao}}(ctx.Request.Context(), condition)
	// 数据总记录数
	total, _ := {{.Module}}.Total{{.Dao}}(ctx.Request.Context(), condition)
	// 分页组件
	objPage := pagination.NewPage(total, pageNum, perNum, gin.URLFor(ctx.GetRoute(), pageURL...))
	// 保存当前页面地址
	// currentUrl := gin.URLFor(ctx.GetRoute(), pageURL...)
	// currentUrl = util.StrReplace(":num", cast.ToString(pageNum), currentUrl, -1)
	// flash.PutUrl(ctx, currentUrl)
	// 数据解析
	ctx.HTML(http.StatusOK, "backstage/{{.View}}/layer_{{.Table}}.html", gin.H{
		"page":     objPage.HTML(),
		"list":     list,
		"allow":    ctx.GetStringSlice("permission"),
		"callback": callback,
		"input":    input,
	})
}`
	return tpl
}
