package util

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
	pb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
)

func labValue(value *pb.Value) (interface{}, error) {
	var err error
	if value == nil {
		return nil, nil
	}
	if structValue, ok := value.GetKind().(*pb.Value_StructValue); ok {
		result := make(map[string]interface{})
		for k, v := range structValue.StructValue.Fields {
			result[k], err = labValue(v)
			if err != nil {
				return nil, err
			}
		}
		return result, err
	}
	if listValue, ok := value.GetKind().(*pb.Value_ListValue); ok {
		result := make([]interface{}, len(listValue.ListValue.Values))
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

func StructToMap(str *pb.Struct) (map[string]interface{}, error) {
	var err error
	result := make(map[string]interface{})
	for k, v := range str.Fields {
		result[k], err = labValue(v)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

func labEntry(entry interface{}) (*pb.Value, error) {
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

		lstValue := &pb.ListValue{Values: make([]*pb.Value, lstEntry.Len(), lstEntry.Len())}
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
		mapEntry := make(map[string]interface{})
		entryValue := reflect.ValueOf(entry)
		for _, k := range entryValue.MapKeys() {
			mapEntry[k.String()] = entryValue.MapIndex(k).Interface()
		}
		structVal, err := MapToStruct(mapEntry)
		return &pb.Value{Kind: &pb.Value_StructValue{StructValue: structVal}}, err
	}
	return nil, errors.New(fmt.Sprintf("Cannot convert [%+v] kind:%s", entry, rt.Kind()))
}

func MapToStruct(input map[string]interface{}) (*pb.Struct, error) {
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

func StructToPbStruct(input interface{}) (*pb.Struct, error) {
	return MapToStruct(structs.Map(input))
}
