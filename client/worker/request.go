package worker

import (
	"encoding/binary"
)

type Request struct {
	DataType uint32
	Data     []byte
	DataLen  uint32

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

//打包内容
func (req *Request) ContentPack(dataType uint32, handle string, params []byte) (content []byte, contentLen uint32) {
	req.DataType = dataType
	req.Handle = handle
	req.HandleLen = uint32(len(handle))
	req.Params = params
	req.ParamsLen = uint32(len(params))
	req.DataLen = uint32(UINT32_SIZE + req.HandleLen + UINT32_SIZE + req.ParamsLen)
	contentLen = req.DataLen

	content = make([]byte, contentLen)
	binary.BigEndian.PutUint32(content[:UINT32_SIZE], req.HandleLen)
	start := UINT32_SIZE
	end := UINT32_SIZE + int(req.HandleLen)
	copy(content[start:end], []byte(req.Handle))
	start = end
	end = start + UINT32_SIZE
	binary.BigEndian.PutUint32(content[start:end], req.ParamsLen)
	start = end
	end = start + int(req.ParamsLen)
	copy(content[start:end], req.Params)
	req.Data = content

	return
}

//打包
func (req *Request) EncodePack() (data []byte) {
	len := MIN_DATA_SIZE + req.DataLen //add 12 bytes head
	data = GetBuffer(int(len))

	binary.BigEndian.PutUint32(data[:4], CONN_TYPE_CLIENT)
	binary.BigEndian.PutUint32(data[4:8], req.DataType)
	binary.BigEndian.PutUint32(data[8:MIN_DATA_SIZE], req.DataLen)
	copy(data[MIN_DATA_SIZE:], req.Data)

	return
}
