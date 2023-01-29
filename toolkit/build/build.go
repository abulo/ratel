package build

import (
	"context"
	"fmt"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/abulo/ratel/v3/watch"
	"github.com/fatih/color"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	CmdNew = &cobra.Command{
		Use:   "build",
		Short: "热编译",
		Long:  "热编译应用程序: toolkit build",
		Run:   Run,
	}
)

type buildOption struct {
	AutoTidy         bool     `toml:"AutoTidy"`
	MainFiles        string   `toml:"MainFiles"`
	OutputName       string   `toml:"OutputName"`
	Exts             []string `toml:"Exts"`
	Excludes         []string `toml:"Excludes"`
	AppArgs          string   `toml:"AppArgs"`
	Recursive        bool     `toml:"Recursive"`
	Dirs             []string `toml:"Dirs"`
	WatcherFrequency string   `toml:"WatcherFrequency"`
	Asm              string   `toml:"Asm"`
	Gccgo            string   `toml:"Gccgo"`
	Gc               string   `toml:"Gc"`
	Ld               string   `toml:"Ld"`
}

func Run(cmd *cobra.Command, args []string) {
	if err := base.InitPath(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	if err := base.InitConfig(); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	// 初始化
	optionInfo := buildOption{}
	if err := base.Config.BindStruct("watch", &optionInfo); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	// 获取配置文件
	option := &watch.Options{}
	option.AutoTidy = optionInfo.AutoTidy
	option.MainFiles = optionInfo.MainFiles
	option.OutputName = optionInfo.OutputName
	option.Exts = optionInfo.Exts
	option.Excludes = optionInfo.Excludes
	option.AppArgs = optionInfo.AppArgs
	option.Recursive = optionInfo.Recursive
	option.Dirs = optionInfo.Dirs
	option.WatcherFrequency = util.Duration(optionInfo.WatcherFrequency)
	optionFlags := watch.Flags{}
	if !util.Empty(cast.ToString(optionInfo.Gc)) {
		optionFlags["gc"] = cast.ToString(optionInfo.Gc)
	}
	if !util.Empty(cast.ToString(optionInfo.Asm)) {
		optionFlags["asm"] = cast.ToString(optionInfo.Asm)
	}
	if !util.Empty(cast.ToString(optionInfo.Gccgo)) {
		optionFlags["gccgo"] = cast.ToString(optionInfo.Gccgo)
	}
	if !util.Empty(cast.ToString(optionInfo.Ld)) {
		optionFlags["ld"] = cast.ToString(optionInfo.Ld)
	}
	option.Flags = optionFlags
	ctx := context.Background()

	if err := watch.Watch(ctx, option); err != nil {
		fmt.Println("失败:", color.RedString(err.Error()))
		return
	}
}
