package worker

func GetBuffer(n int) (buf []byte) {
	buf = make([]byte, n)
	return
}

func GetRetStruct() *RetStruct {
	return &RetStruct{
		Code: 0,
		Msg:  "",
		Data: make([]byte, 0),
	}
}
