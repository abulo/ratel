// Package ftp 提供基于 github.com/jlaffaye/ftp 的 FTP 客户端封装。
//
// Config 配置主机、端口、账号与连接超时，New 完成拨号；Client 在底层连接之上提供目录切换、
// 文件与目录的增删查、上传下载以及目录列举等常用操作。
package ftp
