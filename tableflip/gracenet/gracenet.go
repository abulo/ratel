package gracenet

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/kardianos/osext"
)

const (
	// Used to indicate a graceful restart in the new process.
	envCountKey       = "LISTEN_FDS"
	envCountKeyPrefix = envCountKey + "="
)

// In order to keep the working directory the same as when we started we record
// it at startup.
var originalWD, _ = os.Getwd()

type RestartAttr interface {
	Files() []*os.File
	Envs() []string
	FdStart() int
}

type RestartObserver interface {
	BeforeRestart()
	AfterRestart(pro *os.Process, err error)
}

// Net provides the family of Listen functions and maintains the associated
// state. Typically you will have only once instance of Net per application.
type Net struct {
	inherited   []net.Listener
	active      []net.Listener
	mutex       sync.Mutex
	inheritOnce sync.Once

	// used in tests to override the default behavior of starting from fd 3.
	fdStart   int
	ProcAttr  RestartAttr
	observers []RestartObserver
}

func (n *Net) AddObserver(o RestartObserver) {
	n.observers = append(n.observers, o)
}

func (n *Net) RemoveObserver(o RestartObserver) bool {
	for i, v := range n.observers {
		if v == o {
			n.observers = append(n.observers[:i], n.observers[i+1:]...)
			return true
		}
	}
	return false
}

func (n *Net) NotifyBeforeRestart() {
	for _, o := range n.observers {
		o.BeforeRestart()
	}
}
func (n *Net) NotifyAfterRestart(pro *os.Process, err error) {
	for _, o := range n.observers {
		o.AfterRestart(pro, err)
	}
}

func (n *Net) Inherited() bool {
	return len(n.inherited) > 0
}

func (n *Net) inherit() error {
	var retErr error
	n.inheritOnce.Do(func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()
		countStr := os.Getenv(envCountKey)
		if countStr == "" {
			return
		}
		count, err := strconv.Atoi(countStr)
		if err != nil {
			retErr = fmt.Errorf("found invalid count value: %s=%s", envCountKey, countStr)
			return
		}

		// In tests this may be overridden.
		fdStart := n.fdStart
		if fdStart == 0 {
			// In normal operations if we are inheriting, the listeners will begin at
			// fd 3.
			//fdStart = 3
			if n.ProcAttr != nil {
				fdStart = n.ProcAttr.FdStart()
			} else {
				fdStart = 3
			}
		}

		for i := fdStart; i < fdStart+count; i++ {
			file := os.NewFile(uintptr(i), "listener")
			l, err := net.FileListener(file)
			if err != nil {
				file.Close()
				retErr = fmt.Errorf("error inheriting socket fd %d: %s", i, err)
				return
			}
			if err := file.Close(); err != nil {
				retErr = fmt.Errorf("error closing inherited socket fd %d: %s", i, err)
				return
			}
			n.inherited = append(n.inherited, l)
		}
	})
	return retErr
}

// Listen announces on the local network address laddr. The network net must be
// a stream-oriented network: "tcp", "tcp4", "tcp6", "unix" or "unixpacket". It
// returns an inherited net.Listener for the matching network and address, or
// creates a new one using net.Listen.
func (n *Net) Listen(nett, laddr string) (net.Listener, error) {
	switch nett {
	default:
		return nil, net.UnknownNetworkError(nett)
	case "tcp", "tcp4", "tcp6":
		addr, err := net.ResolveTCPAddr(nett, laddr)
		if err != nil {
			return nil, err
		}
		fmt.Println("net listen 1\n")
		return n.ListenTCP(nett, addr)
	case "unix", "unixpacket", "invalid_unix_net_for_test":
		addr, err := net.ResolveUnixAddr(nett, laddr)
		if err != nil {
			return nil, err
		}
		return n.ListenUnix(nett, addr)
	}
}

// ListenTCP announces on the local network address laddr. The network net must
// be: "tcp", "tcp4" or "tcp6". It returns an inherited net.Listener for the
// matching network and address, or creates a new one using net.ListenTCP.
func (n *Net) ListenTCP(nett string, laddr *net.TCPAddr) (*net.TCPListener, error) {
	if err := n.inherit(); err != nil {
		return nil, err
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	// look for an inherited listener
	for i, l := range n.inherited {
		if l == nil { // we nil used inherited listeners
			continue
		}
		if isSameAddr(l.Addr(), laddr) {
			n.inherited[i] = nil
			n.active = append(n.active, l)
			return l.(*net.TCPListener), nil
		}
	}

	// make a fresh listener
	l, err := net.ListenTCP(nett, laddr)
	if err != nil {
		return nil, err
	}
	n.active = append(n.active, l)
	return l, nil
}

// ListenUnix announces on the local network address laddr. The network net
// must be a: "unix" or "unixpacket". It returns an inherited net.Listener for
// the matching network and address, or creates a new one using net.ListenUnix.
func (n *Net) ListenUnix(nett string, laddr *net.UnixAddr) (*net.UnixListener, error) {
	if err := n.inherit(); err != nil {
		return nil, err
	}

	n.mutex.Lock()
	defer n.mutex.Unlock()

	// look for an inherited listener
	for i, l := range n.inherited {
		if l == nil { // we nil used inherited listeners
			continue
		}
		if isSameAddr(l.Addr(), laddr) {
			n.inherited[i] = nil
			n.active = append(n.active, l)
			return l.(*net.UnixListener), nil
		}
	}

	// make a fresh listener
	l, err := net.ListenUnix(nett, laddr)
	if err != nil {
		return nil, err
	}
	n.active = append(n.active, l)
	return l, nil
}

// activeListeners returns a snapshot copy of the active listeners.
func (n *Net) activeListeners() ([]net.Listener, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	ls := make([]net.Listener, len(n.active))
	copy(ls, n.active)
	return ls, nil
}

func isSameAddr(a1, a2 net.Addr) bool {
	if a1.Network() != a2.Network() {
		return false
	}
	a1s := a1.String()
	a2s := a2.String()
	if a1s == a2s {
		return true
	}

	// This allows for ipv6 vs ipv4 local addresses to compare as equal. This
	// scenario is common when listening on localhost.
	const ipv6prefix = "[::]"
	a1s = strings.TrimPrefix(a1s, ipv6prefix)
	a2s = strings.TrimPrefix(a2s, ipv6prefix)
	const ipv4prefix = "0.0.0.0"
	a1s = strings.TrimPrefix(a1s, ipv4prefix)
	a2s = strings.TrimPrefix(a2s, ipv4prefix)
	return a1s == a2s
}

// StartProcess starts a new process passing it the active listeners. It
// doesn't fork, but starts a new process using the same environment and
// arguments as when it was originally started. This allows for a newly
// deployed binary to be started. It returns the pid of the newly started
// process when successful.
func (n *Net) StartProcess() (int, error) {
	listeners, err := n.activeListeners()
	if err != nil {
		return 0, err
	}

	// Extract the fds from the listeners.
	files := make([]*os.File, len(listeners))
	for i, l := range listeners {
		files[i], err = l.(filer).File()
		if err != nil {
			return 0, err
		}
		syscall.SetNonblock(int(files[i].Fd()), true)
		defer files[i].Close()
	}

	// Use the original binary location. This works with symlinks such that if
	// the file it points to has been changed we will use the updated symlink.
	//argv0, err := exec.LookPath(os.Args[0])
	argv0, err := osext.Executable()
	if err != nil {
		return 0, err
	}

	// Pass on the environment and replace the old count key with the new one.
	var env []string
	for _, v := range os.Environ() {
		if !strings.HasPrefix(v, envCountKeyPrefix) {
			env = append(env, v)
		}
	}
	env = append(env, fmt.Sprintf("%s%d", envCountKeyPrefix, len(listeners)))

	var allFiles []*os.File
	if n.ProcAttr != nil {
		env = append(env, n.ProcAttr.Envs()...)
		allFiles = append(n.ProcAttr.Files(), files...)
	} else {
		allFiles = append([]*os.File{os.Stdin, os.Stdout, os.Stderr}, files...)
	}
	fmt.Printf("allfiles %#v %d\n", allFiles, len(files))

	fmt.Printf("cwd:%s exe:%s args:%s\n", originalWD, argv0, os.Args)
	n.NotifyBeforeRestart()
	process, err := os.StartProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   originalWD,
		Env:   env,
		Files: allFiles,
		Sys: &syscall.SysProcAttr{
			Setsid: true,
		},
	})
	n.NotifyAfterRestart(process, err)
	if err != nil {
		fmt.Printf("start process err. %s", err)
		return 0, err
	}
	go func(proc *os.Process) {
		state, err := proc.Wait()
		if err != nil {
			fmt.Printf("child %d wait err. %s\n", proc.Pid, err)
			return
		}
		fmt.Printf("child %d wait . %v\n", proc.Pid, state)
	}(process)
	fmt.Printf("start process ok. %d\n", process.Pid)
	return process.Pid, nil
}

type filer interface {
	File() (*os.File, error)
}
