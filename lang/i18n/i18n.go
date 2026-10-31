package i18n

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

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
	// data 保存各语言已加载的 INI 数据；Lang()/Export() 等对外接口仍返回原始实例。
	data map[string]*ini.Ini

	// snap 为 data 的扁平化只读快照（PERF-9-17 [P1]）：翻译查询直接读自有 map，
	// 绕开第三方 ini 在只读查询上使用的排他写锁。
	snap map[string]map[string]string

	// mu 保护 data / languages / snap 的读写（PERF-9-20 [P2]）：翻译读路径读锁，动态加载写路径写锁。
	mu sync.RWMutex

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
		snap:        make(map[string]map[string]string),
		langDir:     langDir,
		languages:   languages,
		DefaultLang: defLang,
	}
}

// NewEmpty 构造一个不含目录、默认语言与语言列表的空实例。
func NewEmpty() *I18n {
	return &I18n{
		data:      make(map[string]*ini.Ini),
		snap:      make(map[string]map[string]string),
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
//
// PERF-9-17 [P1]：查询直接走自有扁平化快照，不经第三方 ini 的排他写锁；
// PERF-9-19 [P2]：以快照存在性合并语言注册与键查找两次 map 查询。
func (l *I18n) Tr(lang, key string, args ...any) string {
	l.mu.RLock()
	snap, registered := l.snap[lang]
	var msg string
	ok := false
	if registered {
		msg, ok = lookupSnapshot(snap, key)
	}
	if !ok {
		msg = l.lookupFallbackLocked(key)
	}
	l.mu.RUnlock()

	if msg == "" {
		return key
	}

	if len(args) > 0 {
		msg = l.renderMessage(msg, args...)
	}
	return msg
}

// HasKey 判断已注册的 lang 语言下是否存在 key。
//
// PERF-9-19 [P2]：直接以快照非空判定语言注册，合并 languages 与 data 两次 map 查找。
func (l *I18n) HasKey(lang, key string) (ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	snap, registered := l.snap[lang]
	if !registered {
		return false
	}
	_, ok = lookupSnapshot(snap, key)
	return
}

// transFromFallback 从备用语言中查找 key；备用语言未注册时返回空串。
func (l *I18n) transFromFallback(key string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.lookupFallbackLocked(key)
}

// lookupFallbackLocked 在已持锁状态下从备用语言快照查找 key。
func (l *I18n) lookupFallbackLocked(key string) string {
	fsnap, registered := l.snap[l.FallbackLang]
	if !registered {
		return ""
	}
	msg, ok := lookupSnapshot(fsnap, key)
	if !ok {
		return ""
	}
	return msg
}

// lookupSnapshot 复刻第三方 ini.GetValue 的 formatKey 与 splitSectionAndKey 查找语义，
// 在自有扁平化快照上查询；节内键以 "节.键" 为复合键，默认节键保持裸键名。
func lookupSnapshot(snap map[string]string, key string) (string, bool) {
	key = strings.Trim(strings.TrimSpace(key), ini.SepSection)
	if key == "" {
		return "", false
	}

	if strings.Contains(key, ini.SepSection) {
		ss := strings.SplitN(key, ini.SepSection, 2)
		section, k := strings.TrimSpace(ss[0]), strings.TrimSpace(ss[1])
		v, ok := snap[section+ini.SepSection+k]
		return v, ok
	}

	v, ok := snap[key]
	return v, ok
}

// renderMessage 按 TransMode 用 args 填充消息 msg；ReplaceMode 支持 []string 替换对、
// map[string]any 以及交替出现的键、值参数。
func (l *I18n) renderMessage(msg string, args ...any) string {
	if l.TransMode == SprintfMode {
		return fmt.Sprintf(msg, args...)
	}

	if ss, ok := args[0].([]string); ok {
		// PERF-9-18 [P1]（部分实施）：[]string 为通用原始替换对，保留 NewReplacer 维持
		// 其前缀/最长匹配语义；map 与交替参数分支改为单遍占位符扫描，不再每次编译替换 trie。
		return strings.NewReplacer(ss...).Replace(msg)
	}

	pairs := make(map[string]string)

	if mp, ok := args[0].(map[string]any); ok {
		for k, v := range mp {
			pairs["{"+k+"}"] = toString(v)
		}
	} else if len(args)%2 == 1 {
		// PERF-9-18 [P1]：奇数个交替参数沿用 NewReplacer，保留原有的参数计数校验行为。
		var ss []string
		for i, val := range args {
			str := toString(val)
			if i%2 == 0 {
				str = "{" + str + "}"
			}
			ss = append(ss, str)
		}
		return strings.NewReplacer(ss...).Replace(msg)
	} else {
		var pending string
		for i, val := range args {
			str := toString(val)
			if i%2 == 0 {
				pending = "{" + str + "}"
			} else {
				pairs[pending] = str
			}
		}
	}

	return replacePlaceholders(msg, pairs)
}

// replacePlaceholders 单遍扫描 msg 中的 {占位符} 并按 pairs 替换；未匹配到的花括号原样保留。
//
// PERF-9-18 [P1]：替代 strings.NewReplacer 的每次 trie 编译；占位符由花括号定界，
// 不存在键前缀重叠问题，替换结果与 NewReplacer 在 {key} 形式占位符上一致。
func replacePlaceholders(msg string, pairs map[string]string) string {
	var b strings.Builder
	b.Grow(len(msg))

	for i := 0; i < len(msg); {
		if msg[i] != '{' {
			b.WriteByte(msg[i])
			i++
			continue
		}

		rest := msg[i+1:]
		closeIdx := strings.IndexByte(rest, '}')
		if closeIdx < 0 {
			b.WriteByte('{')
			i++
			continue
		}

		placeholder := rest[:closeIdx+1]
		if val, ok := pairs["{"+placeholder[:closeIdx]+"}"]; ok {
			b.WriteString(val)
			i += closeIdx + 2
		} else {
			b.WriteByte('{')
			i++
		}
	}

	return b.String()
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
	l.mu.Lock()
	defer l.mu.Unlock()

	pathSep := string(os.PathSeparator)

	for lang := range l.languages {
		lData := ini.New()
		err := lData.LoadFiles(l.langDir + pathSep + lang + ".ini")
		if err != nil {
			panic("fail to load language: " + lang + ", error " + err.Error())
		}

		l.data[lang] = lData
		l.refreshSnapshotLocked(lang)
	}
}

// loadDirFiles 在 DirMode 模式下遍历各语言目录，加载其中全部 ini 文件。
func (l *I18n) loadDirFiles() {
	l.mu.Lock()
	defer l.mu.Unlock()

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
		l.refreshSnapshotLocked(lang)
	}
}

// refreshSnapshotLocked 基于原始 ini 实例重建指定语言的扁平化快照；调用方须已持写锁。
func (l *I18n) refreshSnapshotLocked(lang string) {
	ld := l.data[lang]
	if ld == nil {
		l.snap[lang] = nil
		return
	}
	l.snap[lang] = snapshotLang(ld)
}

// snapshotLang 把一个 ini 实例的全部节与键扁平化为 "节.键" → 值 的只读 map。
func snapshotLang(ld *ini.Ini) map[string]string {
	defSec := ld.DefSection()
	snap := make(map[string]string)
	for secName, sec := range ld.Data() {
		for k, v := range sec {
			if secName == defSec {
				snap[k] = v
			} else {
				snap[secName+ini.SepSection+k] = v
			}
		}
	}
	return snap
}

// Add 注册新语言，是 NewLang 的别名。
func (l *I18n) Add(lang string, name string) {
	l.NewLang(lang, name)
}

// NewLang 注册一个新语言及其显示名；语言已存在时直接忽略，不覆盖原有数据。
func (l *I18n) NewLang(lang string, name string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.languages[lang]; ok {
		return
	}

	l.data[lang] = ini.New()
	l.languages[lang] = name
	l.snap[lang] = map[string]string{}
}

// LoadFile 将 ini 文件内容追加到已注册的 lang 语言；语言未注册或文件加载失败时返回错误。
func (l *I18n) LoadFile(lang string, file string) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ld, ok := l.data[lang]; ok {
		err = ld.LoadFiles(file)
		if err != nil {
			return
		}
		l.refreshSnapshotLocked(lang)
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// LoadString 将 INI 格式字符串内容追加到已注册的 lang 语言；
// 语言未注册或内容解析失败时返回错误。
func (l *I18n) LoadString(lang string, data string) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if ld, ok := l.data[lang]; ok {
		err = ld.LoadStrings(data)
		if err != nil {
			return
		}
		l.refreshSnapshotLocked(lang)
	} else {
		err = errors.New("language" + lang + " not exist, please create it before load data")
	}

	return
}

// Lang 返回 lang 语言对应的 INI 实例；语言未注册时返回 nil。
func (l *I18n) Lang(lang string) *ini.Ini {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if _, ok := l.languages[lang]; ok {
		return l.data[lang]
	}

	return nil
}

// Export 将 lang 语言的数据导出为 INI 格式字符串；语言未注册时返回空串，写出失败时 panic。
func (l *I18n) Export(lang string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()

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
	l.mu.RLock()
	defer l.mu.RUnlock()

	_, ok := l.languages[lang]
	return ok
}

// DelLang 移除指定语言及其翻译数据，返回是否确实移除了已存在的语言。
func (l *I18n) DelLang(lang string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, ok := l.languages[lang]
	if ok {
		delete(l.data, lang)
		delete(l.languages, lang)
		delete(l.snap, lang)
	}

	return ok
}

// Languages 返回已注册语言映射（语言名到显示名）；返回的是内部 map，调用方不应修改。
func (l *I18n) Languages() map[string]string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.languages
}
