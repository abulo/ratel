package query

import (
	"database/sql"
	"fmt"
	"reflect"
)

// Row 获取记录
type Row struct {
	rs          *Rows
	lastError   error
	transaction bool
}

//ToArray get Array
func (r *Row) ToArray() (result []string, err error) {
	items, err := r.rs.ToArray()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, sql.ErrNoRows
}

//ToMap get Map
func (r *Row) ToMap() (result map[string]string, err error) {
	items, err := r.rs.ToMap()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, sql.ErrNoRows
}

func (r *Row) ToInterface() (result map[string]interface{}, err error) {
	items, err := r.rs.ToInterface()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, sql.ErrNoRows
}

//ToStruct get Struct
func (r *Row) ToStruct(st interface{}) error {
	//获取变量的类型
	stType := reflect.TypeOf(st)

	stVal := reflect.ValueOf(st)

	if stType.Kind() != reflect.Ptr {
		return fmt.Errorf("the variable type is %v, not a pointer", stType.Kind())
	}

	stTypeInd := stType.Elem()

	if r.rs.rs == nil {
		return r.lastError
	}
	defer r.rs.rs.Close()
	v := reflect.New(stTypeInd)
	tagList, err := extractTagInfo(v)
	if err != nil {
		return err
	}
	fields, err := r.rs.rs.Columns()

	if err != nil {
		r.rs.lastError = err
		return err
	}
	refs := make([]interface{}, len(fields))
	for i, field := range fields {
		if f, ok := tagList[field]; ok {
			refs[i] = f.Addr().Interface()
		} else {
			refs[i] = new(interface{})
		}
	}
	if err := r.rs.rs.Scan(refs...); err != nil {
		return err
	}
	stVal.Elem().Set(v.Elem())
	return nil
}

//Rows get data
type Rows struct {
	rs          *sql.Rows
	lastError   error
	transaction bool
}

//ToArray get Array
func (r *Rows) ToArray() (data [][]string, err error) {
	if r.rs == nil {
		return nil, r.lastError
	}
	defer r.rs.Close()
	//获取查询的字段
	fields, err := r.rs.Columns()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	data = make([][]string, 0)
	num := len(fields)
	//根据查询字段的数量，生成[num]interface{}用于存储Scan的结果
	refs := make([]interface{}, num)
	for i := 0; i < num; i++ {
		var ref interface{}
		refs[i] = &ref
	}
	if !r.rs.Next() {
		if err := r.rs.Err(); err != nil {
			return nil, err
		}
		return data, nil
	}
	for r.rs.Next() {
		result := make([]string, len(fields))
		if err := r.rs.Scan(refs...); err != nil {
			return nil, err
		}
		for i := range fields {
			//把*interface{}转换成strings返回
			if val, err := toString(refs[i]); err == nil {
				result[i] = val
			} else {
				return nil, err
			}
		}
		if err != nil {
			r.lastError = err
			return nil, err
		}
		data = append(data, result)
	}
	if len(data) < 1 {
		return nil, sql.ErrNoRows
	}
	return data, nil
}

// ToInterface []map[interface{}]interface{}
func (r *Rows) ToInterface() (data []map[string]interface{}, err error) {
	if r.rs == nil {
		return nil, r.lastError
	}
	defer r.rs.Close()
	fields, err := r.rs.Columns()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	data = make([]map[string]interface{}, 0)
	num := len(fields)
	refs := make([]interface{}, num)
	for i := 0; i < num; i++ {
		var ref interface{}
		refs[i] = &ref
	}
	for r.rs.Next() {
		result := make(map[string]interface{})
		if err := r.rs.Scan(refs...); err != nil {
			return nil, err
		}
		for i, field := range fields {
			if val, err := toString(refs[i]); err == nil {
				result[field] = interface{}(val)
			} else {
				return nil, err
			}
		}
		data = append(data, result)
	}
	if len(data) < 1 {
		return nil, sql.ErrNoRows
	}
	return data, nil

}

//ToMap get Map
func (r *Rows) ToMap() (data []map[string]string, err error) {
	if r.rs == nil {
		return nil, r.lastError
	}
	defer r.rs.Close()
	fields, err := r.rs.Columns()
	if err != nil {
		r.lastError = err
		return nil, err
	}
	data = make([]map[string]string, 0)
	num := len(fields)
	refs := make([]interface{}, num)
	for i := 0; i < num; i++ {
		var ref interface{}
		refs[i] = &ref
	}
	for r.rs.Next() {
		result := make(map[string]string)
		if err := r.rs.Scan(refs...); err != nil {
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
	if len(data) < 1 {
		return nil, sql.ErrNoRows
	}
	return data, nil
}

//ToStruct get Struct
func (r *Rows) ToStruct(st interface{}) error {
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
	if r.rs == nil {
		return r.lastError
	}
	defer r.rs.Close()
	//初始化struct
	v := reflect.New(stTypeInd.Elem())
	//提取结构体中的tag
	tagList, err := extractTagInfo(v)
	if err != nil {
		return err
	}
	fields, err := r.rs.Columns()
	if err != nil {
		r.lastError = err
		return err
	}
	refs := make([]interface{}, len(fields))
	for i, field := range fields {
		//如果对应的字段在结构体中有映射，则使用结构体成员变量的地址
		if f, ok := tagList[field]; ok {
			refs[i] = f.Addr().Interface()
		} else {
			refs[i] = new(interface{})
		}
	}
	for r.rs.Next() {
		if err := r.rs.Scan(refs...); err != nil {
			return err
		}
		stValInd = reflect.Append(stValInd, v.Elem())
	}
	stVal.Elem().Set(stValInd)
	return nil

}
