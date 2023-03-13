package sql

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Row 获取记录
type Row struct {
	rows *Rows
	err  error
}

// ToInterface ...
func (r *Row) ToAny() (result map[string]any, err error) {
	if r.rows.err != nil {
		err = r.rows.err
		return
	}
	items, err := r.rows.ToAny()
	if err != nil {
		r.err = err
		return nil, err
	}
	return items[0], nil
}

// ToMap get Map
func (r *Row) ToMap() (result map[string]string, err error) {
	if r.rows.err != nil {
		err = r.rows.err
		return
	}
	items, err := r.rows.ToMap()
	if err != nil {
		r.err = err
		return nil, err
	}
	return items[0], nil
}

// ToStruct get Struct
func (r *Row) ToStruct(st any) error {
	if r.rows.err != nil {
		return r.rows.err
	}
	//获取变量的类型
	stType := reflect.TypeOf(st)
	stVal := reflect.ValueOf(st)
	if stType.Kind() != reflect.Ptr {
		return fmt.Errorf("the variable type is %v, not a pointer", stType.Kind())
	}
	stTypeInd := stType.Elem()
	v := reflect.New(stTypeInd)
	tagList, err := extractTagInfo(v)
	if err != nil {
		return err
	}
	fields, err := r.rows.rows.Columns()
	if err != nil {
		return err
	}
	refs := make([]any, len(fields))
	for i, field := range fields {
		if f, ok := tagList[field]; ok {
			refs[i] = f.Addr().Interface()
		} else {
			refs[i] = new(any)
		}
	}
	if err := r.rows.rows.Scan(refs...); err != nil {
		return err
	}
	stVal.Elem().Set(v.Elem())
	return nil
}

// Rows
type Rows struct {
	rows *sql.Rows
	err  error
}

// ToAny 用法
func (r *Rows) ToAny() (data []map[string]any, err error) {
	if r.rows == nil {
		return nil, r.err
	}
	fields, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return nil, err
	}
	data = make([]map[string]any, 0)
	num := len(fields)
	refs := make([]any, num)
	for i := 0; i < num; i++ {
		var ref any
		refs[i] = &ref
	}
	for r.rows.Next() {
		result := make(map[string]any)
		if err := r.rows.Scan(refs...); err != nil {
			return nil, err
		}
		for i, field := range fields {
			if val, err := toString(refs[i]); err == nil {
				result[field] = any(val)
			} else {
				return nil, err
			}
		}
		data = append(data, result)
	}
	return data, nil
}

// ToMap 用法
func (r *Rows) ToMap() (data []map[string]string, err error) {
	if r.rows == nil {
		return nil, r.err
	}
	fields, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return nil, err
	}
	data = make([]map[string]string, 0)
	num := len(fields)
	refs := make([]any, num)
	for i := 0; i < num; i++ {
		var ref any
		refs[i] = &ref
	}
	for r.rows.Next() {
		result := make(map[string]string)
		if err := r.rows.Scan(refs...); err != nil {
			return nil, err
		}
		for i, field := range fields {
			if val, err := toString(refs[i]); err == nil {
				result[field] = val
			} else {
				return nil, err
			}
		}
		data = append(data, result)
	}
	return data, nil
}

// ToStruct 用法
func (r *Rows) ToStruct(st any) error {
	//st->&[]user
	//获取变量的类型,类型为指针
	stType := reflect.TypeOf(st)
	//获取变量的值
	stVal := reflect.ValueOf(st)
	stValInd := reflect.Indirect(stVal)
	//1.参数必须是指针
	if stType.Kind() != reflect.Ptr {
		return fmt.Errorf("the variable type is %v, not a pointer", stType.Kind())
	}
	//指针指向的类型:slice
	stTypeInd := stType.Elem()
	//2.传入的类型必须是slice,slice的成员类型必须是struct
	if stTypeInd.Kind() != reflect.Slice || stTypeInd.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("the variable type is %v, not a slice struct", stType.Elem().Kind())
	}
	if r.rows == nil {
		return r.err
	}
	//初始化struct
	v := reflect.New(stTypeInd.Elem())
	//提取结构体中的tag
	tagList, err := extractTagInfo(v)
	if err != nil {
		return err
	}
	fields, err := r.rows.Columns()
	if err != nil {
		r.err = err
		return err
	}
	refs := make([]any, len(fields))
	for i, field := range fields {
		//如果对应的字段在结构体中有映射，则使用结构体成员变量的地址
		if f, ok := tagList[field]; ok {
			refs[i] = f.Addr().Interface()
		} else {
			refs[i] = new(any)
		}
	}
	for r.rows.Next() {
		if err := r.rows.Scan(refs...); err != nil {
			return err
		}
		stValInd = reflect.Append(stValInd, v.Elem())
	}
	stVal.Elem().Set(stValInd)
	return nil
}
