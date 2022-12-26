	rl := ratelimit.NewLeakyBucket(time.Second, 15) // per second
    rl.TakeAvailable()

    rl = ratelimit.NewTokenBucket(time.Microsecond, 15) // per Microsecond
    rl.TakeAvailable()