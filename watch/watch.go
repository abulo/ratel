package watch

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	viper "github.com/abulo/ratel/v2/config"
	"github.com/fsnotify/fsnotify"
)

var configFile string = "watch.toml"

//AppPath 运行路径
func AppPath() string {

	dir, err := os.Getwd()
	if err != nil {
		Fatalf("不能获取程序执行的目录: [ %s ]", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)

}

//ParseConfig 解析配置文件
func ParseConfig() *viper.Viper {

	AppConfig := viper.New()
	AppConfig.SetConfigName("watch")
	AppConfig.SetConfigType("toml")
	configFile := AppPath() + "/" + configFile
	configFile = strings.Replace(configFile, "\\", "/", -1)
	if !FileExist(configFile) {
		Fatalf("配置文件不存在[ %s ]\n", configFile)
	}
	AppConfig.SetConfigFile(configFile)

	err := AppConfig.ReadInConfig() // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	return AppConfig
}

//FileExist 判断文件是否存在
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

type (
	Level int
)

const (
	LevelFatal = iota
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

var _log *logger = New()

func Fatal(s string) {
	_log.Output(LevelFatal, s)
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_log.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Error(s string) {
	_log.Output(LevelError, s)
}

func Errorf(format string, v ...interface{}) {
	_log.Output(LevelError, fmt.Sprintf(format, v...))
}

func Warn(s string) {
	_log.Output(LevelWarning, s)
}

func Warnf(format string, v ...interface{}) {
	_log.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func Info(s string) {
	_log.Output(LevelInfo, s)
}

func Infof(format string, v ...interface{}) {
	_log.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func Debug(s string) {
	_log.Output(LevelDebug, s)
}

func Debugf(format string, v ...interface{}) {
	_log.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func SetLogLevel(level Level) {
	_log.SetLogLevel(level)
}

type logger struct {
	_log *log.Logger
	//小于等于该级别的level才会被记录
	logLevel Level
}

//NewLogger 实例化，供自定义
func NewLogger() *logger {
	return &logger{_log: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
}

//New 实例化，供外部直接调用 log.XXXX
func New() *logger {
	return &logger{_log: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags), logLevel: LevelDebug}
}

func (l *logger) Output(level Level, s string) error {
	if l.logLevel < level {
		return nil
	}
	formatStr := "[UNKNOWN] %s"
	switch level {
	case LevelFatal:
		formatStr = "\033[35m[FATAL]\033[0m %s"
	case LevelError:
		formatStr = "\033[31m[ERROR]\033[0m %s"
	case LevelWarning:
		formatStr = "\033[33m[WARN]\033[0m %s"
	case LevelInfo:
		formatStr = "\033[32m[INFO]\033[0m %s"
	case LevelDebug:
		formatStr = "\033[36m[DEBUG]\033[0m %s"
	}
	s = fmt.Sprintf(formatStr, s)
	return l._log.Output(3, s)
}

func (l *logger) Fatal(s string) {
	l.Output(LevelFatal, s)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *logger) Error(s string) {
	l.Output(LevelError, s)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

func (l *logger) Warn(s string) {
	l.Output(LevelWarning, s)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarning, fmt.Sprintf(format, v...))
}

func (l *logger) Info(s string) {
	l.Output(LevelInfo, s)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

func (l *logger) Debug(s string) {
	l.Output(LevelDebug, s)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

func (l *logger) SetLogLevel(level Level) {
	l.logLevel = level
}

var (
	cmd          *exec.Cmd
	state        sync.Mutex
	eventTime    = make(map[string]int64)
	scheduleTime time.Time
)

func readAppDirectories(directory string, paths *[]string) {
	fileInfos, err := os.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		// fmt.Println(fileInfo.Name())
		if ok := ignoreDirFunc(fileInfo); ok {
			continue
		}
		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readAppDirectories(directory+"/"+fileInfo.Name(), paths)
			continue
		}
		if useDirectory == true {
			continue
		}
		*paths = append(*paths, directory)
		useDirectory = true
	}
	return
}

func ignoreDirFunc(f fs.DirEntry) bool {
	ignoredir := cfg.Strings("ignoredir")
	if len(ignoredir) == 0 {
		return false
	}
	appPath := AppPath()
	fullpath, _ := filepath.Abs(filepath.Dir(f.Name()))
	fullpath = filepath.FromSlash(strings.Replace(fullpath+"/"+f.Name(), "\\", "/", -1))
	for _, dir := range ignoredir {
		runFile := filepath.FromSlash(strings.Replace(appPath+"/"+dir, "\\", "/", -1))
		if ok := strings.HasPrefix(fullpath, runFile); ok {
			return true
		}
	}
	return false
}

func ignoreFileFun(file string) bool {
	ignorefile := cfg.Strings("ignorefile")
	if len(ignorefile) == 0 {
		return false
	}
	//要判断文件的绝对路径
	fullpath := file
	for _, fileExtra := range ignorefile {
		if ok := strings.HasSuffix(fullpath, fileExtra); ok {
			return true
		}
	}
	return false
}

// getFileModTime retuens unix timestamp of `os.File.ModTime` by given path.
func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		Errorf("文件打开失败[ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		Errorf("文件信息获取失败[ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

//NewWatcher new watcher
func NewWatcher(paths []string, files []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Errorf("创建监控失败[ %s ]\n", err)
		os.Exit(2)
	}
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				isbuild := true
				// Skip ignored files
				if ignoreFileFun(e.Name) {
					continue
				}
				mt := getFileModTime(e.Name)
				if t := eventTime[e.Name]; mt == t {
					//log.Infof("[SKIP] # %s #\n", e.String())
					isbuild = false
				}
				eventTime[e.Name] = mt
				if isbuild {
					go func() {
						scheduleTime = time.Now().Add(1 * time.Second)
						for {
							time.Sleep(scheduleTime.Sub(time.Now()))
							if time.Now().After(scheduleTime) {
								break
							}
							return
						}
						Autobuild(files)
					}()
				}
			case err := <-watcher.Errors:
				Errorf("%v", err)
				Warnf(" %s\n", err.Error()) // No need to exit here
			}
		}
	}()
	Infof("初始化监控\n")
	for _, path := range paths {
		Infof("文件夹( %s )\n", path)
		err = watcher.Add(path)
		if err != nil {
			Errorf("讲课文件夹失败[ %s ]\n", err)
			os.Exit(2)
		}
	}
}

//Autobuild auto build
func Autobuild(files []string) {
	state.Lock()
	defer state.Unlock()
	Infof("开始构建...\n")
	if err := os.Chdir(currpath); err != nil {
		Errorf("目录读取失败: %+v\n", err)
		return
	}

	cmdName := "go"
	var err error
	var ostype = runtime.GOOS
	var build string
	buildPrifix := cfg.String("build")
	if len(buildPrifix) < 1 {
		build = currpath
	} else {
		build = currpath + "/" + buildPrifix
	}
	filenameWithSuffix := path.Base(runcmd)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名
	args := []string{"build"}
	if ostype == "windows" {
		build = build + "/" + filenameOnly + ".exe"
	} else {
		build = build + "/" + filenameOnly
	}
	build = filepath.FromSlash(strings.Replace(build, "\\", "/", -1))
	Output = build
	args = append(args, "-o", Output, runcmd)
	bcmd := exec.Command(cmdName, args...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr

	Infof("编译参数: %s %s", cmdName, strings.Join(args, " "))
	err = bcmd.Run()
	if err != nil {
		Errorf("============== 编译失败 ===================\n")
		Errorf("%+v\n", err)
		return
	}
	Infof("编译成功\n")
	Restart(Output)

}

//Kill kill process
func Kill() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Kill.recover -> ", e)
		}
	}()
	if cmd != nil && cmd.Process != nil {
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println("Kill -> ", err)
		}
	}
}

//Restart restart app
func Restart(appname string) {
	// Debugf("杀掉进程")
	Kill()
	go Start(appname)
}

//Start start app
func Start(appname string) {
	Infof("开始运行 %s ...\n", appname)

	cmd = exec.Command(appname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Args = append([]string{appname}, cmdArgs...)

	Infof("Run %s", strings.Join(cmd.Args, " "))
	go cmd.Run()
	Infof("%s 运行中...\n", appname)
	started <- true
}

// end watch

///main starts
var (
	cfg      *viper.Viper
	currpath string
	exit     chan bool
	runcmd   string
	runArgs  string
	cmdArgs  []string
	Output   string
	started  chan bool
)

func runApp() {
	var paths []string
	files := []string{}
	readAppDirectories(currpath, &paths)
	NewWatcher(paths, files)
	Autobuild(files)
	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}

//Run 运行
func Run(cmd, args string) {
	runcmd = cmd
	runArgs = args
	//config配置
	cfg = ParseConfig()
	//程序运行的当前目录
	currpath = AppPath()

	if len(runcmd) < 1 {
		Fatalf("要运行的Go文件没有填写\n")
	}

	if runArgs != "" {
		cmdArgs = strings.Split(runArgs, ",")
	}

	runFile := currpath + "/" + runcmd
	runFile = strings.Replace(runFile, "\\", "/", -1)
	ok := strings.HasSuffix(runFile, ".go")
	if !ok {
		Fatalf("路径错误[ %s ]\n", runFile)
	}
	if !FileExist(runFile) {
		Fatalf("文件不存在[ %s ]\n", runFile)
	}
	runcmd = runFile

	runApp()
}
