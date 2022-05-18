package mysql2struct

type Config struct {
	// Username      string `json:"username"`
	// Password      string `json:"password"`
	// Protocol      string `json:"protocol"`
	// Address       string `json:"address"`
	// Dbname        string `json:"dbname"`
	// TableName     string `json:"tableName"`
	OutputDir     string `json:"outputDir"`
	OutputPackage string `json:"outputPackage"`
}

// {
// 	"username": "用户名",
// 	"password": "密码",
// 	"protocol": "tcp",
// 	"address": "127.0.0.1:3306",
// 	"dbname": "数据库名称",
// 	"tableName": "表名(可空)",
// 	"outputDir": "输出目录",
// 	"outputPackage": "struct文件的包名"
//   }
