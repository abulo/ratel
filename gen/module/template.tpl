package {{Package}}

import (
	"cloud/dao"
	"cloud/initial"
	"context"
)

// Create 创建数据
func Create(ctx context.Context, data {{Dao}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{TableName}}").Insert(data)
}

// Update 更新数据
func Update(ctx context.Context, id int64, data {{Dao}}) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{TableName}}").Where("id", id).Update(data)
}

// Delete 删除数据
func Delete(ctx context.Context, id int64) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	return db.NewBuilder(ctx).Table("{{TableName}}").Where("id", id).Delete()
}

// Show 获取数据
func Show(ctx context.Context, id int64) ({{Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	var res {{Dao}}
	err := db.NewBuilder(ctx).Table("{{TableName}}").Row().ToStruct(&res)
	return res, err
}

// List 全部数据
func List(ctx context.Context) ([]{{Dao}}, error) {
	db := initial.Core.Store.LoadSQL("mysql").Write()
	var res []{{Dao}}
	err := db.NewBuilder(ctx).Table("{{TableName}}").Rows().ToStruct(&res)
	return res, err
}

// Total 数据记录数
func Total(ctx context.Context) (int64, error) {
	db := initial.Core.Store.LoadSQL("mysql").Read()
	return db.NewBuilder(ctx).Table("{{TableName}}").Count()
}