// Package sftp 提供基于 github.com/pkg/sftp 与 golang.org/x/crypto/ssh 的 SFTP 客户端封装。
//
// Config 配置主机、端口、账号与超时，New 建立 SSH 连接并打开 SFTP 通道；Client 提供单文件与
// 递归目录的上传下载、传输进度统计以及常用远端文件系统操作。
package sftp
