package sftp

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/abulo/ratel/v3/core/logger"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Config 配置
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Timeout  time.Duration
}

// Client Stored Client val
type Client struct {
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

// IOReaderProgress forwards the Read() call
// Addging transferredBytes
type IOReaderProgress struct {
	io.Reader
	TransferredBytes *int64 // Total of bytes transferred
}

// Read ...
func (iorp *IOReaderProgress) Read(p []byte) (int, error) {
	n, err := iorp.Reader.Read(p)
	*iorp.TransferredBytes += int64(n)
	return n, err
}

// New 新建连接
func (config *Config) New() (*Client, error) {
	clientConfig := &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: ssh.HostKeyCallback(
			func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		),
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

// DownFile 下载文件到本地
func (client *Client) DownFile(serverPath, destPath string) (int64, error) {
	return getTransfer(client.SFTPClient, serverPath, destPath, nil, nil)
}

// DownFileWithProgress [Experimental] Get with Display Processing Bytes
func (client *Client) DownFileWithProgress(serverPath, destPath string, transferred *int64, total *int64) (int64, error) {
	return getTransfer(client.SFTPClient, serverPath, destPath, transferred, total)
}

// GetRecursively is Recursively Download entire directories
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

		// if Not Exist Mkdir
		if remoteWalker.Stat().IsDir() {
			localStat, localStatErr := os.Stat(localFilepath)
			// 存在するかつディレクトリではない場合エラー
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

// getTransfer Download Transfer execute
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
	// withProgress
	if totalBytes != nil {
		f, _ := remoteFile.Stat()
		*totalBytes = f.Size()
	}
	if tfBytes != nil {
		remoteFileWithProgress := &IOReaderProgress{Reader: remoteFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.Copy(localFile, remoteFileWithProgress)
	} else {
		bytes, copyErr = io.Copy(localFile, remoteFile)
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

// UploadFile is Single File Upload
func (client *Client) UploadFile(localFilepath string, remoteFilepath string) (int64, error) {
	return putTransfer(client.SFTPClient, localFilepath, remoteFilepath, nil, nil)
}

// UploadFileWithProgress [Experimental] Put with Display Processing Bytes
func (client *Client) UploadFileWithProgress(localFilepath string, remoteFilepath string, transferred *int64, total *int64) (int64, error) {
	return putTransfer(client.SFTPClient, localFilepath, remoteFilepath, transferred, total)
}

// PutRecursively is Recursively Upload entire directories
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

// Upload Transfer execute
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
	// withProgress
	if totalBytes != nil {
		f, _ := localFile.Stat()
		*totalBytes = f.Size()
	}
	if tfBytes != nil {
		localFileWithProgress := &IOReaderProgress{Reader: localFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.Copy(remoteFile, localFileWithProgress)
	} else {
		bytes, copyErr = io.Copy(remoteFile, localFile)
	}

	if copyErr != nil {
		return 0, errors.New("copyErr: " + copyErr.Error())
	}
	return bytes, nil
}

// Close Connection ALL Connection Close
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

func getRecursivelyPath(localPath string, remotePath string, remoteFullFilepath string) (string, error) {
	rel, err := filepath.Rel(filepath.Clean(remotePath), remoteFullFilepath)
	if err != nil {
		return "", err
	}
	localFilepath := filepath.Join(localPath, rel)
	return filepath.ToSlash(localFilepath), nil
}

// Quit alias Close()
func (client *Client) Quit() (errors []error) {
	return client.Close()
}

// Chmod ..
func (client *Client) Chmod(path string, mode os.FileMode) error {
	return client.SFTPClient.Chmod(path, mode)
}

// Chown changes the user and group owners of the named file.
func (client *Client) Chown(path string, uid, gid int) error {
	return client.SFTPClient.Chown(path, uid, gid)
}

// Remove removes the specified file or directory. An error will be returned if no
// file or directory with the specified path exists, or if the specified directory
// is not empty.
func (client *Client) Remove(path string) error {
	return client.SFTPClient.Remove(path)
}

// RemoveDirectory removes a directory path.
func (client *Client) RemoveDirectory(path string) error {
	return client.SFTPClient.RemoveDirectory(path)
}

// Rename renames a file.
func (client *Client) Rename(oldname, newname string) error {
	return client.SFTPClient.Rename(oldname, newname)
}

// PosixRename renames a file using the posix-rename@openssh.com extension
// which will replace newname if it already exists.
func (client *Client) PosixRename(oldname, newname string) error {
	return client.SFTPClient.PosixRename(oldname, newname)
}

// Getwd returns the current working directory of the server. Operations
// involving relative paths will be based at this location.
func (client *Client) Getwd() (string, error) {
	return client.SFTPClient.Getwd()
}

// Mkdir creates the specified directory. An error will be returned if a file or
// directory with the specified path already exists, or if the directory's
// parent folder does not exist (the method cannot create complete paths).
func (client *Client) Mkdir(path string) error {
	return client.SFTPClient.Mkdir(path)
}

// MkdirAll creates a directory named path, along with any necessary parents,
// and returns nil, or else returns an error.
// If path is already a directory, MkdirAll does nothing and returns nil.
// If path contains a regular file, an error is returned
func (client *Client) MkdirAll(path string) error {
	return client.SFTPClient.MkdirAll(path)
}
