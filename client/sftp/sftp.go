package sftp

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// transferBufPool 复用 32KB 传输缓冲，避免每次文件传输都由 io.Copy 内部新建临时缓冲。
//
// PERF-5-21 [P2]：缓冲复用；递归目录传输保持串行——单 SSH/SFTP 通道不支持并发操作，
// 并发改造超出本次范围、风险较高，不实施。
var transferBufPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 32*1024)
		return &buf
	},
}

// Config 保存 SFTP 服务器的连接配置。
type Config struct {
	// Host 为服务器主机名或 IP 地址。
	Host string
	// Port 为服务端口，字符串形式。
	Port string
	// User 为登录用户名。
	User string
	// Password 为登录密码。
	Password string
	// Timeout 为 SSH 拨号超时时间。
	Timeout time.Duration
}

// Client 封装已建立的 SSH 连接与打开的 SFTP 通道，提供远端文件操作。
//
// 零值不可直接使用，须由 Config.New 创建；SSHClient 与 SFTPClient 均为底层句柄，可直接调用。
type Client struct {
	// SSHClient 为底层 SSH 客户端连接。
	SSHClient *ssh.Client
	// SFTPClient 为底层 SFTP 客户端通道。
	SFTPClient *sftp.Client
}

// IOReaderProgress 在读取流的同时累计已传输字节数，供上传下载进度展示使用。
type IOReaderProgress struct {
	io.Reader
	// TransferredBytes 指向已传输字节总数，由 Read 回调持续累加。
	TransferredBytes *int64
}

// Read 转发底层 Reader 的读取，并把本次读取的字节数累加到 TransferredBytes。
func (iorp *IOReaderProgress) Read(p []byte) (int, error) {
	n, err := iorp.Reader.Read(p)
	*iorp.TransferredBytes += int64(n)
	return n, err
}

// New 依据配置建立 SSH 连接并打开 SFTP 通道；拨号或通道创建失败时返回错误。
//
// 当前 ClientConfig 未设置 HostKeyCallback；按 x/crypto/ssh 的校验规则，主机密钥回调为空时
// ssh.Dial 会直接返回 “must specify HostKeyCallback” 错误。
func (config *Config) New() (*Client, error) {
	clientConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		Timeout: config.Timeout,
	}
	conn, err := ssh.Dial("tcp", config.Host+":"+config.Port, clientConfig)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}
	client, err := sftp.NewClient(conn)
	if err != nil {
		logger.Logger.Error(err)
		return nil, err
	}

	return &Client{
		SSHClient:  conn,
		SFTPClient: client,
	}, nil
}

// DownFile 将 serverPath 处的远程文件下载到本地 destPath，返回写入的字节数。
func (client *Client) DownFile(serverPath, destPath string) (int64, error) {
	return getTransfer(client.SFTPClient, serverPath, destPath, nil, nil)
}

// DownFileWithProgress 下载 serverPath 处的远程文件到 destPath，通过 transferred 累计已传输
// 字节、通过 total 回填文件总大小。
func (client *Client) DownFileWithProgress(serverPath, destPath string, transferred *int64, total *int64) (int64, error) {
	return getTransfer(client.SFTPClient, serverPath, destPath, transferred, total)
}

// GetRecursively 递归下载 remotePath 目录到本地 localPath；本地同名文件已存在且非目录时返回错误。
func (client *Client) GetRecursively(remotePath, localPath string) error {
	remoteWalker := client.SFTPClient.Walk(remotePath)
	if remoteWalker == nil {
		return errors.New("SFTP Walker Error")
	}
	for remoteWalker.Step() {
		err := remoteWalker.Err()
		if err != nil {
			return err
		}
		remoteFullFilepath := remoteWalker.Path()
		localFilepath, _ := getRecursivelyPath(localPath, remotePath, remoteFullFilepath)

		if remoteWalker.Stat().IsDir() {
			localStat, localStatErr := os.Stat(localFilepath)
			// 本地已存在同名条目且不是目录时无法创建目录，直接报错。
			if !os.IsNotExist(localStatErr) && !localStat.IsDir() {
				return errors.New("Cannot create a directry when that file already exists")
			}
			mode := remoteWalker.Stat().Mode()
			if os.IsNotExist(localStatErr) {
				mkErr := os.Mkdir(localFilepath, mode)
				if mkErr != nil {
					return mkErr
				}
			}
			continue
		}
		_, getErr := getTransfer(client.SFTPClient, localFilepath, remoteFullFilepath, nil, nil)
		if getErr != nil {
			return getErr
		}
	}
	return nil
}

// getTransfer 执行远程到本地的文件下载：创建本地文件、打开远程文件，按需回填总大小并累计进度，
// 写入完成后同步本地文件；任一环节失败时返回错误。
func getTransfer(client *sftp.Client, remoteFilepath, localFilepath string, tfBytes *int64, totalBytes *int64) (int64, error) {
	localFile, localFileErr := os.Create(localFilepath)
	if localFileErr != nil {
		return 0, errors.New("localFileErr: " + localFileErr.Error())
	}
	defer func() {
		if err := localFile.Close(); err != nil {
			logger.Logger.Error("Error closing localFile: ", err)
		}
	}()

	remoteFile, remoteFileErr := client.Open(remoteFilepath)
	if remoteFileErr != nil {
		return 0, errors.New("remoteFileErr: " + remoteFileErr.Error())
	}
	defer func() {
		if err := remoteFile.Close(); err != nil {
			logger.Logger.Error("Error closing remoteFile: ", err)
		}
	}()

	var bytes int64
	var copyErr error
	// 传入 totalBytes 时先读取远程文件大小回填。
	if totalBytes != nil {
		f, _ := remoteFile.Stat()
		*totalBytes = f.Size()
	}
	// PERF-5-21 [P2]：从缓冲池取 32KB 缓冲配合 io.CopyBuffer 流式拷贝，传输结束归还缓冲池；
	// 读取/写入方向、进度累计与错误返回语义保持不变。
	bufPtr := transferBufPool.Get().(*[]byte)
	defer transferBufPool.Put(bufPtr)
	if tfBytes != nil {
		remoteFileWithProgress := &IOReaderProgress{Reader: remoteFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.CopyBuffer(localFile, remoteFileWithProgress, *bufPtr)
	} else {
		bytes, copyErr = io.CopyBuffer(localFile, remoteFile, *bufPtr)
	}
	if copyErr != nil {
		return 0, errors.New("copyErr: " + copyErr.Error())
	}

	syncErr := localFile.Sync()
	if syncErr != nil {
		return 0, errors.New("syncErr: " + syncErr.Error())
	}
	return bytes, nil
}

// UploadFile 将本地 localFilepath 单文件上传到远程 remoteFilepath，返回写入的字节数。
func (client *Client) UploadFile(localFilepath string, remoteFilepath string) (int64, error) {
	return putTransfer(client.SFTPClient, localFilepath, remoteFilepath, nil, nil)
}

// UploadFileWithProgress 上传本地文件并通过 transferred 累计已传输字节、通过 total 回填文件总大小。
func (client *Client) UploadFileWithProgress(localFilepath string, remoteFilepath string, transferred *int64, total *int64) (int64, error) {
	return putTransfer(client.SFTPClient, localFilepath, remoteFilepath, transferred, total)
}

// PutRecursively 递归上传本地 localPath 目录到远程 remotePath；远程同名文件已存在且非目录时返回错误。
func (client *Client) PutRecursively(localPath string, remotePath string) error {
	localWalkerErr := filepath.Walk(localPath, func(localFullFilepath string, info os.FileInfo, fileErr error) error {
		if fileErr != nil {
			return fileErr
		}
		remoteFilepath, _ := getRecursivelyPath(remotePath, localPath, localFullFilepath)
		if info.IsDir() {
			remoteStat, remoteStatErr := client.SFTPClient.Stat(remoteFilepath)
			if !os.IsNotExist(remoteStatErr) && !remoteStat.IsDir() {
				return errors.New("Cannot create a directry when that file already exists")
			}
			mode := info.Mode()
			if os.IsNotExist(remoteStatErr) {
				mkErr := client.SFTPClient.Mkdir(remoteFilepath)
				if mkErr != nil {
					return mkErr
				}
				chErr := client.SFTPClient.Chmod(remoteFilepath, mode)
				if chErr != nil {
					return chErr
				}
			}
			return nil
		}
		_, getErr := putTransfer(client.SFTPClient, localFullFilepath, remoteFilepath, nil, nil)
		if getErr != nil {
			return getErr
		}
		return nil
	})
	return localWalkerErr
}

// putTransfer 执行本地到远程的文件上传：创建远程文件、打开本地文件，按需回填总大小并累计进度。
func putTransfer(client *sftp.Client, localFilepath string, remoteFilepath string, tfBytes *int64, totalBytes *int64) (int64, error) {
	remoteFile, remoteFileErr := client.Create(remoteFilepath)
	if remoteFileErr != nil {
		return 0, errors.New("remoteFileErr: " + remoteFileErr.Error())
	}

	defer func() {
		if err := remoteFile.Close(); err != nil {
			logger.Logger.Error("Error closing remoteFile: ", err)
		}
	}()

	localFile, localFileErr := os.Open(localFilepath)
	if localFileErr != nil {
		return 0, errors.New("localFileErr: " + localFileErr.Error())
	}

	defer func() {
		if err := localFile.Close(); err != nil {
			logger.Logger.Error("Error closing localFile: ", err)
		}
	}()

	var bytes int64
	var copyErr error
	// 传入 totalBytes 时先读取本地文件大小回填。
	if totalBytes != nil {
		f, _ := localFile.Stat()
		*totalBytes = f.Size()
	}
	// PERF-5-21 [P2]：上传同样走缓冲池复用的 io.CopyBuffer，避免每次传输新建临时缓冲；
	// 进度累计与错误返回语义保持不变。
	bufPtr := transferBufPool.Get().(*[]byte)
	defer transferBufPool.Put(bufPtr)
	if tfBytes != nil {
		localFileWithProgress := &IOReaderProgress{Reader: localFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.CopyBuffer(remoteFile, localFileWithProgress, *bufPtr)
	} else {
		bytes, copyErr = io.CopyBuffer(remoteFile, localFile, *bufPtr)
	}

	if copyErr != nil {
		return 0, errors.New("copyErr: " + copyErr.Error())
	}
	return bytes, nil
}

// Close 依次关闭 SFTP 通道与 SSH 连接，返回各自关闭过程中遇到的全部错误。
func (client *Client) Close() (errors []error) {
	if client.SFTPClient != nil {
		sftpErr := client.SFTPClient.Close()
		if sftpErr != nil {
			errors = append(errors, sftpErr)
		}
	}
	if client.SSHClient != nil {
		sshErr := client.SSHClient.Close()
		if sshErr != nil {
			errors = append(errors, sshErr)
		}
	}
	return
}

// getRecursivelyPath 计算 remotePath 下与 remoteFullFilepath 对应的本地（或远程）目标路径，
// 并统一转换为斜杠分隔形式；路径相对关系无法计算时返回错误。
func getRecursivelyPath(localPath string, remotePath string, remoteFullFilepath string) (string, error) {
	rel, err := filepath.Rel(filepath.Clean(remotePath), remoteFullFilepath)
	if err != nil {
		return "", err
	}
	localFilepath := filepath.Join(localPath, rel)
	return filepath.ToSlash(localFilepath), nil
}

// Quit 是 Close 的别名，关闭全部连接。
func (client *Client) Quit() (errors []error) {
	return client.Close()
}

// Chmod 修改指定远端文件或目录的权限位。
func (client *Client) Chmod(path string, mode os.FileMode) error {
	return client.SFTPClient.Chmod(path, mode)
}

// Chown 修改指定文件的属主用户与属组。
func (client *Client) Chown(path string, uid, gid int) error {
	return client.SFTPClient.Chown(path, uid, gid)
}

// Remove 删除指定文件或目录；路径不存在，或目录非空时返回错误。
func (client *Client) Remove(path string) error {
	return client.SFTPClient.Remove(path)
}

// RemoveDirectory 删除指定的空目录。
func (client *Client) RemoveDirectory(path string) error {
	return client.SFTPClient.RemoveDirectory(path)
}

// Rename 重命名远端文件或目录。
func (client *Client) Rename(oldname, newname string) error {
	return client.SFTPClient.Rename(oldname, newname)
}

// PosixRename 使用 posix-rename@openssh.com 扩展重命名文件；目标已存在时会被直接替换。
func (client *Client) PosixRename(oldname, newname string) error {
	return client.SFTPClient.PosixRename(oldname, newname)
}

// Getwd 返回服务器当前工作目录；涉及相对路径的操作均以该目录为基准。
func (client *Client) Getwd() (string, error) {
	return client.SFTPClient.Getwd()
}

// Mkdir 创建指定目录；同名文件或目录已存在，或父目录不存在时返回错误，不会递归创建完整路径。
func (client *Client) Mkdir(path string) error {
	return client.SFTPClient.Mkdir(path)
}

// MkdirAll 创建 path 目录及其全部父目录；path 已是目录时不做处理并返回 nil，
// 路径中包含普通文件时返回错误。
func (client *Client) MkdirAll(path string) error {
	return client.SFTPClient.MkdirAll(path)
}
