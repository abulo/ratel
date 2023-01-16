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
	"golang.org/x/sync/errgroup"
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
	MainFiles        []string `toml:"MainFiles"`
	OutputName       []string `toml:"OutputName"`
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

	if len(optionInfo.MainFiles) < 1 {
		fmt.Println("初始化:", color.RedString("MainFiles"))
		return
	}

	if len(optionInfo.OutputName) < 1 {
		fmt.Println("初始化:", color.RedString("OutputName"))
		return
	}

	if len(optionInfo.MainFiles) != len(optionInfo.OutputName) {
		fmt.Println("初始化:", color.RedString("MainFiles&OutputName->len"))
		return
	}

	builderNumber := len(optionInfo.MainFiles)
	watchOptionList := make([]*watch.Options, 0)
	for i := 0; i < builderNumber; i++ {
		// 获取配置文件
		option := &watch.Options{}
		option.AutoTidy = optionInfo.AutoTidy
		option.MainFiles = optionInfo.MainFiles[i]
		option.OutputName = optionInfo.OutputName[i]
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
		watchOptionList = append(watchOptionList, option)
	}
	ctx := context.Background()
	var eg errgroup.Group
	for _, s := range watchOptionList {
		s := s
		eg.Go(func() (err error) {
			err = watch.Watch(ctx, s)
			return
		})
	}
	eg.Wait()
}
