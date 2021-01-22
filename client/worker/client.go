package worker

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex

	net, addr string
	conn      net.Conn
	rw        *bufio.ReadWriter

	Req      *Request
	ResQueue chan *Response

	Timeout time.Duration

	ErrHandler   ErrHandler
	RespHandlers *RespHandlerMap
}

func NewClient(network, addr string) (client *Client, err error) {
	client = &Client{
		net:          network,
		addr:         addr,
		Req:          nil,
		ResQueue:     make(chan *Response, QUEUE_SIZE),
		Timeout:      DEFAULT_TIME_OUT,
		RespHandlers: NewResHandlerMap(),
	}
	client.conn, err = net.Dial(client.net, client.addr)
	if err != nil {
		return nil, err
	}

	client.rw = bufio.NewReadWriter(bufio.NewReader(client.conn), bufio.NewWriter(client.conn))

	go client.ClientRead()

	return client, nil
}

func (c *Client) Write() (err error) {
	var n int
	buf := c.Req.EncodePack()
	for i := 0; i < len(buf); i += n {
		n, err = c.rw.Write(buf)
		if err != nil {
			return err
		}
	}

	return c.rw.Flush()
}

func (c *Client) Read(length int) (data []byte, err error) {
	n := 0
	buf := GetBuffer(length)
	for i := length; i > 0 || len(data) < MIN_DATA_SIZE; i -= n {
		if n, err = c.rw.Read(buf); err != nil {
			return
		}
		data = append(data, buf[0:n]...)
		if n < MIN_DATA_SIZE {
			break
		}
	}

	return
}

func (c *Client) ClientRead() {
	var data, leftdata []byte
	var err error
	var res *Response
	var resLen int
Loop:
	for c.conn != nil {
		if data, err = c.Read(MIN_DATA_SIZE); err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if opErr.Timeout() {
					log.Println(err)
				}
				if opErr.Temporary() {
					continue
				}
				break
			}

			//服务端断开
			if err == io.EOF {
				c.ErrHandler(err)
			}

			//断开重连
			log.Println("client read error here:" + err.Error())
			c.Close()
			c.conn, err = net.Dial(c.net, c.addr)
			if err != nil {
				break
			}
			c.rw = bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
			continue
		}

		if len(leftdata) > 0 {
			data = append(leftdata, data...)
			leftdata = nil
		}

		for {
			l := len(data)
			if l < MIN_DATA_SIZE {
				leftdata = data
				continue Loop
			}

			if len(leftdata) == 0 {
				connType := GetConnType(data)
				if connType != CONN_TYPE_SERVER {
					log.Println("read conn type error")
					break
				}
			}

			if res, resLen, err = DecodePack(data); err != nil {
				leftdata = data[:resLen]
				continue Loop
			} else {
				c.ResQueue <- res
			}

			data = data[l:]
			if len(data) > 0 {
				continue
			}
			break
		}
	}
}

func (c *Client) HandlerResp(resp *Response) {
	if resp == nil {
		return
	}
	if len(resp.Handle) == 0 || resp.HandleLen == 0 {
		return
	}

	key := resp.Handle
	if handler, exist := c.RespHandlers.GetResHandlerMap(key); exist {
		handler(resp)
		c.RespHandlers.DelResHandlerMap(key)
		return
	}
}

func (c *Client) ProcessResp() {
	var timer = time.After(c.Timeout)
	select {
	case res := <-c.ResQueue:
		switch res.DataType {
		case PDT_ERROR:
			c.ErrHandler(res.GetResError())
			return
		case PDT_CANT_DO:
			c.ErrHandler(res.GetResError())
			return
		case PDT_S_RETURN_DATA:
			c.HandlerResp(res)
			return
		}
	case <-timer:
		log.Println("time out")
		c.ErrHandler(RESTIMEOUT)
		//c.Close()
		return
	}
}

func (c *Client) Do(funcName string, params []byte, callback RespHandler) (err error) {
	c.Lock()
	defer c.Unlock()

	if c.conn == nil {
		return fmt.Errorf("conn fail")
	}

	c.RespHandlers.PutResHandlerMap(funcName, callback)

	c.Req = NewReq()
	c.Req.ContentPack(PDT_C_DO_JOB, funcName, params)
	if err = c.Write(); err != nil {
		return err
	}

	c.ProcessResp()

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	close(c.ResQueue)
	os.Exit(1)
}
