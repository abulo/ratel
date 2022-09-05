package module

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/abulo/ratel/v3/util"
)

func Run(tableName, outputDir, outputPackage, dao, tpl string) {
	_ = os.MkdirAll(outputDir, os.ModePerm)

	content, _ := util.FileGetContents(tpl)

	content = util.StrReplace("{{Package}}", outputPackage, content, -1)
	content = util.StrReplace("{{Dao}}", dao, content, -1)
	content = util.StrReplace("{{TableName}}", tableName, content, -1)
	builder := strings.Builder{}
	//转换表名
	builder.Reset()

	builder.WriteString(content)
	fileStr := builder.String()
	_ = os.WriteFile(path.Join(outputDir, tableName+".go"), []byte(fileStr), os.ModePerm)

	_ = os.Chdir(outputDir)
	cmd := exec.Command("go", "fmt")
	out, e := cmd.CombinedOutput()
	if e != nil {
		panic(e)
	}

	cmdImport := exec.Command("goimports", "-w", path.Join(outputDir, "*.go"))
	outImport, eI := cmdImport.CombinedOutput()
	if eI != nil {
		panic(eI)
	}
	fmt.Printf("goimports结果:\n%s\n", string(outImport))

	fmt.Printf("格式化结果:\n%s\n", string(out))
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}
