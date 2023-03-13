package sql

import "testing"

func TestClient_mysqlDns(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.mysqlDns(); got != tt.want {
				t.Errorf("Client.mysqlDns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_clickhouseDns(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.clickhouseDns(); got != tt.want {
				t.Errorf("Client.clickhouseDns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_postgresDns(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.postgresDns(); got != tt.want {
				t.Errorf("Client.postgresDns() = %v, want %v", got, tt.want)
			}
		})
	}
}
