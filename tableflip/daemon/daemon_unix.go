//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || plan9 || solaris

package daemon

import (
	"fmt"
	"os"
	"syscall"

	"github.com/kardianos/osext"
)

func (d *Daemon) reborn() (child *os.Process, err error) {
	if !WasReborn() {
		child, err = d.parent()
	} else {
		err = d.child()
	}
	return
}

func (d *Daemon) search() (daemon *os.Process, err error) {
	if len(d.Opts.PidFileName) > 0 {
		var pid int
		if pid, err = ReadPidFile(d.Opts.PidFileName); err != nil {
			return
		}
		daemon, err = os.FindProcess(pid)
	}
	return
}

func (d *Daemon) parent() (child *os.Process, err error) {
	if err = d.prepareEnv(); err != nil {
		return
	}

	defer d.closeFiles()
	if err = d.openFiles(); err != nil {
		return
	}

	attr := &os.ProcAttr{
		Dir:   d.Opts.WorkDir,
		Env:   d.Env,
		Files: d.Files(),
		Sys: &syscall.SysProcAttr{
			//Chroot:     d.Chroot,
			//Credential: d.Opts.Credential,
			Setsid: true,
		},
	}

	if child, err = os.StartProcess(d.abspath, d.Args, attr); err != nil {
		if d.pidFile != nil {
			d.pidFile.Remove()
		}
		return
	}
	//d.RemovePidFile = false

	return
}

func (d *Daemon) openFiles() (err error) {
	if d.Opts.PidFilePerm == 0 {
		d.Opts.PidFilePerm = FILE_PERM
	}
	if d.Opts.LogFilePerm == 0 {
		d.Opts.LogFilePerm = FILE_PERM
	}

	if d.nullFile, err = os.Open(os.DevNull); err != nil {
		return
	}

	if len(d.Opts.PidFileName) > 0 {
		if d.pidFile, err = OpenLockFile(d.Opts.PidFileName, d.Opts.PidFilePerm); err != nil {
			return
		}
		if err = d.pidFile.Lock(); err != nil {
			return
		}
	}

	if len(d.Opts.LogFileName) > 0 {
		if d.logFile, err = os.OpenFile(d.Opts.LogFileName,
			os.O_WRONLY|os.O_CREATE|os.O_APPEND, d.Opts.LogFilePerm); err != nil {
			return
		}
	}

	return
}

func (d *Daemon) closeFiles() (err error) {
	cl := func(file **os.File) {
		if *file != nil {
			(*file).Close()
			*file = nil
		}
	}
	cl(&d.logFile)
	cl(&d.nullFile)
	if d.pidFile != nil {
		d.pidFile.Close()
		d.pidFile = nil
	}
	return
}

func (d *Daemon) prepareEnv() (err error) {
	if d.abspath, err = osext.Executable(); err != nil {
		return
	}

	if len(d.Args) == 0 {
		d.Args = os.Args
	}

	if len(d.Env) == 0 {
		d.Env = os.Environ()
	}
	d.Env = append(d.Env, d.Envs()...)

	return
}

func (d *Daemon) Envs() []string {
	mark := fmt.Sprintf("%s=%d", MARK_NAME, MARK_VALUE)
	return []string{mark}
}

func (d *Daemon) FdStart() int {
	return MARK_VALUE
}

func (d *Daemon) Files() (f []*os.File) {
	log := d.nullFile
	if d.logFile != nil {
		log = d.logFile
	}

	f = []*os.File{
		d.nullFile, // (0) stdin
		log,        // (1) stdout
		log,        // (2) stderr
		d.nullFile, // (3) dup on fd 0 after initialization
	}

	if d.pidFile != nil {
		f = append(f, d.pidFile.File) // (4) pid file
	}
	//fmt.Printf("Files() %#v\n", f)
	return
}

var initialized = false

func (d *Daemon) child() (err error) {
	if initialized {
		return os.ErrInvalid
	}
	initialized = true
	//d.nullFile = os.NewFile(0, "/dev/null")
	d.logFile = os.NewFile(1, d.Opts.LogFileName)
	d.errFile = os.NewFile(2, d.Opts.LogFileName)
	d.nullFile = os.NewFile(3, "/dev/null")

	fmt.Printf("%d child \n", os.Getpid())
	if len(d.Opts.PidFileName) > 0 {
		fmt.Printf("%d child lock file\n", os.Getpid())
		d.pidFile = NewLockFile(os.NewFile(4, d.Opts.PidFileName))
		if err = d.pidFile.WritePid(); err != nil {
			return
		}
	}

	if err = syscall.Close(0); err != nil {
		d.pidFile.Remove()
		return
	}
	if err = syscall.Dup2(3, 0); err != nil {
		d.pidFile.Remove()
		return
	}

	if d.Opts.Umask != 0 {
		syscall.Umask(int(d.Opts.Umask))
	}
	if len(d.Opts.Chroot) > 0 {
		err = syscall.Chroot(d.Opts.Chroot)
		if err != nil {
			d.pidFile.Remove()
			return
		}
	}

	return
}

func (d *Daemon) release() (err error) {
	if !initialized {
		return
	}
	fmt.Printf("daemon release %v\n", d.DontRemovePidFile)
	if d.pidFile != nil && !d.DontRemovePidFile {
		err = d.pidFile.Remove()
	}
	return
}

func (d *Daemon) BeforeRestart() {
}
func (d *Daemon) AfterRestart(pro *os.Process, err error) {
	if err != nil {
		return
	}
	d.DontRemovePidFile = true
}
