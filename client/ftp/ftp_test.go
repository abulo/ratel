package ftp

import (
	"reflect"
	"testing"
)

func TestConfig_New(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.config.New()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Config.New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ChangeDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.ChangeDir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.ChangeDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ChangeDirToParent(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.ChangeDirToParent(); (err != nil) != tt.wantErr {
				t.Errorf("Client.ChangeDirToParent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_CurrentDir(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.CurrentDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CurrentDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.CurrentDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Delete(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Delete(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_FileSize(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.FileSize(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.FileSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.FileSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DownFile(t *testing.T) {
	type args struct {
		serverPath string
		destPath   string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.DownFile(tt.args.serverPath, tt.args.destPath); (err != nil) != tt.wantErr {
				t.Errorf("Client.DownFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UploadFile(t *testing.T) {
	type args struct {
		srcFullPath string
		serverPath  string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.UploadFile(tt.args.srcFullPath, tt.args.serverPath); (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_List(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name          string
		client        *Client
		args          args
		wantFilesData []FileData
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFilesData, err := tt.client.List(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFilesData, tt.wantFilesData) {
				t.Errorf("Client.List() = %v, want %v", gotFilesData, tt.wantFilesData)
			}
		})
	}
}

func TestClient_NameList(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.NameList(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.NameList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NameList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Rename(t *testing.T) {
	type args struct {
		from string
		to   string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Rename(tt.args.from, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("Client.Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_RemoveDirRecur(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.RemoveDirRecur(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveDirRecur() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_MakeDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.MakeDir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.MakeDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_RemoveDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.RemoveDir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Logout(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Logout(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Logout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Quit(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Quit(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Quit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Login(t *testing.T) {
	tests := []struct {
		name    string
		client  *Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.client.Login(); (err != nil) != tt.wantErr {
				t.Errorf("Client.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
