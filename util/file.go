package util

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

//////////// Directory/Filesystem Functions ////////////

// Stat stat()
func Stat(filename string) (os.FileInfo, error) {
	return os.Stat(filename)
}

// PathInfo Pathinfo pathinfo()
// -1: all; 1: dirname; 2: basename; 4: extension; 8: filename
// Usage:
// Pathinfo("/home/go/path/src/php2go/php2go.go", 1|2|4|8)
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
			basename, _ = info["basename"]
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

// FileExists file_exists()
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsFile is_file()
func IsFile(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// IsDir is_dir()
func IsDir(filename string) (bool, error) {
	fd, err := os.Stat(filename)
	if err != nil {
		return false, err
	}
	fm := fd.Mode()
	return fm.IsDir(), nil
}

// FileSize filesize()
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}

// FilePutContents file_put_contents()
func FilePutContents(filename string, data string, mode os.FileMode) error {
	return os.WriteFile(filename, []byte(data), mode)
}

// FileGetContents file_get_contents()
func FileGetContents(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	return string(data), err
}

// Unlink unlink()
func Unlink(filename string) error {
	return os.Remove(filename)
}

// Delete delete()
func Delete(filename string) error {
	return os.Remove(filename)
}

// Copy copy()
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

// IsReadable is_readable()
func IsReadable(filename string) bool {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return false
	}
	syscall.Close(fd)
	return true
}

// IsWriteable is_writeable()
func IsWriteable(filename string) bool {
	fd, err := syscall.Open(filename, syscall.O_WRONLY, 0)
	if err != nil {
		return false
	}
	syscall.Close(fd)
	return true
}

// Rename rename()
func Rename(oldname, newname string) error {
	return os.Rename(oldname, newname)
}

// Touch touch()
func Touch(filename string) (bool, error) {
	fd, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, err
	}
	fd.Close()
	return true, nil
}

// Mkdir mkdir()
func Mkdir(filename string, mode os.FileMode) error {
	return os.Mkdir(filename, mode)
}

// GetCwd Getcwd getcwd()
func GetCwd() (string, error) {
	dir, err := os.Getwd()
	return dir, err
}

// RealPath Realpath realpath()
func RealPath(path string) (string, error) {
	return filepath.Abs(path)
}

// BaseName Basename basename()
func BaseName(path string) string {
	return filepath.Base(path)
}

// Chmod chmod()
func Chmod(filename string, mode os.FileMode) bool {
	return os.Chmod(filename, mode) == nil
}

// Chown chown()
func Chown(filename string, uid, gid int) bool {
	return os.Chown(filename, uid, gid) == nil
}

// Fclose fclose()
func Fclose(handle *os.File) error {
	return handle.Close()
}

// FileMTime Filemtime filemtime()
func FileMTime(filename string) (int64, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer fd.Close()
	fileinfo, err := fd.Stat()
	if err != nil {
		return 0, err
	}
	return fileinfo.ModTime().Unix(), nil
}

// FGetCsv Fgetcsv fgetcsv()
func FGetCsv(handle *os.File, length int, delimiter rune) ([][]string, error) {
	reader := csv.NewReader(handle)
	reader.Comma = delimiter
	// TODO length limit
	return reader.ReadAll()
}

// Glob glob()
func Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}

// MkdirAll 创建文件夹
func MkdirAll(filename string, mode os.FileMode) error {
	return os.MkdirAll(filename, mode)
}
