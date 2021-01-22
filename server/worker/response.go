package worker

import (
	"encoding/binary"
	"fmt"
)

type Response struct {
	DataType   uint32
	Data       []byte
	DataLen    uint32
	Handle     string
	HandleLen  uint32
	ParamsType uint32
	ParamsNum  uint32
	ParamsLen  uint32
	Params     []byte
	StrParams  []string
	Ret        []byte
	RetLen     uint32
	Agent      *Agent
}

func NewRes() (res *Response) {
	res = &Response{
		Data:       make([]byte, 0),
		DataLen:    0,
		Handle:     ``,
		HandleLen:  0,
		ParamsType: PARAMS_TYPE_ONE, //4-one param, 5-multiple params, default 4
		ParamsNum:  0,               //参数个数，一般单个参数只有一个参数，多个参数有相应数量的参数
		ParamsLen:  0,
		Params:     make([]byte, 0),
		StrParams:  make([]string, 0),
		Ret:        make([]byte, 0),
		RetLen:     0,
	}
	return
}

//解包
func DecodePack(data []byte) (resp *Response, resLen int, err error) {
	resLen = len(data)
	if resLen < MIN_DATA_SIZE {
		err = fmt.Errorf("Invalid data: %v", data)
		return
	}
	cl := int(binary.BigEndian.Uint32(data[8:MIN_DATA_SIZE]))
	if resLen < MIN_DATA_SIZE+cl {
		err = fmt.Errorf("Invalid data: %v", data)
		return
	}
	content := data[MIN_DATA_SIZE : MIN_DATA_SIZE+cl]
	if len(content) != cl {
		err = fmt.Errorf("Invalid data: %v", data)
		return
	}

	resp = NewRes()
	resp.DataType = binary.BigEndian.Uint32(data[4:8])
	resp.DataLen = uint32(cl)
	resp.Data = content

	if resp.DataType == PDT_S_GET_DATA {
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
		//resp.ParseParams(data[start:end])

		//新的解包协议
		start := MIN_DATA_SIZE
		end := MIN_DATA_SIZE + UINT32_SIZE
		resp.HandleLen = binary.BigEndian.Uint32(data[start:end])
		start = end
		end = start + UINT32_SIZE
		resp.ParamsLen = binary.BigEndian.Uint32(data[start:end])
		start = end
		end = start + int(resp.HandleLen)
		resp.Handle = string(data[start:end])
		start = end
		end = start + int(resp.ParamsLen)
		resp.ParseParams(data[start:end])
	}

	return
}

func (resp *Response) GetResponse() *Response {
	return resp
}

func (resp *Response) ParseParams(params []byte) {
	resp.Params = params
	strArrParams := GetStrParamsArr(params)
	if strArrParams != nil {
		resp.StrParams = strArrParams
		resp.ParamsNum = uint32(len(strArrParams))
		if resp.ParamsNum > 1 {
			resp.ParamsType = PARAMS_TYPE_MUL
		}
	}

	return
}

func (resp *Response) GetParams() []byte {
	if resp.ParamsLen == 0 {
		return nil
	}

	return resp.Params
}

func (resp *Response) GetStrParams() []string {
	if resp.ParamsLen == 0 {
		return nil
	}

	return resp.StrParams
}
