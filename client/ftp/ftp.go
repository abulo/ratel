package ftp

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/jlaffaye/ftp"
)

// Config 保存 FTP 服务器的连接配置。
type Config struct {
	// Host 为服务器主机名或 IP 地址。
	Host string
	// Port 为服务端口，字符串形式。
	Port string
	// User 为登录用户名。
	User string
	// Password 为登录密码。
	Password string
	// Timeout 为拨号连接的超时时间。
	Timeout time.Duration
}

// FileData 描述目录列举出的一个文件或目录条目。
type FileData struct {
	// Name 为条目名称。
	Name string
	// Size 为条目大小（字节）。
	Size uint64
	// Type 为条目类型，取值 file、directory 或 other。
	Type string
}

// Client 封装底层 FTP 连接及其配置，提供文件与目录操作。
//
// 零值不可直接使用，须由 Config.New 创建；连接本身不是并发安全的，同一 Client 的操作应串行调用。
type Client struct {
	ftp    *ftp.ServerConn
	config *Config
}

// New 依据配置拨号连接 FTP 服务器；拨号失败时返回错误。连接成功后须再调用 Login 完成账号认证。
func (config *Config) New() (*Client, error) {
	conn, err := ftp.Dial(config.Host+":"+config.Port, ftp.DialWithTimeout(config.Timeout))
	if err != nil {
		return nil, err
	}
	return &Client{
		ftp:    conn,
		config: config,
	}, nil
}

// ChangeDir 将当前工作目录切换到指定路径。
func (client *Client) ChangeDir(path string) error {
	return client.ftp.ChangeDir(path)
}

// ChangeDirToParent 将当前工作目录切换到父目录，等价于 ChangeDir("..")。
func (client *Client) ChangeDirToParent() error {
	return client.ftp.ChangeDirToParent()
}

// CurrentDir 返回当前工作目录的远程路径。
func (client *Client) CurrentDir() (string, error) {
	return client.ftp.CurrentDir()
}

// Delete 删除服务器上指定的远程文件。
func (client *Client) Delete(path string) error {
	return client.ftp.Delete(path)
}

// FileSize 返回指定远程文件的大小（字节）。
func (client *Client) FileSize(path string) (int64, error) {
	return client.ftp.FileSize(path)
}

// DownFile 下载 serverPath 处的远程文件并写入本地 destPath；下载或写入失败时返回错误。
func (client *Client) DownFile(serverPath, destPath string) error {
	resp, err := client.ftp.Retr(serverPath)
	if err != nil {
		return err
	}
	// PERF-5-20 [P1]：defer 提前注册，后续打开本地文件或流式拷贝出错时响应流同样会被关闭，
	// 顺带修复原实现 ReadAll 失败路径 resp 未关闭的资源泄漏。
	defer func() {
		if err := resp.Close(); err != nil {
			logger.Logger.Error("Error closing response: ", err)
		}
	}()

	// PERF-5-20 [P1]：改为 os.OpenFile + io.Copy 流式拷贝，内存峰值固定为拷贝缓冲大小，
	// 不再经 io.ReadAll 把整个文件读入内存；本地文件权限 0664 与原 os.WriteFile 保持一致。
	file, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			logger.Logger.Error("Error closing file: ", closeErr)
		}
	}()
	_, err = io.Copy(file, resp)
	return err
}

// UploadFile 将本地 srcFullPath 文件上传到 serverPath；上传前会尝试创建远程父目录，该目录创建不做递归。
func (client *Client) UploadFile(srcFullPath, serverPath string) error {
	file, err := os.Open(srcFullPath)
	if err != nil {
		return err
	}
	err = client.MakeDir(filepath.Dir(serverPath))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Logger.Error("Error closing file: ", err)
		}
	}()
	return client.ftp.Stor(serverPath, file)
}

// List 切换到 path 并列举其下全部条目，返回文件与目录信息；目录切换或列举失败时返回错误。
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

// NameList 切换到 path 并返回其下所有条目的名称列表。
func (client *Client) NameList(path string) ([]string, error) {
	err := client.ChangeDir(path)
	if err != nil {
		return nil, err
	}
	return client.ftp.NameList("")
}

// Rename 将远程文件或目录从 from 重命名为 to。
func (client *Client) Rename(from, to string) error {
	return client.ftp.Rename(from, to)
}

// RemoveDirRecur 递归删除指定目录及其全部内容。
func (client *Client) RemoveDirRecur(path string) error {
	return client.ftp.RemoveDirRecur(path)
}

// MakeDir 创建指定远程目录；父目录不存在时由底层命令决定是否报错。
func (client *Client) MakeDir(path string) error {
	return client.ftp.MakeDir(path)
}

// RemoveDir 删除指定的空目录。
func (client *Client) RemoveDir(path string) error {
	return client.ftp.RemoveDir(path)
}

// Logout 从 FTP 服务器登出。
func (client *Client) Logout() error {
	return client.ftp.Logout()
}

// Quit 关闭与 FTP 服务器的连接。
func (client *Client) Quit() error {
	return client.ftp.Quit()
}

// Login 使用配置中的用户名与密码完成认证。
func (client *Client) Login() error {
	return client.ftp.Login(client.config.User, client.config.Password)
}
