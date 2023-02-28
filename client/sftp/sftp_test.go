package sftp

import (
	"os"
	"reflect"
	"testing"

	"github.com/pkg/sftp"
)

func TestIOReaderProgress_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		iorp    *IOReaderProgress
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.iorp.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("IOReaderProgress.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IOReaderProgress.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestClient_DownFile(t *testing.T) {
	type args struct {
		serverPath string
		destPath   string
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
			got, err := tt.client.DownFile(tt.args.serverPath, tt.args.destPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DownFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.DownFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DownFileWithProgress(t *testing.T) {
	type args struct {
		serverPath  string
		destPath    string
		transferred *int64
		total       *int64
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
			got, err := tt.client.DownFileWithProgress(tt.args.serverPath, tt.args.destPath, tt.args.transferred, tt.args.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DownFileWithProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.DownFileWithProgress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetRecursively(t *testing.T) {
	type args struct {
		remotePath string
		localPath  string
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
			if err := tt.client.GetRecursively(tt.args.remotePath, tt.args.localPath); (err != nil) != tt.wantErr {
				t.Errorf("Client.GetRecursively() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getTransfer(t *testing.T) {
	type args struct {
		client         *sftp.Client
		remoteFilepath string
		localFilepath  string
		tfBytes        *int64
		totalBytes     *int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTransfer(tt.args.client, tt.args.remoteFilepath, tt.args.localFilepath, tt.args.tfBytes, tt.args.totalBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UploadFile(t *testing.T) {
	type args struct {
		localFilepath  string
		remoteFilepath string
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
			got, err := tt.client.UploadFile(tt.args.localFilepath, tt.args.remoteFilepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.UploadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_UploadFileWithProgress(t *testing.T) {
	type args struct {
		localFilepath  string
		remoteFilepath string
		transferred    *int64
		total          *int64
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
			got, err := tt.client.UploadFileWithProgress(tt.args.localFilepath, tt.args.remoteFilepath, tt.args.transferred, tt.args.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UploadFileWithProgress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.UploadFileWithProgress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_PutRecursively(t *testing.T) {
	type args struct {
		localPath  string
		remotePath string
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
			if err := tt.client.PutRecursively(tt.args.localPath, tt.args.remotePath); (err != nil) != tt.wantErr {
				t.Errorf("Client.PutRecursively() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_putTransfer(t *testing.T) {
	type args struct {
		client         *sftp.Client
		localFilepath  string
		remoteFilepath string
		tfBytes        *int64
		totalBytes     *int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := putTransfer(tt.args.client, tt.args.localFilepath, tt.args.remoteFilepath, tt.args.tfBytes, tt.args.totalBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("putTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("putTransfer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Close(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantErrors []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrors := tt.client.Close(); !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("Client.Close() = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}

func Test_getRecursivelyPath(t *testing.T) {
	type args struct {
		localPath          string
		remotePath         string
		remoteFullFilepath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getRecursivelyPath(tt.args.localPath, tt.args.remotePath, tt.args.remoteFullFilepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("getRecursivelyPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getRecursivelyPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Quit(t *testing.T) {
	tests := []struct {
		name       string
		client     *Client
		wantErrors []error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErrors := tt.client.Quit(); !reflect.DeepEqual(gotErrors, tt.wantErrors) {
				t.Errorf("Client.Quit() = %v, want %v", gotErrors, tt.wantErrors)
			}
		})
	}
}

func TestClient_Chmod(t *testing.T) {
	type args struct {
		path string
		mode os.FileMode
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
			if err := tt.client.Chmod(tt.args.path, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("Client.Chmod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Chown(t *testing.T) {
	type args struct {
		path string
		uid  int
		gid  int
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
			if err := tt.client.Chown(tt.args.path, tt.args.uid, tt.args.gid); (err != nil) != tt.wantErr {
				t.Errorf("Client.Chown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Remove(t *testing.T) {
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
			if err := tt.client.Remove(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_RemoveDirectory(t *testing.T) {
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
			if err := tt.client.RemoveDirectory(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.RemoveDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Rename(t *testing.T) {
	type args struct {
		oldname string
		newname string
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
			if err := tt.client.Rename(tt.args.oldname, tt.args.newname); (err != nil) != tt.wantErr {
				t.Errorf("Client.Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_PosixRename(t *testing.T) {
	type args struct {
		oldname string
		newname string
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
			if err := tt.client.PosixRename(tt.args.oldname, tt.args.newname); (err != nil) != tt.wantErr {
				t.Errorf("Client.PosixRename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Getwd(t *testing.T) {
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
			got, err := tt.client.Getwd()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Getwd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.Getwd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_Mkdir(t *testing.T) {
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
			if err := tt.client.Mkdir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.Mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_MkdirAll(t *testing.T) {
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
			if err := tt.client.MkdirAll(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Client.MkdirAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
