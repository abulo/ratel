# crontab


功能完善的协程班计划任务执行库。

```go
package main

import "github.com/abulo/ratel/crontab"

func something() {
	//do

	return
}
func main() {
	sched := crontab.NewScheduler()
	sched.Schedule().Every(10).Seconds().Do(something)
	sched.Schedule().Every(3).Minutes().Do(something)
	sched.Schedule().Every(4).Hours().Do(something)
	sched.Schedule().Every(2).Days().At("12:32").Do(something)
	sched.Schedule().Every(12).Weeks().Do(something)
	sched.Schedule().Every(1).Monday().Do(something)
	sched.Schedule().Every(1).Saturday().At("8:00").Do(something)
	sched.Run()
}
```