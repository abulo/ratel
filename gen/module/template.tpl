package {{.Package}}

import (
	"safety/dao"
	"safety/initial"
	"context"

	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/pkg/errors"
)

// Create{{.Dao}} 创建数据
func Create{{.Dao}}(ctx context.Context, data dao.{{.Dao}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("`{{.TableName}}`").Insert(data)
}

// Update{{.Dao}} 更新数据
func Update{{.Dao}}(ctx context.Context, id int64, data dao.{{.Dao}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("`{{.TableName}}`").Where("id", id).Update(data)
}

// Delete{{.Dao}} 删除数据
func Delete{{.Dao}}(ctx context.Context, id int64) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("`{{.TableName}}`").Where("id", id).Delete()
}

// Show{{.Dao}} 获取数据
func Show{{.Dao}}(ctx context.Context, id int64) (dao.{{.Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{.Dao}}
	err := db.NewBuilder(ctx).Table("`{{.TableName}}`").Where("id", id).Row().ToStruct(&res)
	return res, err
}

// List{{.Dao}} 全部数据
func List{{.Dao}}(ctx context.Context, condition map[string]interface{}) ([]dao.{{.Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{.Dao}}
	builder :=  db.NewBuilder(ctx).Table("`{{.TableName}}`")

	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}

	if !util.Empty(condition["pageOffset"]){
		builder.Offset(condition["pageOffset"])
	}

	if !util.Empty(condition["pageSize"]){
		builder.Limit(condition["pageOffset"])
	}
	err := builder.OrderBy("id", query.DESC).Rows().ToStruct(&res)
	return res, err
}

// Item{{.Dao}} 条件数据
func Item{{.Dao}}(ctx context.Context, condition map[string]interface{}) (dao.{{.Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{.Dao}}
	builder :=  db.NewBuilder(ctx).Table("`{{.TableName}}`")
	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}
	err := builder.Row().ToStruct(&res)
	return res, err
}

// Total{{.Dao}} 数据记录数
func Total{{.Dao}}(ctx context.Context, condition map[string]interface{}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder :=  db.NewBuilder(ctx).Table("`{{.TableName}}`")
	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}
	return builder.Count()
}



{{range .Func}}


// List{{$.Dao}}By{{.FuncName}} 获取数据
func List{{$.Dao}}By{{.FuncName}}(ctx  context.Context, condition map[string]interface{})([]dao.{{$.Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res []dao.{{$.Dao}}
	builder :=  db.NewBuilder(ctx).Table("`{{$.TableName}}`")

	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}

	if !util.Empty(condition["pageOffset"]){
		builder.Offset(condition["pageOffset"])
	}

	if !util.Empty(condition["pageSize"]){
		builder.Limit(condition["pageOffset"])
	}
	err := builder.OrderBy("id", query.DESC).Rows().ToStruct(&res)
	return res, err
}


// Total{{$.Dao}}By{{.FuncName}} 获取数据总量
func Total{{$.Dao}}By{{.FuncName}}(ctx  context.Context, condition map[string]interface{})(int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	builder :=  db.NewBuilder(ctx).Table("`{{$.TableName}}`")

	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}

	return builder.Count()
}



// Item{{$.Dao}}By{{.FuncName}} 获取数据
func Item{{$.Dao}}By{{.FuncName}}(ctx  context.Context, condition map[string]interface{})(dao.{{$.Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	var res dao.{{$.Dao}}
	builder :=  db.NewBuilder(ctx).Table("`{{$.TableName}}`")
	{{range  $elem :=   .CondiTion}}
	if !util.Empty(condition["{{Helper $elem}}"]){
		builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	}
	{{end}}
	err := builder.OrderBy("id", query.DESC).Row().ToStruct(&res)
	return res, err
}

{{if eq .NonUnique 0}}
// Delete{{$.Dao}}By{{.FuncName}} 获取数据
func Delete{{$.Dao}}By{{.FuncName}}(ctx  context.Context, condition map[string]interface{})(int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	builder :=  db.NewBuilder(ctx).Table("`{{$.TableName}}`")
	{{range  $elem :=   .CondiTion}}
	if util.Empty(condition["{{Helper $elem}}"]){
		return 0,errors.New("CondiTion {{Helper $elem}} Is Empty")
	}
	builder.Where("`{{$elem}}`",condition["{{Helper $elem}}"])
	{{end}}
	return builder.Delete()
}
{{end}}
{{end}}