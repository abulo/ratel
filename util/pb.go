package util

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	pb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
)

// labValue 递归将 protobuf Value 解码为普通 Go 值。
func labValue(value *pb.Value) (any, error) {
	var err error
	if value == nil {
		return nil, nil
	}
	if structValue, ok := value.GetKind().(*pb.Value_StructValue); ok {
		result := make(map[string]any)
		for k, v := range structValue.StructValue.Fields {
			result[k], err = labValue(v)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if listValue, ok := value.GetKind().(*pb.Value_ListValue); ok {
		result := make([]any, len(listValue.ListValue.Values))
		for i, el := range listValue.ListValue.Values {
			result[i], err = labValue(el)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if _, ok := value.GetKind().(*pb.Value_NullValue); ok {
		return nil, nil
	}
	if numValue, ok := value.GetKind().(*pb.Value_NumberValue); ok {
		return numValue.NumberValue, nil
	}
	if strValue, ok := value.GetKind().(*pb.Value_StringValue); ok {
		return strValue.StringValue, nil
	}
	if boolValue, ok := value.GetKind().(*pb.Value_BoolValue); ok {
		return boolValue.BoolValue, nil
	}
	return errors.New(fmt.Sprintf("Cannot convert the value %+v", value)), nil
}

// StructToMap 将 protobuf Struct 递归转换为 map[string]any。
//
// 嵌套结构体转为嵌套 map，列表转为切片，null 转为 nil；递归过程中出现错误时原样返回。
func StructToMap(str *pb.Struct) (map[string]any, error) {
	var err error
	result := make(map[string]any)
	for k, v := range str.Fields {
		result[k], err = labValue(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

// labEntry 递归将普通 Go 值编码为 protobuf Value；不支持的类型返回错误。
func labEntry(entry any) (*pb.Value, error) {
	var err error
	if entry == nil {
		return &pb.Value{Kind: &pb.Value_NullValue{}}, nil
	}
	rt := reflect.TypeOf(entry)
	switch rt.Kind() {
	case reflect.String:
		if realValue, ok := entry.(string); ok {
			return &pb.Value{Kind: &pb.Value_StringValue{StringValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert string value")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Int())}}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: float64(reflect.ValueOf(entry).Uint())}}, nil
	case reflect.Float32, reflect.Float64:
		return &pb.Value{Kind: &pb.Value_NumberValue{NumberValue: reflect.ValueOf(entry).Float()}}, nil
	case reflect.Bool:
		if realValue, ok := entry.(bool); ok {
			return &pb.Value{Kind: &pb.Value_BoolValue{BoolValue: realValue}}, nil
		}
		return nil, errors.New("cannot convert boolean value")
	case reflect.Array, reflect.Slice:
		lstEntry := reflect.ValueOf(entry)

		lstValue := &pb.ListValue{Values: make([]*pb.Value, lstEntry.Len())}
		for i := 0; i < lstEntry.Len(); i++ {
			lstValue.Values[i], err = labEntry(lstEntry.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return &pb.Value{Kind: &pb.Value_ListValue{ListValue: lstValue}}, nil
	case reflect.Struct:
		return labEntry(structs.Map(entry))
	case reflect.Map:
		mapEntry := make(map[string]any)
		entryValue := reflect.ValueOf(entry)
		for _, k := range entryValue.MapKeys() {
			mapEntry[k.String()] = entryValue.MapIndex(k).Interface()
		}
		structVal, err := MapToStruct(mapEntry)
		return &pb.Value{Kind: &pb.Value_StructValue{StructValue: structVal}}, err
	}
	return nil, errors.New(fmt.Sprintf("Cannot convert [%+v] kind:%s", entry, rt.Kind()))
}

// MapToStruct 将 map[string]any 递归转换为 protobuf Struct。
//
// 结构体先通过反射转为 map；不支持的值类型会导致转换失败并返回错误。
func MapToStruct(input map[string]any) (*pb.Struct, error) {
	var err error
	result := &pb.Struct{Fields: make(map[string]*pb.Value)}
	for k, v := range input {
		result.Fields[k], err = labEntry(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

// StructToPbStruct 将任意结构体转换为 protobuf Struct。
//
// 先通过反射将结构体转为 map[string]any，再交由 MapToStruct 编码；转换失败返回错误。
func StructToPbStruct(input any) (*pb.Struct, error) {
	return MapToStruct(structs.Map(input))
}
