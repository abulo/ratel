package config

import (
	"fmt"

	"github.com/abulo/ratel/v2/util"
	"github.com/fsnotify/fsnotify"
)

type watcher = fsnotify.Watcher

func newWatcher() (*watcher, error) {
	return fsnotify.NewWatcher()
}

func WatchConfig(suffix string) {
	dc.WatchConfig(suffix)
}

func OnConfigChange(run func(in fsnotify.Event)) { dc.OnConfigChange(run) }
func (v *Config) OnConfigChange(run func(in fsnotify.Event)) {
	v.onConfigChange = run
}

func (c *Config) ConfigDir() []string {
	dir := make([]string, 0)
	for _, v := range c.loadedFiles {
		dirname := util.GetParentDirectory(v)
		if !util.InArray(dirname, dir) {
			dir = append(dir, util.GetParentDirectory(v))
		}
	}
	return dir
}

// loadDir
func (c *Config) WatchConfig(suffix string) {
	watcher, err := newWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	done := make(chan bool)
	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				//do something
				if ev.Op&fsnotify.Create == fsnotify.Create {
					if c.onConfigChange != nil {
						c.onConfigChange(ev)
					}
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if c.onConfigChange != nil {
						c.onConfigChange(ev)
					}
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if c.onConfigChange != nil {
						c.onConfigChange(ev)
					}
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if c.onConfigChange != nil {
						c.onConfigChange(ev)
					}
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if c.onConfigChange != nil {
						c.onConfigChange(ev)
					}
				}
			case err := <-watcher.Errors:
				fmt.Println(err)
			}
		}
	}()
	dir := c.ConfigDir()
	for _, v := range dir {
		fmt.Println(v)
		err = watcher.Add(v)
		if err != nil {
			fmt.Println(err)
		}
	}
	<-done
	watcher.Close()
}
