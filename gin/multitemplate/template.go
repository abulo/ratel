package multitemplate

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/abulo/ratel/gin"
)

func getFilelist(path string, stuffix string) (files []string) {
	// 遍历目录
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		// 将模板后缀的文件放到列表
		if strings.HasSuffix(path, stuffix) {
			files = append(files, path)
		}
		return nil
	})
	return
}

// LoadTemplateFiles 加载模板
func LoadTemplateFiles(templateDir, stuffix string, funcMap template.FuncMap) Renderer {

	if gin.IsDebugging() {
		fmt.Println("========================")
		fmt.Println("[debug] loading template files ")
		fmt.Println("========================")
	}
	r := NewRenderer()
	rd, _ := ioutil.ReadDir(templateDir)
	for _, fi := range rd {
		if fi.IsDir() {
			// 如果是目录
			for _, f := range getFilelist(path.Join(templateDir, fi.Name()), stuffix) {
				// 添加到模板的时候，去掉跟路径
				if len(funcMap) > 0 {
					r.AddFromFilesFuncs(f[len(templateDir)+1:], funcMap, f)
				} else {
					r.AddFromFiles(f[len(templateDir)+1:], f)
				}
				if gin.IsDebugging() {
					fmt.Println(f[len(templateDir)+1:])
				}
			}
		} else {
			if strings.HasSuffix(fi.Name(), stuffix) {
				// 如果再根目录底下的文件直接添加
				if len(funcMap) > 0 {
					r.AddFromFilesFuncs(fi.Name(), funcMap, path.Join(templateDir, fi.Name()))
				} else {
					r.AddFromFiles(fi.Name(), path.Join(templateDir, fi.Name()))
				}
				if gin.IsDebugging() {
					fmt.Println(fi.Name())
				}
			}
		}
	}
	if gin.IsDebugging() {
		fmt.Println("========================")
	}
	return r
}
