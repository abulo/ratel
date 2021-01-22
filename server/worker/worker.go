package worker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type Worker struct {
	sync.Mutex

	Agents   []*Agent
	Funcs    map[string]*Function
	FuncsNum int
	Resps    chan *Response

	ready   bool
	running bool
}

func NewWorker() *Worker {
	return &Worker{
		Agents:   make([]*Agent, 0),
		Funcs:    make(map[string]*Function),
		FuncsNum: 0,
		Resps:    make(chan *Response, QUEUE_SIZE),
		ready:    false,
		running:  false,
	}
}

func (w *Worker) AddServer(net, addr string) (err error) {
	var agent *Agent
	agent = NewAgent(net, addr, w)

	if agent == nil {
		return fmt.Errorf("agent nil")
	}
	w.Agents = append(w.Agents, agent)

	return nil
}

func (w *Worker) AddFunction(funcName string, jobFunc JobFunc) (err error) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.Funcs[funcName]; ok {
		return fmt.Errorf("function already exist")
	}

	w.Funcs[funcName] = NewFunction(jobFunc, funcName)
	w.FuncsNum++
	if w.running {
		go w.FuncBroadcast(funcName, PDT_W_ADD_FUNC)
	}

	return nil
}

func (w *Worker) DelFunction(funcName string) (err error) {
	w.Lock()
	defer w.Unlock()

	if _, ok := w.Funcs[funcName]; !ok {
		return fmt.Errorf("function not exist")
	}

	delete(w.Funcs, funcName)
	w.FuncsNum--
	if w.running {
		go w.FuncBroadcast(funcName, PDT_W_DEL_FUNC)
	}

	return nil
}

func (w *Worker) GetFunction(funcName string) (function *Function, err error) {
	w.Lock()
	defer w.Unlock()

	if w.FuncsNum == 0 {
		return nil, fmt.Errorf("worker have no funcs")
	}

	f, ok := w.Funcs[funcName]
	if !ok {
		return nil, fmt.Errorf("not found")
	}

	if f.FuncName != funcName {
		return nil, fmt.Errorf("not found")
	}

	function = f

	return function, nil
}

func (w *Worker) DoFunction(resp *Response) (err error) {
	if resp.DataType == PDT_S_GET_DATA {
		funcName := resp.Handle
		if function, err := w.GetFunction(funcName); err != nil {
			return err
		} else if function != nil {
			if function.FuncName != funcName {
				return fmt.Errorf("funcname error")
			}
			if resp.ParamsLen == 0 {
				return fmt.Errorf("params error")
			}

			var ret []byte
			if ret, err = function.Func(resp); err == nil {
				resp.Agent.Req.HandleLen = resp.HandleLen
				resp.Agent.Req.Handle = resp.Handle
				resp.Agent.Req.ParamsLen = resp.ParamsLen
				resp.Agent.Req.Params = resp.Params

				resp.Agent.Req.RetPack(ret)
				resp.Agent.Write()
			}
		}

	}

	return nil
}

func (w *Worker) FuncBroadcast(funcName string, flag int) {
	for _, a := range w.Agents {
		switch flag {
		case PDT_W_ADD_FUNC:
			a.Req.AddFunctionPack(funcName)
		case PDT_W_DEL_FUNC:
			a.Req.DelFunctionPack(funcName)
		default:
			a.Req.AddFunctionPack(funcName)
		}
		a.Write()
	}
}

func (w *Worker) WorkerReady() (err error) {
	if len(w.Agents) == 0 {
		return errors.New("None active agents")
	}
	if w.FuncsNum == 0 || len(w.Funcs) == 0 {
		return errors.New("None functions")
	}

	for _, a := range w.Agents {
		if err = a.Connect(); err != nil {
			return err
		}
	}

	for fn := range w.Funcs {
		w.FuncBroadcast(fn, PDT_W_ADD_FUNC)
	}
	w.ready = true

	return nil
}

func (w *Worker) WorkerDo() {
	if !w.ready {
		err := w.WorkerReady()
		if err != nil {
			panic(err)
		}
	}

	w.Lock()
	w.running = true
	w.Unlock()

	for _, a := range w.Agents {
		a.Grab()
	}

	n := 1
	for resp := range w.Resps {
		switch resp.DataType {
		case PDT_TOSLEEP:
			time.Sleep(time.Duration(2) * time.Second)
			go resp.Agent.Wakeup()

			//fallthrough
		case PDT_S_GET_DATA:
			fmt.Println("######server get data")
			fmt.Println("######resps num n:", n)
			go func() {
				if err := w.DoFunction(resp); err != nil {
					log.Println(err)
				}
			}()

			//fallthrough
		case PDT_NO_JOB:
			fmt.Println("######no job Grab")
			fmt.Println("######resps num n:", n)
			go resp.Agent.Grab()

		case PDT_WAKEUPED:
		default:
			fmt.Println("######default Grab")
			fmt.Println("######resps num n:", n)
			go resp.Agent.Grab()
		}

		n++
	}

}

func (w *Worker) WorkerClose() error {
	if w.running == true {
		for _, a := range w.Agents {
			a.Close()
		}

		w.running = false
		close(w.Resps)
	}
	return nil
}
