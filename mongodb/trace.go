package mongodb

var trace bool

//SetTrace 设置跟踪
func (configs *Configs) SetTrace(t bool) {
	trace = t
}
