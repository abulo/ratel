package tableflip

import (
	"context"
	"log"
	"net"
	"net/http"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const (
	PRE_SIGNAL = iota
	POST_SIGNAL

	STATE_INIT
	STATE_RUNNING
	STATE_SHUTTING_DOWN
	STATE_TERMINATE
)

var (
	runningServerReg sync.RWMutex

	DefaultReadTimeOut    time.Duration
	DefaultWriteTimeOut   time.Duration
	DefaultMaxHeaderBytes int
	DefaultHammerTime     time.Duration
)

func init() {
	DefaultHammerTime = 30 * time.Second
}

type httpServer struct {
	http.Server

	tableflipListener net.Listener

	wg    sync.WaitGroup
	state uint8
	lock  *sync.RWMutex
}

func NewHttpServer(addr string, handler http.Handler) *httpServer {
	runningServerReg.Lock()
	defer runningServerReg.Unlock()

	srv := &httpServer{
		wg:    sync.WaitGroup{},
		state: STATE_INIT,
		lock:  &sync.RWMutex{},
	}

	srv.Server.ReadTimeout = DefaultReadTimeOut
	srv.Server.WriteTimeout = DefaultWriteTimeOut
	srv.Server.MaxHeaderBytes = DefaultMaxHeaderBytes
	srv.Server.Handler = handler

	return srv
}

func (srv *httpServer) Serve() error {
	srv.setState(STATE_RUNNING)
	err := srv.Server.Serve(srv.tableflipListener)
	srv.wg.Wait()
	srv.setState(STATE_TERMINATE)
	return err
}

func (srv *httpServer) getState() uint8 {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	return srv.state
}

func (srv *httpServer) setState(st uint8) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	srv.state = st
}

func (srv *httpServer) Shutdown() {
	if srv.getState() != STATE_RUNNING {
		return
	}

	srv.setState(STATE_SHUTTING_DOWN)

	go srv.hammerTime(DefaultHammerTime)
	//srv.SetKeepAlivesEnabled(false)
	//err := srv.tableflipListener.Close()
	ctx, _ := context.WithTimeout(context.TODO(), DefaultHammerTime)
	err := srv.Server.Shutdown(ctx)

	if err != nil {
		log.Println(syscall.Getpid(), "http server shutdown error:", err)
	} else {
		log.Println(syscall.Getpid(), srv.tableflipListener.Addr(), "http server shutdown.")
	}
}
func (srv *httpServer) hammerTime(d time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("WaitGroup at 0", r)
		}
	}()
	if srv.getState() != STATE_SHUTTING_DOWN {
		return
	}
	time.Sleep(d)
	log.Println("[STOP - Hammer Time] Forcefully shutting down parent")
	//err := srv.Server.Close()
	//if err != nil {
	//	log.Println("http server close err. %s", err)
	//}
	for {
		if srv.getState() == STATE_TERMINATE {
			break
		}
		srv.wg.Done()
		runtime.Gosched()
	}
}
