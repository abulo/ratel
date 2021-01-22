package worker

import (
	"github.com/abulo/ratel/logger"
	"github.com/vmihailenco/msgpack"
)

func GetBuffer(n int) (buf []byte) {
	buf = make([]byte, n)
	return
}

func GetStrParamsArr(params []byte) []string {
	var strParamsArr []string

	err := msgpack.Unmarshal(params, &strParamsArr)
	if err != nil {
		logger.Logger.Info("msgpack unmarshal error:", err)
		return nil
	}

	return strParamsArr
}

func GetRetStruct() *RetStruct {
	return &RetStruct{
		Code: 0,
		Msg:  "",
		Data: make([]byte, 0),
	}
}
