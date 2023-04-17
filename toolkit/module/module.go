package module

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	// CmdNew represents the new command.
	CmdNew = &cobra.Command{
		Use:   "mp",
		Short: "数据模型层",
		Long:  "创建数据库模型层: toolkit mp dir table_name",
		Run:   Run,
	}
)

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
	if err := base.InitBase(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	// 创建文件夹
	dirModule := path.Join(base.Path, "module")
	_ = os.MkdirAll(dirModule, os.ModePerm)
	// 创建文件夹
	dirProto := path.Join(base.Path, "proto")
	_ = os.MkdirAll(dirProto, os.ModePerm)

	//创建数据
	dir := ""
	tableName := ""
	page := ""
	if len(args) == 0 {
		if err := survey.AskOne(&survey.Input{
			Message: "模型路径",
			Help:    "文件夹路径",
		}, &dir); err != nil || dir == "" {
			return
		}
		if err := survey.AskOne(&survey.Input{
			Message: "表名称",
			Help:    "数据库中某个表名称",
		}, &tableName); err != nil || tableName == "" {
			return
		}
		if err := survey.AskOne(&survey.Select{
			Message: "列表分页",
			Help:    "列表分页",
			Options: []string{"yes", "no"},
		}, &page); err != nil || page == "" {
			return
		}
	} else {
		dir = args[0]
		tableName = args[1]
		page = args[2]
	}
	if tableName == "" || dir == "" || page == "" {
		fmt.Println("初始化:", color.RedString("模型层名称 & 表名称 && 分页 必须填写"))
		return
	}
	// 文件夹的路径
	fullModuleDir := path.Join(base.Path, "module", dir)
	_ = os.MkdirAll(fullModuleDir, os.ModePerm)

	// 文件夹的路径
	fullProtoDir := path.Join(base.Path, "proto", dir)
	_ = os.MkdirAll(fullProtoDir, os.ModePerm)

	// 文件夹的路径
	fullServiceDir := path.Join(base.Path, "service", dir)
	_ = os.MkdirAll(fullServiceDir, os.ModePerm)

	// 初始化上下文
	timeout := "60s"
	t, err := time.ParseDuration(timeout)
	if err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	// 表结构信息
	tableColumn, err := base.TableColumn(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表结构信息:", color.RedString(err.Error()))
		return
	}
	tableColumnMap := make(map[string]base.Column)
	for _, item := range tableColumn {
		tableColumnMap[item.ColumnName] = item
	}
	// 表信息
	tableItem, err := base.TableItem(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表信息:", color.RedString(err.Error()))
		return
	}
	// 表索引
	tableIndex, err := base.TableIndex(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表索引:", color.RedString(err.Error()))
		return
	}
	// 表主键
	tablePrimary, err := base.TablePrimary(ctx, base.Config.String("db.Database"), tableName)
	if err != nil {
		fmt.Println("表主键:", color.RedString(err.Error()))
		return
	}

	//获取 go.mod
	mod, err := base.ModulePath(path.Join(base.Path, "go.mod"))
	if err != nil {
		fmt.Println("go.mod文件不存在:", color.RedString(err.Error()))
		return
	}
	// 数字长度
	strLen := strings.LastIndex(dir, "/")

	var methodList []base.Method

	var pageBool bool

	if page == "yes" {
		pageBool = true
	}

	//获取的索引信息没有
	if err != nil {
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           "List",
			Default:        true,
			Condition:      nil,
			ConditionTotal: 0,
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
			Page:           pageBool,
		}
		methodList = append(methodList, method)
	} else {
		//存储条件信息
		field := make([]string, 0)
		//有索引信息
		for _, v := range tableIndex {
			//查询条件
			condition := make([]base.Column, 0)
			//数据库索引
			indexField := v.Field
			indexFieldSlice := util.Explode(",", indexField)
			for _, fieldValue := range indexFieldSlice {
				//构造查询条件
				positionIndex := cast.ToInt64(len(condition)) + 1
				currentColumn := tableColumnMap[fieldValue]
				currentColumn.PosiTion = positionIndex
				condition = append(condition, currentColumn)
				if !util.InArray(fieldValue, field) {
					field = append(field, fieldValue)
				}
			}
			// 数据库中的索引名称
			indexName := v.IndexName
			// 拆分字符串,得到索引类型和索引名称
			indexNameSlice := util.Explode(":", indexName)
			if len(indexNameSlice) < 2 {
				continue
			}
			// 自定义函数名称和索引信息
			customIndexType := util.UCWords(indexNameSlice[0])
			customIndexName := util.UCWords(indexNameSlice[1])
			method := base.Method{
				Table:          tableItem,
				TableColumn:    tableColumn,
				Type:           customIndexType,
				Name:           customIndexName,
				Default:        false,
				Condition:      condition,
				ConditionTotal: len(condition),
				Primary:        tablePrimary,
				Pkg:            dir[strLen+1:],
				PkgPath:        dir,
				ModName:        mod,
				Page:           pageBool,
			}
			//添加到集合中
			methodList = append(methodList, method)
		}
		condition := make([]base.Column, 0)
		for _, fieldValue := range field {
			//构造查询条件
			positionIndex := cast.ToInt64(len(condition)) + 1
			currentColumn := tableColumnMap[fieldValue]
			currentColumn.PosiTion = positionIndex
			condition = append(condition, currentColumn)
			// condition = append(condition, tableColumnMap[fieldValue])
		}
		method := base.Method{
			Table:          tableItem,
			TableColumn:    tableColumn,
			Type:           "List",
			Name:           "List",
			Default:        true,
			Condition:      condition,
			ConditionTotal: len(condition),
			Primary:        tablePrimary,
			Pkg:            dir[strLen+1:],
			PkgPath:        dir,
			ModName:        mod,
			Page:           pageBool,
		}
		methodList = append(methodList, method)
	}
	// 数据模型
	moduleParam := base.ModuleParam{
		Pkg:         dir[strLen+1:],
		PkgPath:     dir,
		ModName:     mod,
		Primary:     tablePrimary,
		Table:       tableItem,
		TableColumn: tableColumn,
		Method:      methodList,
		Page:        pageBool,
	}
	GenerateModule(moduleParam, fullModuleDir, tableName)
	GenerateProto(moduleParam, fullProtoDir, fullServiceDir, tableName)
	GenerateService(moduleParam, fullServiceDir, tableName)

	strLenSlot := util.Explode("/", dir)
	pkgName := strLenSlot[len(strLenSlot)-1]
	if len(strLenSlot) > 1 {
		pkgName = strLenSlot[len(strLenSlot)-2]
	}

	builder := strings.Builder{}
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("package %s", pkgName))
	builder.WriteString("\n")
	builder.WriteString("import (")
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("	\"%s/service/%s\"", moduleParam.ModName, moduleParam.PkgPath))
	builder.WriteString("\n")
	builder.WriteString("\n")
	builder.WriteString("	\"github.com/abulo/ratel/v3/server/xgrpc\"")
	builder.WriteString("\n")
	builder.WriteString(")")
	builder.WriteString("\n")
	builder.WriteString("func Registry(server *xgrpc.Server) {")
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("	// %s->%s", moduleParam.Table.TableComment, moduleParam.Table.TableName))
	builder.WriteString("\n")
	builder.WriteString(fmt.Sprintf("	%s.Register%sServiceServer(server.Server, &%s.Srv%sServiceServer{", moduleParam.Pkg, base.CamelStr(moduleParam.Table.TableName), moduleParam.Pkg, base.CamelStr(moduleParam.Table.TableName)))
	builder.WriteString("\n")
	builder.WriteString("		Server: server,")
	builder.WriteString("\n")
	builder.WriteString("	})")
	builder.WriteString("\n")
	builder.WriteString("}")
	fmt.Println(builder.String())
}
