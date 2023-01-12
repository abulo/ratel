package build

import (
	"context"
	"fmt"
	"reflect"
	"time"

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
	MainFiles        string   `toml:"MainFiles"`
	OutputName       string   `toml:"OutputName"`
	Exts             []string `toml:"Exts"`
	Excludes         []string `toml:"Excludes"`
	AppArgs          string   `toml:"AppArgs"`
	Recursive        bool     `toml:"Recursive"`
	Dirs             []string `toml:"Dirs"`
	WatcherFrequency string   `toml:"WatcherFrequency"`
	Flags            struct {
		Asm   string `toml:"Asm"`
		Gccgo string `toml:"Gccgo"`
		Gc    string `toml:"Gc"`
		Ld    string `toml:"Ld"`
	} `toml:"Flags"`
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
	// ctx := context.Background()

	optionList := []buildOption{}
	if err := base.Config.BindStruct("watch", &optionList); err != nil {
		fmt.Println("初始化:", color.RedString(err.Error()))
		return
	}

	watchOptionList := make([]*watch.Options, 0)
	for _, optionConfig := range optionList {
		// 获取配置文件
		option := &watch.Options{}
		option.AutoTidy = optionConfig.AutoTidy
		option.MainFiles = optionConfig.MainFiles
		option.OutputName = optionConfig.OutputName
		option.Exts = optionConfig.Exts
		option.Excludes = optionConfig.Excludes
		option.AppArgs = optionConfig.AppArgs
		option.Recursive = optionConfig.Recursive
		option.Dirs = optionConfig.Dirs
		option.WatcherFrequency = 10 * time.Second
		flags := optionConfig.Flags
		optionFlags := watch.Flags{}
		sVal := reflect.ValueOf(flags)
		sType := reflect.TypeOf(flags)
		if sType.Kind() == reflect.Ptr {
			//用Elem()获得实际的value
			sVal = sVal.Elem()
			sType = sType.Elem()
		}
		num := sVal.NumField()
		for i := 0; i < num; i++ {
			f := sType.Field(i)
			val := sVal.Field(i).Interface()
			if !util.Empty(cast.ToString(val)) {
				optionFlags[util.StrToLower(f.Name)] = cast.ToString(val)
			}
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
