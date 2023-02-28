package query

import (
	"database/sql"
	"reflect"
	"testing"
)

func Test_getCachedSQLConn(t *testing.T) {
	type args struct {
		driverName string
		server     string
		opt        *Opt
	}
	tests := []struct {
		name    string
		args    args
		want    *pingedDB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCachedSQLConn(tt.args.driverName, tt.args.server, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCachedSQLConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCachedSQLConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSQLConn(t *testing.T) {
	type args struct {
		driverName string
		server     string
		opt        *Opt
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSQLConn(tt.args.driverName, tt.args.server, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSQLConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSQLConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDBConnection(t *testing.T) {
	type args struct {
		driverName string
		datasource string
		opt        *Opt
	}
	tests := []struct {
		name    string
		args    args
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newDBConnection(tt.args.driverName, tt.args.datasource, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("newDBConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newDBConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
