package mongodb

var trace bool

//SetTrace 设置跟踪
func (config *Config) SetTrace(t bool) {
	trace = t
}
