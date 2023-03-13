package sql

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func Test_clickhouseOpen(t *testing.T) {
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
			got, err := clickhouseOpen(tt.args.driverName, tt.args.dns)
			if (err != nil) != tt.wantErr {
				t.Errorf("clickhouseOpen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clickhouseOpen() = %v, want %v", got, tt.want)
			}
		})
	}
}
