package tableflip

import (
	"net"
	"os"
	"sync"
	"syscall"
	"time"
)

func newTableflipListener(l net.Listener, wg *sync.WaitGroup) *tableflipListener {
	return &tableflipListener{
		Listener: l,
		wg:       wg,
	}
}

type tableflipListener struct {
	net.Listener
	stopped bool
	wg      *sync.WaitGroup
}

func (this *tableflipListener) Accept() (net.Conn, error) {
	conn, err := this.Listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return nil, err
	}
	conn.SetKeepAlive(true)                  // see http.tcpKeepAliveListener
	conn.SetKeepAlivePeriod(3 * time.Minute) // see http.tcpKeepAliveListener

	nxconn := tableflipConn{
		Conn: conn,
		wg:   this.wg,
	}
	this.wg.Add(1)
	return nxconn, nil
}

func (this *tableflipListener) Close() error {
	if this.stopped {
		return syscall.EINVAL
	}
	this.stopped = true
	return this.Listener.Close()
}

func (this *tableflipListener) File() (*os.File, error) {
	// returns a dup(2) - FD_CLOEXEC flag *not* set
	//tl := this.Listener.(*net.TCPListener)
	//fl, _ := tl.File()
	tl, err := this.Listener.(filer).File()
	return tl, err
}

type filer interface {
	File() (*os.File, error)
}

type tableflipTlsListenner struct {
	tableflipListener
}

type tableflipConn struct {
	net.Conn
	wg *sync.WaitGroup
}

func (nxc tableflipConn) Close() error {
	err := nxc.Conn.Close()
	if err == nil {
		nxc.wg.Done()
	}
	return err
}
