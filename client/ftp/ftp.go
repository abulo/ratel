package ftp

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/abulo/ratel/v2/logger"
	"github.com/jlaffaye/ftp"
)

// Config 配置
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}

//FileData 文件
type FileData struct {
	Name string
	Size uint64
	Type string
}

// Client 客户端连接
type Client struct {
	ftp    *ftp.ServerConn
	config *Config
}

//New 新建连接
func (config *Config) New() *Client {
	conn, err := ftp.Dial(config.Host+":"+config.Port, ftp.DialWithTimeout(config.Timeout))
	if err != nil {
		logger.Logger.Error(err)
		return nil
	}
	return &Client{
		ftp:    conn,
		config: config,
	}
}

// ChangeDir 将当前目录更改为指定的路径
func (client *Client) ChangeDir(path string) error {
	return client.ftp.ChangeDir(path)
}

//ChangeDirToParent 将当前目录更改为父目录。这类似于对ChangeDir的调用，路径设置为“..”
func (client *Client) ChangeDirToParent() error {
	return client.ftp.ChangeDirToParent()
}

//CurrentDir 返回当前目录的路径
func (client *Client) CurrentDir() (string, error) {
	return client.ftp.CurrentDir()
}

//Delete 从远程FTP服务器删除指定的文件
func (client *Client) Delete(path string) error {
	return client.ftp.Delete(path)
}

//FileSize 返回文件的大小
func (client *Client) FileSize(path string) (int64, error) {
	return client.ftp.FileSize(path)
}

// DownFile 下载文件到本地
func (client *Client) DownFile(serverPath, destPath string) error {
	resp, err := client.ftp.Retr(serverPath)
	if err != nil {
		return err
	}
	fileData, err := io.ReadAll(resp)
	if err != nil {
		return err
	}
	resp.Close()
	return os.WriteFile(destPath, fileData, 0664)
}

// UploadFile 上传本地文件到服务器
func (client *Client) UploadFile(srcFullPath, serverPath string) error {
	file, err := os.Open(srcFullPath)
	if err != nil {
		return err
	}
	err = client.MakeDir(filepath.Dir(serverPath))
	if err != nil {
		return err
	}
	defer file.Close()
	return client.ftp.Stor(serverPath, file)
}

//List 返回文件夹下所有文件
func (client *Client) List(path string) (filesData []FileData, err error) {
	err = client.ChangeDir(path)
	if err != nil {
		return nil, err
	}
	entries, err := client.ftp.List("")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fileType := "other"
		switch entry.Type {
		case ftp.EntryTypeFile:
			fileType = "file"
		case ftp.EntryTypeFolder:
			fileType = "directory"
		}
		filesData = append(filesData, FileData{Name: entry.Name, Size: entry.Size, Type: fileType})
	}
	return
}

//NameList 返回文件夹下所有文件
func (client *Client) NameList(path string) ([]string, error) {
	err := client.ChangeDir(path)
	if err != nil {
		return nil, err
	}
	return client.ftp.NameList("")
}

// Rename ..
func (client *Client) Rename(from, to string) error {
	return client.ftp.Rename(from, to)
}

// RemoveDirRecur ..
func (client *Client) RemoveDirRecur(path string) error {
	return client.ftp.RemoveDirRecur(path)
}

// MakeDir ..
func (client *Client) MakeDir(path string) error {
	return client.ftp.MakeDir(path)
}

// RemoveDir ..
func (client *Client) RemoveDir(path string) error {
	return client.ftp.RemoveDir(path)
}

// Logout ..
func (client *Client) Logout() error {
	return client.ftp.Logout()
}

// Quit ..
func (client *Client) Quit() error {
	return client.ftp.Quit()
}

//Login 验证登录
func (client *Client) Login() error {
	return client.ftp.Login(client.config.User, client.config.Password)
}
