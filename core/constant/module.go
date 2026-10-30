package constant

// Module 标识单例缓存所属的客户端、注册中心或存储模块，用于按模块隔离缓存键空间。
type Module int

// 模块标识常量。
const (
	// ModuleInvalid 为未指定模块时的非法取值。
	ModuleInvalid = Module(iota)

	// ModuleClientGrpc 为 gRPC 客户端模块。
	ModuleClientGrpc
	// ModuleClientResty 为 Resty HTTP 客户端模块。
	ModuleClientResty
	// ModuleClientRedis 为 Redis 客户端模块。
	ModuleClientRedis
	// ModuleClientRocketMQ 为 RocketMQ 客户端模块。
	ModuleClientRocketMQ
	// ModuleClientEtcd 为 Etcd 客户端模块。
	ModuleClientEtcd

	// ModuleRegistryEtcd 为 Etcd 注册中心模块。
	ModuleRegistryEtcd

	// ModuleStoreMongoDB 为 MongoDB 存储模块。
	ModuleStoreMongoDB
	// ModuleStoreGorm 为 GORM 关系型数据库存储模块。
	ModuleStoreGorm
	// ModuleStoreTableStore 为 TableStore 存储模块。
	ModuleStoreTableStore
)
