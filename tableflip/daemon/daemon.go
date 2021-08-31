package daemon

import (
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
)

var errNotSupported = errors.New("daemon: Non-POSIX OS is not supported")

// Mark of daemon process - system environment variable _GO_DAEMON=1
const (
	MARK_NAME  = "GOLANG_DAEMON"
	MARK_VALUE = 5
)

// Default file permissions for log and pid files.
const FILE_PERM = os.FileMode(0644)
const DEF_UMASK = int(022)

type Option func(*Options) error

type Options struct {
	// If PidFileName is non-empty, parent process will try to create and lock
	// pid file with given name. Child process writes process id to file.
	PidFileName string
	// Permissions for new pid file.
	PidFilePerm os.FileMode

	// If LogFileName is non-empty, parent process will create file with given name
	// and will link to fd 2 (stderr) for child process.
	LogFileName string
	// Permissions for new log file.
	LogFilePerm os.FileMode

	// If WorkDir is non-empty, the child changes into the directory before
	// creating the process.
	WorkDir string
	// If Chroot is non-empty, the child changes root directory
	Chroot string

	// Credential holds user and group identities to be assumed by a daemon-process.
	//Credential *syscall.Credential
	// If Umask is non-zero, the daemon-process call Umask() func with given value.
	Umask int
}

// A Context describes daemon context.
type Daemon struct {
	Opts *Options

	// If Env is non-nil, it gives the environment variables for the
	// daemon-process in the form returned by os.Environ.
	// If it is nil, the result of os.Environ will be used.
	Env []string
	// If Args is non-nil, it gives the command-line args for the
	// daemon-process. If it is nil, the result of os.Args will be used
	// (without program name).
	Args []string

	// Struct contains only serializable public fields (!!!)
	abspath  string
	pidFile  *LockFile
	logFile  *os.File
	errFile  *os.File
	nullFile *os.File

	// golang bool default value is false
	DontRemovePidFile bool
}

func PidFile(name string, perm os.FileMode) Option {
	return func(opt *Options) error {
		opt.PidFileName = name
		opt.PidFilePerm = perm
		return nil
	}
}

func LogFile(name string, perm os.FileMode) Option {
	return func(opt *Options) error {
		opt.LogFileName = name
		opt.LogFilePerm = perm
		return nil
	}
}
func LogFileName(name string) Option {
	return func(opt *Options) error {
		opt.LogFileName = name
		return nil
	}
}
func LogFilePerm(perm os.FileMode) Option {
	return func(opt *Options) error {
		opt.LogFilePerm = perm
		return nil
	}
}

func WorkDir(dir string) Option {
	return func(opt *Options) error {
		opt.WorkDir = dir
		return nil
	}
}
func Umask(umask int) Option {
	return func(opt *Options) error {
		opt.Umask = umask
		return nil
	}
}

// func SignalConfig(cmd string, sig os.Signal, handler SignalHandlerFunc) Option {
// 	return func(opt *Options) error {
// 		if !flag.Parsed() {
// 			flag.Parse()
// 		}
// 		AddCommand(StringFlag(sigflag, cmd), sig, handler)
// 		return nil
// 	}
// }

func DefaultOption(opts ...Option) (*Options, error) {
	abspath, err := osext.Executable()
	if err != nil {
		return nil, err
	}
	op := &Options{
		PidFileName: abspath + ".pid",
		PidFilePerm: FILE_PERM,

		LogFileName: abspath + ".log",
		LogFilePerm: FILE_PERM,

		WorkDir: filepath.Dir(abspath),
		//Umask:   022,

		//flagSignal: map[string]os.Signal{
		//	"restart": syscall.SIGHUP,
		//	"quit":    syscall.SIGQUIT,
		//},

		//signalHandlers: map[os.Signal]SignalHandlerFunc{
		//	syscall.SIGHUP:  sigtermDefaultHandler,
		//	syscall.SIGQUIT: sigtermDefaultHandler,
		//},
	}

	for _, fun := range opts {
		if err := fun(op); err != nil {
			return nil, err
		}
	}
	return op, nil
}

// WasReborn returns true in child process (daemon) and false in parent process.
func WasReborn() bool {
	return os.Getenv(MARK_NAME) != ""
}

// Reborn runs second copy of current process in the given context.
// function executes separate parts of code in child process and parent process
// and provides demonization of child process. It look similar as the
// fork-daemonization, but goroutine-safe.
// In success returns *os.Process in parent process and nil in child process.
// Otherwise returns error.
func (d *Daemon) Reborn() (child *os.Process, err error) {
	return d.reborn()
}

// Search search daemons process by given in context pid file name.
// If success returns pointer on daemons os.Process structure,
// else returns error. Returns nil if filename is empty.
func (d *Daemon) Search() (daemon *os.Process, err error) {
	return d.search()
}

// Release provides correct pid-file release in daemon.
func (d *Daemon) Release() (err error) {
	return d.release()
}
