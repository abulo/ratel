package util

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// Stat 返回 filename 指向条目的 FileInfo；条目不存在或状态获取失败时返回错误。
func Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// PathInfo 解析路径并按 options 位掩码返回组成部分。
//
// 位掩码取值：1 为目录 dirname，2 为文件名 basename，4 为扩展名 extension，
// 8 为不含扩展名的文件名 filename；options 为 -1 时返回全部四项。
func PathInfo(path string, options int) map[string]string {
	if options == -1 {
		options = 1 | 2 | 4 | 8
	}
	info := make(map[string]string)
	if (options & 1) == 1 {
		info["dirname"] = filepath.Dir(path)
	}
	if (options & 2) == 2 {
		info["basename"] = filepath.Base(path)
	}
	if ((options & 4) == 4) || ((options & 8) == 8) {
		basename := ""
		if (options & 2) == 2 {
			if val, ok := info["basename"]; !ok {
				basename = val
			}
		} else {
			basename = filepath.Base(path)
		}
		p := strings.LastIndex(basename, ".")
		filename, extension := "", ""
		if p > 0 {
			filename, extension = basename[:p], basename[p+1:]
		} else if p == -1 {
			filename = basename
		} else if p == 0 {
			extension = basename[p+1:]
		}
		if (options & 4) == 4 {
			info["extension"] = extension
		}
		if (options & 8) == 8 {
			info["filename"] = filename
		}
	}
	return info
}

// FileExists 判断文件或目录是否存在；状态获取出错（非不存在错误）时同样返回 true。
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFile 判断路径指向的条目是否存在。
//
// 当前实现仅检查条目是否存在，不区分文件与目录。
func IsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDir 判断路径是否为目录；状态获取失败时返回 false 与错误。
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

// FileSize 返回文件大小（字节）；状态获取失败时返回 0 与错误。
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}

// FilePutContents 将 data 以 mode 权限模式写入 filename，覆盖已有内容；失败返回错误。
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return os.WriteFile(filename, []byte(data), mode)
}

// FileGetContents 读取整个文件内容并以字符串返回；读取失败时返回空字符串与错误。
func FileGetContents(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	return string(data), err
}

// Unlink 删除指定文件或空目录；失败返回错误。
func Unlink(filename string) error {
	return os.Remove(filename)
}

// Delete 删除指定文件或空目录；失败返回错误，语义与 Unlink 相同。
func Delete(filename string) error {
	return os.Remove(filename)
}

// Copy 将 source 复制到 dest；目标文件不存在时创建，打开或拷贝失败时返回 false 与错误。
func Copy(source, dest string) (bool, error) {
	fd1, err := os.Open(source)
	if err != nil {
		return false, err
	}
	defer fd1.Close()
	fd2, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer fd2.Close()
	_, e := io.Copy(fd2, fd1)
	if e != nil {
		return false, e
	}
	return true, nil
}

// IsReadable 尝试以只读方式打开文件以判断是否可读；无法打开时返回 false。
func IsReadable(filename string) bool {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}
	syscall.Close(fd)
	return true
}

// IsWriteable 尝试以只写方式打开文件以判断是否可写；无法打开时返回 false。
func IsWriteable(filename string) bool {
	fd, err := syscall.Open(filename, syscall.O_WRONLY, 0)
	if err != nil {
		return false
	}
	syscall.Close(fd)
	return true
}

// Rename 将 oldname 重命名或移动为 newname；失败返回错误。
func Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Touch 创建不存在的文件（权限 0666）；已存在时仅打开不修改访问时间。失败返回 false 与错误。
func Touch(filename string) (bool, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	fd.Close()
	return true, nil
}

// Mkdir 按 mode 权限模式创建单级目录，父目录不存在时返回错误。
func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

// GetCwd 返回当前工作目录的绝对路径；获取失败时返回错误。
func GetCwd() (string, error) {
	dir, err := os.Getwd()
	return dir, err
}

// RealPath 校验路径存在后返回其绝对路径；路径不存在或状态获取失败时返回错误。
func RealPath(path string) (string, error) {
	_, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return filepath.Abs(path)
}

// BaseName 返回路径的最后一段（文件名或目录名）。
func BaseName(path string) string {
	return filepath.Base(path)
}

// Chmod 修改 filename 的权限模式，操作成功返回 true。
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown 修改 filename 的属主 uid 与属组 gid，操作成功返回 true。
func Chown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// Fclose 关闭已打开的文件句柄 handle；关闭失败返回错误。
func Fclose(handle *os.File) error {
	return handle.Close()
}

// FileMTime 返回文件最后修改时间的 Unix 时间戳（秒）；打开或获取状态失败时返回 0 与错误。
func FileMTime(filename string) (int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	fi, err := fd.Stat()
	if err != nil {
		return 0, err
	}
	return fi.ModTime().Unix(), nil
}

// FGetCsv 以 delimiter 为分隔符读取整个 CSV 文件，返回二维字符串数组；读取失败返回错误。
func FGetCsv(handle *os.File, length int, delimiter rune) ([][]string, error) {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter
	// TODO(owner): 尚未实现 length 参数对应的单行字段长度限制。
	return reader.ReadAll()
}

// Glob 返回匹配 pattern 的全部文件路径；pattern 非法时返回错误。
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// MkdirAll 递归创建目录及其父目录，已存在的目录不报错；失败返回错误。
func MkdirAll(filename string, mode os.FileMode) error {
	return os.MkdirAll(filename, mode)
}
