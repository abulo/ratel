## 功能简介

- 使用简单，支持加载多个语言，多个文件
- 两种数据加载模式：单文件 `FileMode` 、文件夹 `DirMode`；默认是文件夹模式
- 支持设置默认语言，备用语言；当在默认语言数据没找到时，自动尝试到备用语言查找
- 支持参数替换，也有两种模式：`SprintfMode` 通过 `fmt.Sprintf` 替换参数，`ReplaceMode` 则使用 `strings.Replacer` 替换

## 快速使用

```text
lang/
    en/
        default.ini
        ...
    zh-CN/
        default.ini
        ...
```

### 初始化

```go
    import "github.com/abulo/ratel/v1/i18n"

    languages := map[string]string{
        "en": "English",
        "zh-CN": "简体中文",
        // "zh-TW": "繁体中文",
    }

    // 这里直接初始化的默认实例
    i18n.Init("conf/lang", "en", languages)

    // 或者创建自定义的新实例
    // i18n.New(langDir string, defLang string, languages)
    // i18n.NewEmpty()
```


### 翻译数据

```go
    // 从指定的语言翻译
    msg := i18n.Tr("en", "key")

    // 从默认语言翻译
    msg = i18n.DefTr("key")
    // with arguments. 
    msg = i18n.DefTr("key1", "arg1", "arg2")
```


### 参数替换模式

使用 `SprintfMode`(**defaults**) 模式:

```ini
# en.ini
desc = I am %s, age is %d
```

按kv顺序传入参数使用：

```go
msg := i18n.Tr("en", "desc", "name", "tom", "age", 22)
// Output: "I am tom, age is 22"
```

使用 `ReplaceMode` 替换模式:

```ini
# en.ini
desc = I am {name}, age is {age}
```

传入 `map[string]interface{}` 参数使用：

```go
i18n.TransMode = i18n.ReplaceMode

msg := i18n.Tr("en", "desc", "desc", map[string]interface{}{
    "name": "tom",
    "age": 22,
})
// Output: "I am tom, age is 22"
```