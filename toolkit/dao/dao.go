package dao

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/toml"
	"github.com/abulo/ratel/v3/stores/mysql"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	// CmdNew represents the new command.
	CmdNew = &cobra.Command{
		Use:   "dao",
		Short: "Create a dao",
		Long:  "Create a dao using the repository template. Example: ratel dao",
		Run:   run,
	}
	AppConfig *config.Config
	Link      *query.Query
)

func run(cmd *cobra.Command, args []string) {
	timeout := "60s"
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	mysqlConfig := "mysql.toml"
	configFile := wd + "/" + mysqlConfig
	if !util.FileExists(configFile) {
		fmt.Println("The mysql configuration file does not exist.")
		return
	}

	daoName := "dao"
	dir := wd + "/" + daoName
	_ = os.MkdirAll(dir, os.ModePerm)
	//加载配置文件
	AppConfig = config.New("dao")
	AppConfig.AddDriver(toml.Driver)
	AppConfig.LoadFiles(configFile)
	//创建数据链接
	opt := &mysql.Config{}

	if Username := cast.ToString(AppConfig.String("mysql.Username")); Username != "" {
		opt.Username = Username
	}
	if Password := cast.ToString(AppConfig.String("mysql.Password")); Password != "" {
		opt.Password = Password
	}
	if Host := cast.ToString(AppConfig.String("mysql.Host")); Host != "" {
		opt.Host = Host
	}
	if Port := cast.ToString(AppConfig.String("mysql.Port")); Port != "" {
		opt.Port = Port
	}
	if Charset := cast.ToString(AppConfig.String("mysql.Charset")); Charset != "" {
		opt.Charset = Charset
	}
	if Database := cast.ToString(AppConfig.String("mysql.Database")); Database != "" {
		opt.Database = Database
	}

	// # MaxOpenConns 连接池最多同时打开的连接数
	// MaxOpenConns = 128
	// # MaxIdleConns 连接池里最大空闲连接数。必须要比maxOpenConns小
	// MaxIdleConns = 32
	// # MaxLifetime 连接池里面的连接最大存活时长(分钟)
	// MaxLifetime = 10
	// # MaxIdleTime 连接池里面的连接最大空闲时长(分钟)
	// MaxIdleTime = 5

	if MaxLifetime := cast.ToInt(AppConfig.Int("mysql.MaxLifetime")); MaxLifetime > 0 {
		opt.MaxLifetime = time.Duration(MaxLifetime) * time.Minute
	}
	if MaxIdleTime := cast.ToInt(AppConfig.Int("mysql.MaxIdleTime")); MaxIdleTime > 0 {
		opt.MaxIdleTime = time.Duration(MaxIdleTime) * time.Minute
	}
	if MaxIdleConns := cast.ToInt(AppConfig.Int("mysql.MaxIdleConns")); MaxIdleConns > 0 {
		opt.MaxIdleConns = cast.ToInt(MaxIdleConns)
	}
	if MaxOpenConns := cast.ToInt(AppConfig.Int("mysql.MaxOpenConns")); MaxOpenConns > 0 {
		opt.MaxOpenConns = cast.ToInt(MaxOpenConns)
	}
	opt.DriverName = "mysql"
	opt.DisableMetric = cast.ToBool(AppConfig.Bool("mysql.DisableMetric"))
	opt.DisableTrace = cast.ToBool(AppConfig.Bool("mysql.DisableTrace"))
	Link = mysql.NewClient(opt)
	//获取表信息
	tableList, err := QueryTable(ctx, AppConfig.String("mysql.Database"))
	if err != nil {
		fmt.Println("QueryTable is Error:", err)
		return
	}

	builder := strings.Builder{}
	for _, table := range tableList {
		columns, err := QueryColumn(ctx, AppConfig.String("mysql.Database"), table.TableName)
		if err != nil {
			continue
		}
		//转换表名
		builder.Reset()
		packageTime := false
		packageSQL := false
		builder.WriteString(fmt.Sprintf("//%s\ntype %s struct {\n", table.TableComment, CamelStr(table.TableName)))

		for _, column := range columns {
			//转换列名
			dataType := strings.ToUpper(column.DataType)
			value, ok := DataTypeMap[dataType]
			if ok {
				if column.IsNullable == "YES" {
					dataType = value[1]
					packageSQL = true
				} else {
					dataType = value[0]
				}
				//是否需要 sql 包
				if dataType == "time.Time" {
					packageTime = true
				}
			} else {
				dataType = "string"
			}

			//拼接字符串
			camelStr := CamelStr(column.ColumnName)
			builder.WriteString(fmt.Sprintf("	%s %s `db:\"%s\" json:\"%s\" form:\"%s\"` //%s", camelStr, dataType, column.ColumnName, strings.ToLower(string(camelStr[0]))+camelStr[1:], strings.ToLower(string(camelStr[0]))+camelStr[1:], column.ColumnComment))
			if column.ColumnKey != "" {
				builder.WriteString("(" + column.ColumnKey + ")")
			}
			builder.WriteString("\n")
		}
		builder.WriteString("}\n")
		fileStr := "package " + daoName
		fileStr += "\nimport ("
		if packageSQL {
			fileStr += "\"github.com/abulo/ratel/v3/stores/query\"\n\n"
		}
		if packageTime {
			fileStr += "\"time\"\n\n"
		}
		fileStr += ")\n"
		fileStr += builder.String()

		outFile := path.Join(dir, table.TableName+".go")
		if util.FileExists(outFile) {
			util.Delete(outFile)
		}
		if err := os.WriteFile(outFile, []byte(fileStr), os.ModePerm); err == nil {
			fmt.Printf("\n🍺 Create  "+dir+" %s\n", color.GreenString(dir+"/"+table.TableName+".go"))
		}
	}

	_ = os.Chdir(dir)
	cmdShell := exec.Command("go", "fmt")
	if _, err := cmdShell.CombinedOutput(); err != nil {
		fmt.Println("go fmt is Error:", err)
		return
	}
	cmdImport := exec.Command("goimports", "-w", path.Join(dir, "*.go"))
	cmdImport.CombinedOutput()
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}
