package change

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var CmdNew = &cobra.Command{
	Use:   "changelog",
	Short: "更新日志",
	Long:  "更新日志: toolkit changelog",
	Run:   Run,
}

func getTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "--sort=-v:refname")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	tags := strings.Split(strings.TrimSpace(out.String()), "\n")
	return tags, nil
}

// getTagDate returns the date of a specific tag
func getTagDate(tag string) (string, error) {
	cmd := exec.Command("git", "log", "-1", "--format=%ad", "--date=short", tag)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	date := strings.TrimSpace(out.String())
	return date, nil
}

// getCommits returns a list of commits between two tags
func getCommits(tag1, tag2 string) ([]string, error) {
	var cmd *exec.Cmd
	if tag2 == "" {
		cmd = exec.Command("git", "log", "--pretty=format:%s %h", tag1)
	} else {
		cmd = exec.Command("git", "log", "--pretty=format:%s %h", fmt.Sprintf("%s...%s", tag2, tag1))
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	commits := strings.Split(strings.TrimSpace(out.String()), "\n")
	return commits, nil
}

// getRepoURL returns the URL of the remote repository
func getRepoURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	url := strings.TrimSpace(out.String())

	// Convert SSH URL to HTTPS URL if necessary
	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@github.com", "https://github.com", 1)
	}
	// Ensure it ends with .git and remove it for display purposes
	url = strings.TrimSuffix(url, ".git")
	return url, nil
}

// classifyCommits classifies commits into different types
func classifyCommits(commits []string, repoURL string) map[string][]string {
	classified := map[string][]string{
		"feat✨:引入新功能":            {},
		"fix🐛:修复 bug":            {},
		"style💄:更新UI样式文件":        {},
		"format🥚:格式化代码":          {},
		"docs📝:添加/更新文档":          {},
		"perf👌:提高性能/优化":          {},
		"test✅:增加测试代码":           {},
		"refactor🎨:改进代码结构/代码格式":  {},
		"patch🚑:添加重要补丁":          {},
		"file📦:添加新文件":            {},
		"publish🚀:发布新版本":         {},
		"tag📌:发布新标签":             {},
		"config🔧:修改配置文件":         {},
		"git🙈:添加或修改.gitignore文件": {},
		"others:其他":              {},
		"init🎉:初次提交/初始化项目":       {},
	}

	re := regexp.MustCompile(`^(feat|fix|style|perf|docs|init|test|refactor|patch|file|publish|tag|config|git)([^:]*):(.+)([a-f0-9]{7,40})$`)
	// match := re.FindStringSubmatch(commit)

	for _, commit := range commits {
		match := re.FindStringSubmatch(commit)
		if len(match) > 4 {
			commitType, emoji, message, hash := match[1], match[2], match[3], match[4]
			formattedCommit := fmt.Sprintf("%s%s:%s ([%s](%s/commit/%s))", commitType, emoji, message, hash, repoURL, hash)
			switch commitType {
			case "feat":
				classified["feat✨:引入新功能"] = append(classified["feat✨:引入新功能"], formattedCommit)
			case "fix":
				classified["fix🐛:修复 bug"] = append(classified["fix🐛:修复 bug"], formattedCommit)
			case "style":
				classified["style💄:更新UI样式文件"] = append(classified["style💄:更新UI样式文件"], formattedCommit)
			case "format":
				classified["format🥚:格式化代码"] = append(classified["format🥚:格式化代码"], formattedCommit)
			case "docs":
				classified["docs📝:添加/更新文档"] = append(classified["docs📝:添加/更新文档"], formattedCommit)
			case "perf":
				classified["perf👌:提高性能/优化"] = append(classified["perf👌:提高性能/优化"], formattedCommit)
			case "init":
				classified["init🎉:初次提交/初始化项目"] = append(classified["init🎉:初次提交/初始化项目"], formattedCommit)
			case "test":
				classified["test✅:增加测试代码"] = append(classified["test✅:增加测试代码"], formattedCommit)
			case "refactor":
				classified["refactor🎨:改进代码结构/代码格式"] = append(classified["refactor🎨:改进代码结构/代码格式"], formattedCommit)
			case "patch":
				classified["patch🚑:添加重要补丁"] = append(classified["patch🚑:添加重要补丁"], formattedCommit)
			case "file":
				classified["file📦:添加新文件"] = append(classified["file📦:添加新文件"], formattedCommit)
			case "publish":
				classified["publish🚀:发布新版本"] = append(classified["publish🚀:发布新版本"], formattedCommit)
			case "tag":
				classified["tag📌:发布新标签"] = append(classified["tag📌:发布新标签"], formattedCommit)
			case "config":
				classified["config🔧:修改配置文件"] = append(classified["config🔧:修改配置文件"], formattedCommit)
			case "git":
				classified["git🙈:添加或修改.gitignore文件"] = append(classified["git🙈:添加或修改.gitignore文件"], formattedCommit)
			default:
				classified["others:其他"] = append(classified["others:其他"], formattedCommit)
			}
		} else {
			classified["others:其他"] = append(classified["others:其他"], commit)
		}
	}

	return classified
}

func Run(cmd *cobra.Command, args []string) {
	// 数据初始化
	if err := base.InitPath(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	tags, err := getTags()
	if err != nil {
		fmt.Println("获取 Tag:", color.RedString(err.Error()))
		return
	}

	repoURL, err := getRepoURL()
	if err != nil {
		fmt.Println("获取地址:", color.RedString(err.Error()))
		return
	}
	var changelogBuilder strings.Builder
	changelogBuilder.WriteString("# 更新日志\n\n")

	for i := 0; i < len(tags); i++ {
		tag1 := tags[i]
		date, err := getTagDate(tag1)
		if err != nil {
			fmt.Println("获取日期:", color.RedString(err.Error()))
			return
		}
		var tag2 string
		if i+1 < len(tags) {
			tag2 = tags[i+1]
		}
		if tag2 == "" {
			changelogBuilder.WriteString(fmt.Sprintf("## %s (%s)\n", tag1, date))
		} else {
			changelogBuilder.WriteString(fmt.Sprintf("## [%s](%s/compare/%s...%s) (%s)\n", tag1, repoURL, tag2, tag1, date))
		}
		commits, err := getCommits(tag1, tag2)
		if err != nil {
			log.Fatal(err)
		}
		categories := []string{
			"feat✨:引入新功能",
			"fix🐛:修复 bug",
			"style💄:更新UI样式文件",
			"format🥚:格式化代码",
			"docs📝:添加/更新文档",
			"perf👌:提高性能/优化",
			"test✅:增加测试代码",
			"refactor🎨:改进代码结构/代码格式",
			"patch🚑:添加重要补丁",
			"file📦:添加新文件",
			"publish🚀:发布新版本",
			"tag📌:发布新标签",
			"config🔧:修改配置文件",
			"git🙈:添加或修改.gitignore文件",
			"others:其他",
			"init🎉:初次提交/初始化项目",
		}

		if len(commits) == 0 {
			changelogBuilder.WriteString("无更新.\n")
		} else {
			classifiedCommits := classifyCommits(commits, repoURL)
			for _, category := range categories {
				if commits, exists := classifiedCommits[category]; exists && len(commits) > 0 {
					changelogBuilder.WriteString(fmt.Sprintf("### %s\n", category))
					for _, commit := range commits {
						changelogBuilder.WriteString(fmt.Sprintf("- %s\n", commit))
					}
					changelogBuilder.WriteString("\n")
				}
			}
		}
		changelogBuilder.WriteString("\n")
	}
	outApiFile := path.Join(cast.ToString(base.Path), "CHANGELOG.md")
	file, err := os.OpenFile(outApiFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}
	_, err = file.WriteString(changelogBuilder.String())
	if err != nil {
		fmt.Println("写入文件失败:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(outApiFile))
}
