package fx

import "github.com/abulo/ratel/v2/thread"

// Parallel runs fns parallelly and waits for done.
func Parallel(fns ...func()) {
	group := thread.NewRoutineGroup()
	for _, fn := range fns {
		group.RunSafe(fn)
	}
	group.Wait()
}
