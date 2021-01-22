package worker

import "encoding/binary"

type Request struct {
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

func NewReq() (req *Request) {
	req = &Request{
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

//打包内容-添加方法
func (req *Request) AddFunctionPack(funcName string) (content []byte, err error) {
	req.DataType = PDT_W_ADD_FUNC
	req.DataLen = uint32(len(funcName))
	req.Data = []byte(funcName)
	content = req.Data

	return
}

//打包内容-删除方法
func (req *Request) DelFunctionPack(funcName string) (content []byte, err error) {
	req.DataType = PDT_W_DEL_FUNC
	req.DataLen = uint32(len(funcName))
	req.Data = []byte(funcName)
	content = req.Data

	return
}

//打包内容-抓取任务
func (req *Request) GrabDataPack() (content []byte, err error) {
	req.DataType = PDT_W_GRAB_JOB
	req.DataLen = 0
	req.Data = []byte(``)
	content = req.Data

	return
}

//打包内容-唤醒
func (req *Request) WakeupPack() {
	req.DataType = PDT_WAKEUP
	req.DataLen = 0
	req.Data = []byte(``)

	return
}

//打包内容-返回结果
func (req *Request) RetPack(ret []byte) (content []byte, err error) {
	req.Ret = ret
	req.RetLen = uint32(len(ret))

	req.DataType = PDT_W_RETURN_DATA
	req.DataLen = UINT32_SIZE + req.HandleLen + UINT32_SIZE + req.ParamsLen + UINT32_SIZE + req.RetLen

	length := int(req.DataLen)
	content = GetBuffer(length)
	binary.BigEndian.PutUint32(content[:UINT32_SIZE], req.HandleLen)
	start := UINT32_SIZE
	end := int(UINT32_SIZE + req.HandleLen)
	copy(content[start:end], []byte(req.Handle))
	start = end
	end = start + UINT32_SIZE
	binary.BigEndian.PutUint32(content[start:end], uint32(req.ParamsLen))
	start = end
	end = start + int(req.ParamsLen)
	copy(content[start:end], req.Params)
	start = end
	end = start + UINT32_SIZE
	binary.BigEndian.PutUint32(content[start:end], req.RetLen)
	start = end
	end = start + int(req.RetLen)
	copy(content[start:end], req.Ret)
	req.Data = content

	return
}

//打包
func (req *Request) EncodePack() (data []byte) {
	len := MIN_DATA_SIZE + req.DataLen //add 12 bytes head
	data = GetBuffer(int(len))

	binary.BigEndian.PutUint32(data[:4], CONN_TYPE_WORKER)
	binary.BigEndian.PutUint32(data[4:8], req.DataType)
	binary.BigEndian.PutUint32(data[8:MIN_DATA_SIZE], req.DataLen)
	copy(data[MIN_DATA_SIZE:], req.Data)

	return
}
