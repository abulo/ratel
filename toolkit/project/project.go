package project

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v2/toolkit/base"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "项目创建",
	Long:  "新建微服务项目: toolkit new helloworld",
	Run:   run,
}

var (
	repoURL string
	branch  string
	timeout string
	nomod   bool
)

func init() {
	if repoURL = os.Getenv("RATEL_LAYOUT_REPO"); repoURL == "" {
		repoURL = "https://github.com/abulo/layout.git"
	}
	timeout = "60s"
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdNew.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
	CmdNew.Flags().BoolVarP(&nomod, "nomod", "", nomod, "retain go mod")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "项目名称",
			Help:    "项目命名:字母小写",
		}
		err = survey.AskOne(prompt, &name)
		if err != nil || name == "" {
			return
		}
	} else {
		name = args[0]
	}
	p := &Project{Name: path.Base(name), Path: name}
	done := make(chan error, 1)
	go func() {
		if !nomod {
			done <- p.New(ctx, wd, repoURL, branch)
			return
		}
		if _, e := os.Stat(path.Join(wd, "go.mod")); os.IsNotExist(e) {
			done <- fmt.Errorf("🚫 未在 %s 中找到 go.mod 文件", wd)
			return
		}

		mod, e := base.ModulePath(path.Join(wd, "go.mod"))
		if e != nil {
			panic(e)
		}
		done <- p.Add(ctx, wd, repoURL, branch, mod)
	}()
	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			fmt.Fprint(os.Stderr, "\033[31mERROR: 项目创建超时 \033[m\n")
			return
		}
		fmt.Fprintf(os.Stderr, "\033[31mERROR: 项目创建失败(%s)\033[m\n", ctx.Err().Error())
	case err = <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\033[31mERROR: 项目创建失败(%s)\033[m\n", err.Error())
		}
	}
}
