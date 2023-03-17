package frontend

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/project"
	"github.com/spf13/cobra"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "frontend",
	Short: "前端项目",
	Long:  "前端项目: toolkit frontend helloworld",
	Run:   run,
}

var (
	repoURL string
	branch  string
	timeout string
)

func init() {
	repoURL = "https://github.com/abulo/layout-vue.git"
	timeout = "60s"
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
	CmdNew.Flags().StringVarP(&timeout, "timeout", "t", timeout, "time out")
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
	p := &project.Project{Name: path.Base(name), Path: name}
	done := make(chan error, 1)
	go func() {
		done <- p.NewFront(ctx, wd, repoURL, branch)
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
