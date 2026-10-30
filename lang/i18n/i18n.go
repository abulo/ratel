package i18n

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gookit/ini/v2"
	"github.com/pkg/errors"
)

// 语言文件加载模式与消息渲染模式常量。
const (
	// FileMode 表示语言名即文件名，例如 "en" 对应 lang/en.ini。
	FileMode uint8 = 0
	// DirMode 表示语言名即目录名，加载该目录下全部 ini 文件，例如 "en" 对应 lang/en/*.ini。
	DirMode uint8 = 1

	// SprintfMode 使用 fmt.Sprintf 按顺序填充消息中的占位参数。
	SprintfMode uint8 = 0
	// ReplaceMode 使用字符串替换方式填充消息中的 {键} 占位参数。
	ReplaceMode uint8 = 1
)

// I18n 是多语言管理器，负责加载语言文件、注册语言并按语言翻译消息。
//
// 零值不可直接使用，须通过 New、NewEmpty、NewWithInit 或 Init 构造；
// 语言文件须先经 Init 加载完成后再调用翻译方法。
type I18n struct {
	// data 保存各语言已加载的 INI 数据。
	data map[string]*ini.Ini

	// langDir 为语言文件所在目录。
	langDir string
	// languages 为已注册语言映射，键为语言名，值为显示名（如 en:English、zh-CN:简体中文）。
	languages map[string]string

	// LoadMode 为语言文件加载模式，取值 FileMode 或 DirMode，Init 校验非法值时 panic。
	LoadMode uint8
	// TransMode 为消息参数渲染模式，取值 SprintfMode 或 ReplaceMode。
	TransMode uint8
	// DefaultLang 为默认语言名，Dt、DefTr 等方法按该语言翻译。
	DefaultLang string
	// FallbackLang 为备用语言名，当前语言缺失键时回退到该语言查找。
	FallbackLang string
}

// defI18n 为包级默认管理器实例。
var defI18n = NewEmpty()

// Default 返回包级默认的 I18n 实例。
func Default() *I18n {
	return defI18n
}

// T 通过默认实例按 lang 翻译 key，args 用于填充消息参数；等价于默认实例的 T 方法。
func T(lang string, key string, args ...any) string {
	return defI18n.T(lang, key, args...)
}

// Tr 通过默认实例按 lang 翻译 key，args 用于填充消息参数；等价于默认实例的 Tr 方法。
func Tr(lang string, key string, args ...any) string {
	return defI18n.Tr(lang, key, args...)
}

// Dt 通过默认实例按默认语言翻译 key，args 用于填充消息参数。
func Dt(key string, args ...any) string {
	return defI18n.DefTr(key, args...)
}

// DefTr 通过默认实例按默认语言翻译 key，args 用于填充消息参数。
func DefTr(key string, args ...any) string {
	return defI18n.DefTr(key, args...)
}

// Init 配置默认实例的语言目录、默认语言与语言列表，并按当前 LoadMode 加载语言文件；
// 返回加载完成的默认实例。
func Init(langDir string, defLang string, languages map[string]string) *I18n {
	defI18n.langDir = langDir
	defI18n.languages = languages
	defI18n.DefaultLang = defLang

	return defI18n.Init()
}

// New 按给定语言目录、默认语言与语言列表构造一个空实例，不会自动加载语言文件。
func New(langDir string, defLang string, languages map[string]string) *I18n {
	return &I18n{
		data:        make(map[string]*ini.Ini),
		langDir:     langDir,
		languages:   languages,
		DefaultLang: defLang,
	}
}

// NewEmpty 构造一个不含目录、默认语言与语言列表的空实例。
func NewEmpty() *I18n {
	return &I18n{
		data:      make(map[string]*ini.Ini),
		languages: make(map[string]string),
	}
}

// NewWithInit 构造实例并立即调用 Init 加载语言文件，返回加载完成的实例。
func NewWithInit(langDir string, defLang string, languages map[string]string) *I18n {
	m := New(langDir, defLang, languages)

	return m.Init()
}

// Dt 按默认语言翻译 key，args 用于填充消息参数。
func (l *I18n) Dt(key string, args ...any) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// DefTr 按默认语言翻译 key，args 用于填充消息参数。
func (l *I18n) DefTr(key string, args ...any) string {
	return l.Tr(l.DefaultLang, key, args...)
}

// T 按指定语言翻译 key，args 用于填充消息参数。
func (l *I18n) T(lang, key string, args ...any) string {
	return l.Tr(lang, key, args...)
}

// Tr 按指定语言翻译 key，args 用于填充消息参数。
//
// 语言未注册或当前语言缺失键时，回退到 FallbackLang 查找；备用语言仍取不到时
// 原样返回 key。键路径可使用节名前缀，例如 site.name 对应 INI 的 [site] 节。
func (l *I18n) Tr(lang, key string, args ...any) string {
	if !l.HasLang(lang) {
		msg := l.transFromFallback(key)
		if msg == "" {
			return key
		}

		if len(args) > 0 {
			msg = l.renderMessage(msg, args...)
		}
		return msg
	}

	msg, ok := l.data[lang].GetValue(key)
	if !ok {
		msg = l.transFromFallback(key)
	}

	if msg == "" {
		return key
	}

	if len(args) > 0 {
		msg = l.renderMessage(msg, args...)
	}
	return msg
}

// HasKey 判断已注册的 lang 语言下是否存在 key。
func (l *I18n) HasKey(lang, key string) (ok bool) {
	if !l.HasLang(lang) {
		return
	}

	_, ok = l.data[lang].GetValue(key)
	return
}

// transFromFallback 从备用语言中查找 key；备用语言未注册时返回空串。
func (l *I18n) transFromFallback(key string) string {
	fl := l.FallbackLang
	if !l.HasLang(fl) {
		return ""
	}

	return l.data[fl].String(key)
}

// renderMessage 按 TransMode 用 args 填充消息 msg；ReplaceMode 支持 []string 替换对、
// map[string]any 以及交替出现的键、值参数。
func (l *I18n) renderMessage(msg string, args ...any) string {
	if l.TransMode == SprintfMode {
		return fmt.Sprintf(msg, args...)
	}

	if ss, ok := args[0].([]string); ok {
		return strings.NewReplacer(ss...).Replace(msg)
	}

	var ss []string

	if mp, ok := args[0].(map[string]any); ok {
		for k, v := range mp {
			ss = append(ss, "{"+k+"}")
			ss = append(ss, toString(v))
		}
	} else {
		// 交替传入 键、值、键、值……，偶数下标补充为 {键} 形式。
		for i, val := range args {
			str := toString(val)
			if i%2 == 0 {
				str = "{" + str + "}"
			}

			ss = append(ss, str)
		}
	}

	return strings.NewReplacer(ss...).Replace(msg)
}

// toString 将参数转换为字符串；不支持的类型返回占位标记 "CANNOT-TO-STRING"。
func toString(val any) (str string) {
	switch tVal := val.(type) {
	case int:
		str = strconv.Itoa(tVal)
	case int8:
		str = strconv.Itoa(int(tVal))
	case int16:
		str = strconv.Itoa(int(tVal))
	case int32:
		str = strconv.Itoa(int(tVal))
	case int64:
		str = strconv.Itoa(int(tVal))
	case uint:
		str = strconv.Itoa(int(tVal))
	case uint8:
		str = strconv.Itoa(int(tVal))
	case uint16:
		str = strconv.Itoa(int(tVal))
	case uint32:
		str = strconv.Itoa(int(tVal))
	case uint64:
		str = strconv.Itoa(int(tVal))
	case float32:
		str = fmt.Sprint(tVal)
	case float64:
		str = fmt.Sprint(tVal)
	case string:
		str = tVal
	case nil:
		str = ""
	default:
		str = "CANNOT-TO-STRING"
	}
	return
}

// Init 按 LoadMode 加载全部已注册语言的文件；LoadMode 取值非法时直接 panic，
// 语言文件加载失败时同样 panic；返回实例自身以便链式调用。
func (l *I18n) Init() *I18n {
	if l.LoadMode == FileMode {
		l.loadSingleFiles()
	} else if l.LoadMode == DirMode {
		l.loadDirFiles()
	} else {
		panic("invalid load mode setting. only allow 0, 1")
	}

	return l
}

// loadSingleFiles 在 FileMode 模式下逐个加载各语言对应的单个 ini 文件。
func (l *I18n) loadSingleFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		lData := ini.New()
		err := lData.LoadFiles(l.langDir + pathSep + lang + ".ini")
		if err != nil {
			panic("fail to load language: " + lang + ", error " + err.Error())
		}

		l.data[lang] = lData
	}
}

// loadDirFiles 在 DirMode 模式下遍历各语言目录，加载其中全部 ini 文件。
func (l *I18n) loadDirFiles() {
	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		dirPath := l.langDir + pathSep + lang
		files, err := os.ReadDir(dirPath)
		if err != nil {
			panic("read dir fail: " + dirPath + ", error " + err.Error())
		}

		sl := l.data[lang]

		for _, fi := range files {
			if fi.IsDir() || !strings.HasSuffix(fi.Name(), ".ini") {
				continue
			}

			var err error
			if sl != nil {
				err = sl.LoadFiles(dirPath + pathSep + fi.Name())
			} else {
				sl = ini.New()
				err = sl.LoadFiles(dirPath + pathSep + fi.Name())
				l.data[lang] = sl
			}

			if err != nil {
				panic("fail to load language file: " + lang + ", error " + err.Error())
			}
		}
	}
}

// Add 注册新语言，是 NewLang 的别名。
func (l *I18n) Add(lang string, name string) {
	l.NewLang(lang, name)
}

// NewLang 注册一个新语言及其显示名；语言已存在时直接忽略，不覆盖原有数据。
func (l *I18n) NewLang(lang string, name string) {
	if _, ok := l.languages[lang]; ok {
		return
	}

	l.data[lang] = ini.New()
	l.languages[lang] = name
}

// LoadFile 将 ini 文件内容追加到已注册的 lang 语言；语言未注册或文件加载失败时返回错误。
func (l *I18n) LoadFile(lang string, file string) (err error) {
	if ld, ok := l.data[lang]; ok {
		err = ld.LoadFiles(file)
		if err != nil {
			return
		}
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// LoadString 将 INI 格式字符串内容追加到已注册的 lang 语言；
// 语言未注册或内容解析失败时返回错误。
func (l *I18n) LoadString(lang string, data string) (err error) {
	if ld, ok := l.data[lang]; ok {
		err = ld.LoadStrings(data)
		if err != nil {
			return
		}
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// Lang 返回 lang 语言对应的 INI 实例；语言未注册时返回 nil。
func (l *I18n) Lang(lang string) *ini.Ini {
	if _, ok := l.languages[lang]; ok {
		return l.data[lang]
	}

	return nil
}

// Export 将 lang 语言的数据导出为 INI 格式字符串；语言未注册时返回空串，写出失败时 panic。
func (l *I18n) Export(lang string) string {
	if _, ok := l.languages[lang]; !ok {
		return ""
	}

	var buf bytes.Buffer

	_, err := l.data[lang].WriteTo(&buf)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// HasLang 判断管理器中是否已注册 lang 语言。
func (l *I18n) HasLang(lang string) bool {
	_, ok := l.languages[lang]
	return ok
}

// DelLang 移除指定语言及其翻译数据，返回是否确实移除了已存在的语言。
func (l *I18n) DelLang(lang string) bool {
	_, ok := l.languages[lang]
	if ok {
		delete(l.data, lang)
		delete(l.languages, lang)
	}

	return ok
}

// Languages 返回已注册语言映射（语言名到显示名）；返回的是内部 map，调用方不应修改。
func (l *I18n) Languages() map[string]string {
	return l.languages
}
