package base

import "testing"

func TestInitPath(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitPath(); (err != nil) != tt.wantErr {
				t.Errorf("InitPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitConfig(); (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitQuery(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitQuery(); (err != nil) != tt.wantErr {
				t.Errorf("InitQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitBase(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitBase(); (err != nil) != tt.wantErr {
				t.Errorf("InitBase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
