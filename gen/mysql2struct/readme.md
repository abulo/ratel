使用方法


```golang

dbInfo := &mysql2struct.DBInfo{
    DBType:   dbType,
    Host:     host,
    UserName: username,
    Password: password,
    Charset:  charset,
}
dbModel := mysql2struct.NewDBModel(dbInfo)
err := dbModel.Connect()
if err != nil {
    log.Fatalf("dbModel.Connect err: %v", err)
}
columns, err := dbModel.GetColumns(dbName, tableName)
if err != nil {
    log.Fatalf("dbModel.GetColumns err:%v", err)
}

template := mysql2struct.NewStructTemplate()
templateColumns := template.AssemblyColumns(columns)
err = template.Generate(tableName, templateColumns)
if err != nil {
    log.Fatalf("template.Generate err: %v", err)
}


```