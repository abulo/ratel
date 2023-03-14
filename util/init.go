package util

import (
	"math/rand"
	"time"
)

func Do(attempts int, sleep time.Duration, f func() error) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if err := f(); err != nil {
		if attempts--; attempts >= 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(r.Int63n(int64(sleep)))
			sleep = sleep + jitter/2
			time.Sleep(sleep)
			return Do(attempts, 2*sleep, f)
		}
		return err
	}
	return nil
}
