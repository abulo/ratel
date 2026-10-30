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

// Project 描述待创建项目的名称与模块路径，仅在脚手架命令执行期间使用，不保证并发安全。
//
// 零值字段均为空，不可直接用于创建，须先填充 Name 与 Path。
type Project struct {
	// Name 为项目目录名（路径末级）。
	Name string
	// Path 为项目在模块内的相对路径。
	Path string
}

// New 从 layout 远程仓库克隆后端项目模板到 dir/Name；目标目录已存在时先弹出确认，
// 用户同意覆盖后删除原目录再复制；ctx 用于控制整个拉取过程的超时。
//
// 仓库复制失败时返回对应错误；用户拒绝覆盖或确认流程被中断时返回错误。
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

// NewFront 从 layout 远程仓库克隆前端项目模板，目录覆盖确认与超时约定同 New；区别在于
// 复制的是前端模板内容。
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
