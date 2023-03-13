package sql

import (
	"database/sql"
	"reflect"
	"testing"
)

func Test_getSqlConn(t *testing.T) {
	type args struct {
		driverName string
		server     string
		pool       *pool
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
			got, err := getSqlConn(tt.args.driverName, tt.args.server, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSqlConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSqlConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCachedSqlConn(t *testing.T) {
	type args struct {
		driverName string
		server     string
		pool       *pool
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
			got, err := getCachedSqlConn(tt.args.driverName, tt.args.server, tt.args.pool)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCachedSqlConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCachedSqlConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newDBConnection(t *testing.T) {
	type args struct {
		driverName string
		dns        string
		pool       *pool
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
			got, err := newDBConnection(tt.args.driverName, tt.args.dns, tt.args.pool)
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

func Test_dbConnection(t *testing.T) {
	type args struct {
		driverName string
		dns        string
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
			got, err := dbConnection(tt.args.driverName, tt.args.dns)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
