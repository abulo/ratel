package watch

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// MinWatcherFrequency 为监视器更新频率允许的最小值，低于该值的配置会被判定为非法。
const MinWatcherFrequency = 1 * time.Second

// Options 保存热编译服务的全部配置项。
//
// 零值不可直接使用，须先调用 Watch 由 sanitize 完成校验与默认值填充；
// 配置应在启动监视前一次性写好，运行期间不保证并发安全。
type Options struct {
	// Logger 为日志输出对象；为空时初始化为控制台日志器。
	Logger Logger

	// AutoTidy 控制 go.mod 发生变化时是否自动执行 go mod tidy。
	AutoTidy bool

	// MainFiles 为传给 go build 的文件参数；为空表示编译当前目录。
	MainFiles string

	// OutputName 为可执行文件输出路径。
	//
	// 为空时取主目录名；不含路径信息时附加到第一个监视目录下；Windows 下无需指定
	// .exe 后缀；若带路径信息，则以该文件所在目录作为程序的运行工作目录。
	OutputName string
	appName    string

	// Flags 为传给各工具链的 flags 参数，键对应 asmflags、gccgoflags、gcflags、ldflags 系列前缀。
	Flags Flags

	// Exts 为需要监视的文件扩展名列表；为空表示不监视任何文件，指定 "*" 表示监视所有文件类型（含无扩展名文件）。
	Exts    []string
	anyExts bool

	// Excludes 为忽略文件的路径模式列表，采用 filepath.Match 匹配语法。
	Excludes []string

	// AppArgs 为传给编译后程序启动命令的参数字符串。
	AppArgs string
	appArgs []string

	// Recursive 控制是否递归监视子目录。
	Recursive bool

	// Dirs 为需要监视的目录列表，至少指定一个；第一个目录为主目录，同时作为默认工作目录。
	Dirs  []string
	paths []string

	// WatcherFrequency 为判定文件更新的时间间隔，间隔内的频繁修改会被忽略，避免重复编译。
	//
	// 不能小于 MinWatcherFrequency；为 0 时取默认值 MinWatcherFrequency。
	WatcherFrequency time.Duration

	goCmdArgs []string
}

// Flags 保存传递给 go build 的各类 flags 参数，键为工具链前缀（asm、gccgo、gc、ld），值为对应参数。
type Flags map[string]string

// sanitize 校验配置合法性、填充默认值并展开内部派生字段；校验失败时返回对应错误。
func (opt *Options) sanitize() error {

	if opt.Logger == nil {
		opt.Logger = NewConsoleLogger(true, os.Stderr, os.Stdout)
	}

	// 校验 Excludes 中的模式是否为合法的 filepath.Match 语法。
	for _, p := range opt.Excludes {
		if _, err := filepath.Match(p, "abc"); err != nil {
			return err
		}
	}

	if len(opt.Dirs) == 0 {
		return errors.New("字段 Dirs 不能为空")
	}
	wd, err := filepath.Abs(opt.Dirs[0])
	if err != nil {
		return err
	}
	opt.Dirs[0] = wd

	if opt.appName, err = getAppName(opt.OutputName, wd); err != nil {
		return err
	}

	opt.sanitizeExts()

	opt.appArgs = splitArgs(opt.AppArgs)

	if opt.paths, err = recursivePaths(opt.Recursive, opt.Dirs); err != nil {
		return err
	}

	if opt.WatcherFrequency == 0 {
		opt.WatcherFrequency = MinWatcherFrequency
	} else if opt.WatcherFrequency < MinWatcherFrequency {
		return errors.New("watcherFrequency 值过小")
	}

	args := []string{"build", "-o", opt.appName}
	for k, v := range opt.Flags {
		args = append(args, "-"+k+"flags", v)
	}
	args = append(args, "-v")
	if len(opt.MainFiles) > 0 {
		args = append(args, opt.MainFiles)
	}
	opt.goCmdArgs = args

	return nil
}

func (opt *Options) sanitizeExts() {
	exts := make([]string, 0, len(opt.Exts))
	for _, ext := range opt.Exts {
		ext = strings.TrimSpace(ext)
		if len(ext) == 0 {
			continue
		}

		if ext == "*" {
			opt.anyExts = true
			return
		}

		if ext[0] != '.' {
			ext = "." + ext
		}
		exts = append(exts, ext)
	}
	opt.Exts = exts
}

func getAppName(outputName, wd string) (string, error) {
	if outputName == "" {
		outputName = filepath.Base(wd)
	}

	goexe := os.Getenv("GOEXE")
	if goexe != "" && !strings.HasSuffix(outputName, goexe) {
		outputName += goexe
	}

	// 不含路径分隔符时视为纯文件名，拼接到主目录下。
	if strings.IndexByte(outputName, '/') < 0 || strings.IndexByte(outputName, filepath.Separator) < 0 {
		outputName = filepath.Join(wd, outputName)
	}

	outputName, err := filepath.Abs(outputName)
	if err != nil {
		return "", err
	}

	return outputName, nil
}

// recursivePaths 按 recursive 决定是否展开 paths 中每个目录的子目录；点开头的隐藏目录被跳过。
func recursivePaths(recursive bool, paths []string) ([]string, error) {
	if !recursive {
		return paths, nil
	}

	ret := []string{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			return nil
		}

		if fi.Name()[0] == '.' {
			return fs.SkipDir
		}
		ret = append(ret, path)
		return nil
	}

	for _, path := range paths {
		if err := filepath.Walk(path, walk); err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func splitArgs(args string) []string {
	ret := make([]string, 0, 10)
	var state byte
	var start, index int

	for index = 0; index < len(args); index++ {
		b := args[index]
		switch b {
		case ' ':
			if state == '"' {
				break
			}

			if state != ' ' {
				ret = appendArg(ret, args[start:index])
				state = ' '
			}
			start = index + 1
		case '=':
			if state == '"' {
				break
			}

			if state != '=' {
				ret = appendArg(ret, args[start:index])
				state = '='
			}
			start = index + 1
			state = 0
		case '"':
			if state == '"' {
				ret = appendArg(ret, args[start:index])
				state = 0
				start = index + 1
				break
			}

			if start != index {
				ret = appendArg(ret, args[start:index])
			}
			state = '"'
			start = index + 1
		default:
			if state == ' ' {
				state = 0
				start = index
			}
		}
	}

	if start < len(args) {
		ret = appendArg(ret, args[start:])
	}

	return ret
}

func appendArg(args []string, arg string) []string {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return args
	}

	return append(args, arg)
}
