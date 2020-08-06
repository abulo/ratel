package redis

// AutoConfigRedisClient 合并配置文件和环境
func AutoConfigRedisClient(rwType RWType) (*Client, error) {
	opts, err := customizedOptionsFromFullVariable(rwType)
	if opts != nil {
		return NewClient(*opts), err
	}
	return nil, err
}

// AutoConfigRedisClientFromVolume 使用配置文件中的参数创建redisclient
func AutoConfigRedisClientFromVolume(rwType RWType) (*Client, error) {
	opts, err := customizedOptionsFromVolume(rwType)
	if opts != nil {
		return NewClient(*opts), err
	}
	return nil, err
}

// AutoConfigRedisClientFromEnv 使用纯环境变量参数创建redisclient
func AutoConfigRedisClientFromEnv(rwType RWType) (*Client, error) {
	opts, err := customizedOptionsFromEnv(rwType)
	if opts != nil {
		return NewClient(*opts), err
	}
	return nil, err
}
