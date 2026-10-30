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

// repoAddIgnores 为追加服务时从模板目录中排除的版本控制与忽略规则文件。
var repoAddIgnores = []string{
	".git", ".github", ".gitignore",
}

// Add 在现有项目 dir 下追加一个服务目录，从 layout 仓库复制模板到 dir/Path，并将模板中
// 的模块路径替换为 mod/Path；目标目录已存在时先弹出确认，用户同意覆盖后删除原目录再复制，
// ctx 用于控制整个拉取过程的超时。
//
// 仓库复制失败时返回对应错误；用户拒绝覆盖或确认流程被中断时返回错误。
func (p *Project) Add(ctx context.Context, dir string, layout string, branch string, mod string) error {
	to := path.Join(dir, p.Path)

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

	fmt.Printf("🚀 添加服务 %s, 代码仓库是 %s, 请稍候.\n\n", p.Name, layout)

	repo := base.NewRepo(layout, branch)

	if err := repo.CopyToV2(ctx, to, path.Join(mod, p.Path), repoAddIgnores, []string{path.Join(p.Path, "api"), "api"}); err != nil {
		return err
	}
	base.Tree(to, dir)
	fmt.Printf("\n🍺 服务添加成功 %s\n", color.GreenString(p.Name))
	fmt.Print("💻 使用以下命令进入项目 👇:\n\n")
	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println("			🤝 感谢使用 Ratel")
	return nil
}
