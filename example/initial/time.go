package initial

import "time"

// InitLaunchTime set binary start time
func (initial *Initial) InitLaunchTime(launchTime time.Time) *Initial {
	initial.LaunchTime = launchTime
	return initial
}
