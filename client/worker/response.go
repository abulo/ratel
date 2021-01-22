package worker

import (
	"encoding/binary"
	"fmt"
	"sync"
)

type Response struct {
	DataType   uint32
	Data       []byte
	DataLen    uint32
	Handle     string
	HandleLen  uint32
	ParamsType uint32
	ParamsLen  uint32
	Params     []byte
	Ret        []byte
	RetLen     uint32
}

type ErrHandler func(err error)

type RespHandler func(resp *Response)

type RespHandlerMap struct {
	sync.Mutex
	holder map[string]RespHandler
}

func NewRes() (res *Response) {
	res = &Response{
		Data:      make([]byte, 0),
		DataLen:   0,
		Handle:    ``,
		HandleLen: 0,
		ParamsLen: 0,
		Params:    make([]byte, 0),
		Ret:       make([]byte, 0),
		RetLen:    0,
	}
	return
}

func (resp *Response) GetResError() (err error) {
	if resp.DataType == PDT_ERROR {
		return fmt.Errorf("request error")
	} else if resp.DataType == PDT_CANT_DO {
		return fmt.Errorf("have no job do")
	}

	return nil
}

func (resp *Response) GetResResult() (data []byte, err error) {
	if resp.DataType == PDT_S_RETURN_DATA {
		return resp.Ret, nil
	}

	return []byte(``), fmt.Errorf("data nil")
}

func GetConnType(data []byte) (connType uint32) {
	if len(data) == 0 {
		return 0
	}

	if len(data) < 4 {
		return 0
	}

	connType = uint32(binary.BigEndian.Uint32(data[:4]))

	return
}

//解包
func DecodePack(data []byte) (resp *Response, resLen int, err error) {
	resLen = len(data)
	if resLen < MIN_DATA_SIZE {
		err = fmt.Errorf("Invalid data1: %v", data)
		return
	}
	cl := int(binary.BigEndian.Uint32(data[8:MIN_DATA_SIZE]))
	if resLen < MIN_DATA_SIZE+cl {
		err = fmt.Errorf("Invalid data2: %v", data)
		return
	}
	content := data[MIN_DATA_SIZE : MIN_DATA_SIZE+cl]
	if len(content) != cl {
		err = fmt.Errorf("Invalid data3: %v", data)
		return
	}

	resp = NewRes()
	resp.DataType = binary.BigEndian.Uint32(data[4:8])
	resp.DataLen = uint32(cl)
	resp.Data = content

	if resp.DataType == PDT_S_RETURN_DATA {
		//旧的解包协议
		//start := MIN_DATA_SIZE
		//end   := MIN_DATA_SIZE + UINT32_SIZE
		//resp.HandleLen = binary.BigEndian.Uint32(data[start:end])
		//start = end
		//end   = start + int(resp.HandleLen)
		//resp.Handle = string(data[start:end])
		//start = end
		//end   = start + UINT32_SIZE
		//resp.ParamsLen = binary.BigEndian.Uint32(data[start:end])
		//start = end
		//end   = start + int(resp.ParamsLen)
		//resp.Params = data[start:end]
		//start = end
		//end   = start + UINT32_SIZE
		//resp.RetLen = binary.BigEndian.Uint32(data[start:end])
		//start = end
		//end   = start + int(resp.RetLen)
		//resp.Ret = data[start:end]

		//新的解包协议
		start := MIN_DATA_SIZE
		end := MIN_DATA_SIZE + UINT32_SIZE
		resp.HandleLen = binary.BigEndian.Uint32(data[start:end])
		start = end
		end = start + UINT32_SIZE
		resp.ParamsLen = binary.BigEndian.Uint32(data[start:end])
		start = end
		end = start + UINT32_SIZE
		resp.RetLen = binary.BigEndian.Uint32(data[start:end])
		start = end
		end = start + int(resp.HandleLen)
		resp.Handle = string(data[start:end])
		start = end
		end = start + int(resp.ParamsLen)
		resp.Params = data[start:end]
		start = end
		end = start + int(resp.RetLen)
		resp.Ret = data[start:end]
	}

	return
}

func NewResHandlerMap() *RespHandlerMap {
	return &RespHandlerMap{
		holder: make(map[string]RespHandler, QUEUE_SIZE),
	}
}

func (rhm *RespHandlerMap) PutResHandlerMap(key string, handler RespHandler) {
	rhm.Lock()
	rhm.holder[key] = handler
	rhm.Unlock()
}

func (rhm *RespHandlerMap) GetResHandlerMap(key string) (handler RespHandler, exist bool) {
	rhm.Lock()
	handler, exist = rhm.holder[key]
	rhm.Unlock()
	return handler, exist
}

func (rhm *RespHandlerMap) DelResHandlerMap(key string) {
	rhm.Lock()
	delete(rhm.holder, key)
	rhm.Unlock()
}
