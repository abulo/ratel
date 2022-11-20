package backstage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/config"
	"github.com/abulo/ratel/v3/config/toml"
	"github.com/abulo/ratel/v3/stores/mysql"
	"github.com/abulo/ratel/v3/stores/query"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var (
	CmdNew = &cobra.Command{
		Use:   "backstage",
		Short: "管理端接口",
		Long:  "管理端接口: toolkit backstage dir dao ",
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

	tableName := ""
	alias := ""
	moduleDir := ""
	view := ""
	viewDir := ""
	proto := ""
	if err = survey.AskOne(&survey.Input{
		Message: "表名称",
		Help:    "数据库中某个表名称",
	}, &tableName); err != nil || tableName == "" {
		return
	}
	if err = survey.AskOne(&survey.Input{
		Message: "别名",
		Help:    "别名",
	}, &alias); err != nil || alias == "" {
		return
	}
	if err = survey.AskOne(&survey.Input{
		Message: "模块的文件夹",
		Help:    "模块的文件夹路径",
	}, &moduleDir); err != nil || moduleDir == "" {
		return
	}
	if err = survey.AskOne(&survey.Select{
		Message: "gRPC协议",
		Help:    "是否使用gRPC协议",
		Options: []string{"yes", "no"},
	}, &proto); err != nil || proto == "" {
		return
	}
	if err = survey.AskOne(&survey.Select{
		Message: "管理视图",
		Help:    "是否生成管理视图",
		Options: []string{"yes", "no"},
	}, &view); err != nil || view == "" {
		return
	}
	if view == "yes" {
		if err = survey.AskOne(&survey.Input{
			Message: "视图的文件夹",
			Help:    "视图的文件夹路径",
		}, &viewDir); err != nil || viewDir == "" {
			return
		}
	}

	mod, err := base.ModulePath(path.Join(wd, "go.mod"))
	if err != nil {
		fmt.Println("go.mod 文件不存在")
		return
	}

	ColumnInfo, err := QueryColumn(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表字段信息不准确")
		return
	}
	ColumnMap := make(map[string]Column)

	for _, v := range ColumnInfo {
		ColumnMap[v.ColumnName] = v
	}

	tableInfo, err := QueryTable(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		fmt.Println("表信息不准确")
		return
	}

	//获取表信息
	functionList := make([]Function, 0)
	fieldList := make([]string, 0)
	indexList, err := QueryIndex(ctx, AppConfig.String("mysql.Database"), tableName)
	if err != nil {
		function := Function{}
		function.Type = "list"
		function.Name = CamelStr("list")
		function.TableName = tableName
		function.Mark = CamelStr(tableName)
		function.Default = true
		argument := make([]Column, 0)
		if len(fieldList) > 0 {
			for _, field := range fieldList {
				arg := util.Explode(":", field)
				fieldKey := arg[0]
				argument = append(argument, ColumnMap[fieldKey])
			}
		}
		function.Column = argument
		functionList = append(functionList, function)
	} else {

		for _, index := range indexList {
			// index.IndexName
			indexName := util.Explode(":", index.IndexName)
			if len(indexName) < 2 {
				continue
			}
			function := Function{}
			function.Type = indexName[0]
			function.Name = CamelStr(indexName[1])
			function.TableName = tableName
			function.Mark = CamelStr(tableName)
			function.Default = false
			//获取参数
			fields := util.Explode(",", index.Field)
			for _, field := range fields {
				if !util.InArray(field, fieldList) {
					fieldList = append(fieldList, field)
				}
			}

			argument := make([]Column, 0)
			for _, field := range fields {
				arg := util.Explode(":", field)
				fieldKey := arg[0]
				argument = append(argument, ColumnMap[fieldKey])
			}
			function.Column = argument
			functionList = append(functionList, function)
		}

		function := Function{}
		function.Type = "list"
		function.Name = CamelStr("list")
		function.TableName = tableName
		function.Mark = CamelStr(tableName)
		function.Default = true
		argument := make([]Column, 0)
		if len(fieldList) > 0 {
			for _, field := range fieldList {
				arg := util.Explode(":", field)
				fieldKey := arg[0]
				argument = append(argument, ColumnMap[fieldKey])
			}
		}
		function.Column = argument
		functionList = append(functionList, function)
	}

	res := BackStage{
		TableName:  tableName,
		Mark:       CamelStr(tableName),
		Alias:      alias,
		Column:     ColumnInfo,
		Table:      tableInfo,
		PrimaryKey: tableName + "_id",
		Mod:        mod,
	}
	n := strings.LastIndex(moduleDir, "/")
	res.Pkg = moduleDir[n+1:]
	res.FunctionList = functionList
	if proto == "yes" {
		res.Proto = true
	}
	if view == "yes" {
		res.View = true
	}
	res.ViewDir = viewDir

	jsonString, _ := json.Marshal(res)
	fmt.Println(cast.ToString(jsonString))

}

var DataTypeMap = map[string][]string{
	//整型
	"TINYINT":   {"int64", "query.NullInt64"},
	"SMALLINT":  {"int64", "query.NullInt64"},
	"MEDIUMINT": {"int64", "query.NullInt64"},
	"INT":       {"int64", "query.NullInt64"},
	"INTEGER":   {"int64", "query.NullInt64"},
	"BIGINT":    {"int64", "query.NullInt64"},
	//浮点数
	"FLOAT":   {"float64", "query.NullFloat64"},
	"DOUBLE":  {"float64", "query.NullFloat64"},
	"DECIMAL": {"float64", "query.NullFloat64"},
	//时间
	"DATE":      {"query.NullDate", "query.NullDate"},
	"TIME":      {"query.NullTime", "query.NullTime"},
	"YEAR":      {"query.NullYear", "query.NullYear"},
	"DATETIME":  {"query.NullDateTime", "query.NullDateTime"},
	"TIMESTAMP": {"query.NullTimeStamp", "query.NullTimeStamp"},
	//字符串
	"CHAR":       {"string", "query.NullString"},
	"VARCHAR":    {"string", "query.NullString"},
	"TINYBLOB":   {"string", "query.NullString"},
	"TINYTEXT":   {"string", "query.NullString"},
	"BLOB":       {"string", "query.NullString"},
	"TEXT":       {"string", "query.NullString"},
	"MEDIUMBLOB": {"string", "query.NullString"},
	"MEDIUMTEXT": {"string", "query.NullString"},
	"LONGBLOB":   {"string", "query.NullString"},
	"LONGTEXT":   {"string", "query.NullString"},
	"JSON":       {"string", "query.NullString"},
}

type BackStage struct {
	Proto        bool       //是否采用微服务模式
	View         bool       //是否包含视图
	ViewDir      string     //视图的文件夹
	Pkg          string     //接口所在包名
	TableName    string     //表名名
	Mark         string     //
	Table        Table      //表信息
	Alias        string     //别名
	Column       []Column   //表信息
	PrimaryKey   string     //主键
	FunctionList []Function //函数
	Mod          string     //模型名称
}

// Column 表明信息
type Column struct {
	ColumnName    string `db:"COLUMN_NAME"`
	IsNullable    string `db:"IS_NULLABLE"`
	DataType      string `db:"DATA_TYPE"`
	ColumnKey     string `db:"COLUMN_KEY"`
	ColumnComment string `db:"COLUMN_COMMENT"`
}

type Table struct {
	TableName    string `db:"TABLE_NAME"`
	TableComment string `db:"TABLE_COMMENT"`
	Mark         string
}

type Index struct {
	IndexName string `db:"INDEX_NAME"`
	Field     string `db:"FIELD"`
}

type Function struct {
	Type      string
	Name      string
	Column    []Column
	TableName string
	Mark      string
	Default   bool
}

// QueryColumn 获取数据中表中字段的信息
func QueryColumn(ctx context.Context, DbName, TableName string) ([]Column, error) {
	var res []Column
	builder := Link.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Rows().ToStruct(&res)
	return res, err
}

// QueryTable 获取数据中表的信息
func QueryTable(ctx context.Context, DbName string, TableName string) (Table, error) {
	var res Table
	builder := Link.NewBuilder(ctx).Select("TABLE_NAME", "TABLE_COMMENT").Table("information_schema.TABLES").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName)
	err := builder.Row().ToStruct(&res)
	if err == nil {
		res.Mark = CamelStr(TableName)
	}
	return res, err
}

func QueryIndex(ctx context.Context, DbName, TableName string) ([]Index, error) {
	var res []Index
	err := Link.NewBuilder(ctx).Select("statistics.INDEX_NAME", "GROUP_CONCAT(CONCAT(statistics.COLUMN_NAME,':',`columns`.DATA_TYPE )) AS FIELD").Table("`information_schema`.`STATISTICS` AS statistics").LeftJoin("information_schema.`COLUMNS` AS `columns`", "statistics.COLUMN_NAME = `columns`.COLUMN_NAME").Where("statistics.TABLE_SCHEMA", DbName).Where("statistics.TABLE_NAME", TableName).Where("`columns`.TABLE_SCHEMA", DbName).Where("`columns`.TABLE_NAME", TableName).NotEqual("statistics.INDEX_NAME", "PRIMARY").GroupBy("statistics.TABLE_NAME", "statistics.INDEX_NAME").OrderBy("statistics.NON_UNIQUE", query.ASC).OrderBy("statistics.SEQ_IN_INDEX", query.ASC).Rows().ToStruct(&res)
	return res, err
}

// QueryColumn 获取数据中表中字段的信息
func QueryColumnPrimaryKey(ctx context.Context, DbName, TableName string) (Column, error) {
	var res Column
	builder := Link.NewBuilder(ctx).Select("COLUMN_NAME", "IS_NULLABLE", "DATA_TYPE", "COLUMN_KEY", "COLUMN_COMMENT").Table("information_schema.COLUMNS").Where("TABLE_SCHEMA", DbName).Where("TABLE_NAME", TableName).Where("COLUMN_KEY", "PRI").OrderBy("ORDINAL_POSITION", query.ASC)
	err := builder.Row().ToStruct(&res)
	return res, err
}

// CamelStr 下划线转驼峰
func CamelStr(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = util.UCWords(name)
	return strings.Replace(name, " ", "", -1)
}

func Helper(name string) string {
	name = CamelStr(name)
	return strings.ToLower(string(name[0])) + name[1:]
}

func Convert(column Column, name string) string {
	dataType := strings.ToUpper(column.DataType)
	value, ok := DataTypeMap[dataType]
	if ok {
		if column.IsNullable == "YES" {
			dataType = value[1]
		} else {
			dataType = value[0]
		}
	} else {
		dataType = "string"
	}
	var res string
	switch dataType {
	case "string":
		res = "cast.ToString(ctx.PostForm(\"" + Helper(name) + "\"))"
	case "query.NullString":
		res = "query.NewString(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "int64":
		res = "cast.ToInt64(ctx.PostForm(\"" + Helper(name) + "\"))"
	case "query.NullInt64":
		res = "query.NewInt64(cast.ToInt64(ctx.PostForm(\"" + Helper(name) + "\"))"
	case "float64":
		res = "cast.ToFloat64(ctx.PostForm(\"" + Helper(name) + "\"))"
	case "query.NullFloat64":
		res = "query.NewFloat64(cast.ToFloat64(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "query.NullDate":
		res = "query.NewDate(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "query.NullTime":
		res = "query.NewTime(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "query.NullYear":
		res = "query.NewYear(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "query.NullDateTime":
		res = "query.NewDateTime(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	case "query.NullTimeStamp":
		res = "query.NewTimeStamp(cast.ToString(ctx.PostForm(\"" + Helper(name) + "\")))"
	}
	return res
}

var BackStageTpl = `
package {{.Pkg}}

// {{.Alias}} 模块操作



`
