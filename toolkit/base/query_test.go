package base

import (
	"context"
	"reflect"
	"testing"
)

func TestNewDataType(t *testing.T) {
	tests := []struct {
		name string
		want map[string]DataType
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableList(t *testing.T) {
	type args struct {
		ctx    context.Context
		DbName string
	}
	tests := []struct {
		name    string
		args    args
		want    []Table
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableList(tt.args.ctx, tt.args.DbName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableItem(t *testing.T) {
	type args struct {
		ctx       context.Context
		DbName    string
		TableName string
	}
	tests := []struct {
		name    string
		args    args
		want    Table
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableItem(tt.args.ctx, tt.args.DbName, tt.args.TableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableColumn(t *testing.T) {
	type args struct {
		ctx       context.Context
		DbName    string
		TableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []Column
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableColumn(tt.args.ctx, tt.args.DbName, tt.args.TableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableColumn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableIndex(t *testing.T) {
	type args struct {
		ctx       context.Context
		DbName    string
		TableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []Index
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TableIndex(tt.args.ctx, tt.args.DbName, tt.args.TableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TableIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTablePrimary(t *testing.T) {
	type args struct {
		ctx       context.Context
		DbName    string
		TableName string
	}
	tests := []struct {
		name    string
		args    args
		want    Column
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TablePrimary(tt.args.ctx, tt.args.DbName, tt.args.TableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("TablePrimary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TablePrimary() = %v, want %v", got, tt.want)
			}
		})
	}
}
