package ecode

// ErrKind* 为错误类别标签，写入日志字段以归类错误来源；同一次日志记录中类别应保持唯一。
const (
	// ErrKindUnmarshalConfigErr 为配置反序列化失败的类别标签。
	ErrKindUnmarshalConfigErr = "unmarshal config err"
	// ErrKindRegisterErr 为注册中心操作失败的类别标签。
	ErrKindRegisterErr = "register err"
	// ErrKindURIErr 为 URI 解析失败的类别标签。
	ErrKindURIErr = "uri err"
	// ErrKindRequestErr 为请求执行失败的类别标签。
	ErrKindRequestErr = "request err"
	// ErrKindFlagErr 为命令行参数解析失败的类别标签。
	ErrKindFlagErr = "flag err"
	// ErrKindListenErr 为端口监听失败的类别标签。
	ErrKindListenErr = "listen err"
	// ErrKindAny 为未归类错误使用的通用类别标签。
	ErrKindAny = "any"
)

// Msg* 为统一的日志与提示消息文本，模块内消息文本保持唯一。
const (
	// MsgRegisterParseConfigErr1 为注册配置解析失败的首条消息。
	MsgRegisterParseConfigErr1 = "register parse config err1"
	// MsgRegisterParseConfigErr2 为注册配置解析失败的第二条消息。
	MsgRegisterParseConfigErr2 = "register parse config err2"
	// MsgRegisterETCDErr 为注册服务到 Etcd 失败的消息。
	MsgRegisterETCDErr = "register service err"
	// MsgRegisterETCDOk 为注册服务到 Etcd 成功的消息。
	MsgRegisterETCDOk = "register service ok"
	// MsgDeleteParseProviderURIErr 为删除提供方 URI 解析失败的消息。
	MsgDeleteParseProviderURIErr = "delete parse provider uri err"
	// MsgDeleteParseConfiguratorsURIErr 为删除配置中心 URI 解析失败的消息。
	MsgDeleteParseConfiguratorsURIErr = "delete parse configurators uri err"
	// MsgUpdateParseProviderURIErr 为更新提供方 URI 解析失败的消息。
	MsgUpdateParseProviderURIErr = "update parse provider uri err"
	// MsgUpdateParseConfiguratorsURIErr 为更新配置中心 URI 解析失败的消息。
	MsgUpdateParseConfiguratorsURIErr = "update parse configurators uri err"
	// MsgUpdateParseConfigErr 为更新配置失败的消息。
	MsgUpdateParseConfigErr = "update parse config err"
	// MsgUpdateResolverOk 为解析器更新成功的消息。
	MsgUpdateResolverOk = "update resolver ok"
	// MsgWatchRequestErr 为监听请求失败的消息。
	MsgWatchRequestErr = "watch request err"
	// MsgDeregisterETCDOk 为从 Etcd 注销服务成功的消息。
	MsgDeregisterETCDOk = "deregister etcd service ok"
	// MsgDeregisterETCDErr 为从 Etcd 注销服务失败的消息。
	MsgDeregisterETCDErr = "deregister etcd service err"
	// MsgRegistryResolverOk 为解析器构建成功的消息。
	MsgRegistryResolverOk = "resolver build ok"
	// MsgRegistryResolverNow 为解析器当前状态的消息。
	MsgRegistryResolverNow = "resolver now"
	// MsgRegistryResolverClose 为解析器关闭的消息。
	MsgRegistryResolverClose = "resolver close"

	// 配置加载相关消息。
	// MsgConfigParseFlagPanic 为解析命令行参数触发 panic 的消息。
	MsgConfigParseFlagPanic = "parse flag panic"
	// MsgConfigLoadFromRemoteDataSourcePanic 为从远程数据源加载配置触发 panic 的消息。
	MsgConfigLoadFromRemoteDataSourcePanic = "load from remote datasource panic"
	// MsgConfigLoadFromRemoteDataSourceOK 为从远程数据源加载配置成功的消息。
	MsgConfigLoadFromRemoteDataSourceOK = "load from remote datasource ok"
	// MsgConfigLoadFromFilePanic 为从文件加载配置触发 panic 的消息。
	MsgConfigLoadFromFilePanic = "load from file panic"
	// MsgConfigLoadFromFileOK 为从文件加载配置成功的消息。
	MsgConfigLoadFromFileOK = "load from file ok"

	// 应用启动相关消息。
	// MsgAppStartServerOk 为服务端启动成功的消息。
	MsgAppStartServerOk = "server start"
	// MsgAppStartGovernorOk 为管控组件启动成功的消息。
	MsgAppStartGovernorOk = "governor start"

	// 进程资源限制相关消息。
	// MsgProcSetPanic 为设置最大处理器数触发 panic 的消息。
	MsgProcSetPanic = "set max procs panic"
	// MsgProcSetOk 为设置最大处理器数成功的消息。
	MsgProcSetOk = "set max procs ok"

	// gRPC 服务端相关消息。
	// MsgGrpcServerNewErr 为创建 gRPC 服务端失败的消息。
	MsgGrpcServerNewErr = "new grpc server err"
	// MsgGrpcServerRecover 为 gRPC 服务端捕获 panic 的消息。
	MsgGrpcServerRecover = "grpc server recover"

	// MySQL 客户端相关消息。
	// MsgClientMysqlOpenStart 为打开 MySQL 连接开始的消息。
	MsgClientMysqlOpenStart = "client mysql start"
	// MsgClientMysqlOpenPanic 为打开 MySQL 连接触发 panic 的消息。
	MsgClientMysqlOpenPanic = "mysql open panic"
	// MsgClientMysqlPingPanic 为 MySQL 连通性探测触发 panic 的消息。
	MsgClientMysqlPingPanic = "mysql ping panic"
)

// Mod* 为模块名称标识，作为日志字段标记当前操作的来源模块。
const (
	// ModConfig 为配置模块标识。
	ModConfig = "config"
	// ModApp 为应用模块标识。
	ModApp = "app"
	// ModProc 为进程模块标识。
	ModProc = "proc"
	// ModGrpcServer 为 gRPC 服务端模块标识。
	ModGrpcServer = "server.grpc"
	// ModRegistryETCD 为 Etcd 注册中心模块标识。
	ModRegistryETCD = "registry.etcd"
	// ModClientETCD 为 Etcd 客户端模块标识。
	ModClientETCD = "client.etcd"
	// ModClientGrpc 为 gRPC 客户端模块标识。
	ModClientGrpc = "client.grpc"
	// ModClientMySQL 为 MySQL 客户端模块标识。
	ModClientMySQL = "client.mysql"
	// ModXcronETCD 为基于 Etcd 的定时任务模块标识。
	ModXcronETCD = "xcron.etcd"
)
