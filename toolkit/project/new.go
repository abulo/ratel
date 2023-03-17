package project

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/fatih/color"
)

// Project is a project template.
type Project struct {
	Name string
	Path string
}

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s 已经存在\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 您想要覆盖文件夹吗 ?",
			Help:    "删除现有文件夹并创建项目.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 正在创建服务 %s, 代码仓库是 %s, 请稍候.\n\n", p.Name, layout)
	repo := base.NewRepo(layout, branch)
	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}
	base.Tree(to, dir)

	fmt.Printf("\n🍺 项目创建成功 %s\n", color.GreenString(p.Name))
	fmt.Print("💻 使用以下命令进入项目 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println("			🤝 感谢使用 Ratel")
	return nil
}

// New new a project from remote repo.
func (p *Project) NewFront(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s 已经存在\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 您想要覆盖文件夹吗 ?",
			Help:    "删除现有文件夹并创建项目.",
		}
		e := survey.AskOne(prompt, &override)
		if e != nil {
			return e
		}
		if !override {
			return err
		}
		os.RemoveAll(to)
	}
	fmt.Printf("🚀 正在创建服务 %s, 代码仓库是 %s, 请稍候.\n\n", p.Name, layout)
	repo := base.NewRepo(layout, branch)
	if err := repo.CopyFrontTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}
	base.Tree(to, dir)

	fmt.Printf("\n🍺 项目创建成功 %s\n", color.GreenString(p.Name))
	fmt.Print("💻 使用以下命令进入项目 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println("			🤝 感谢使用 Ratel")
	return nil
}
