package initial

import "github.com/abulo/ratel/v3/stores/session"

// Session load session && returns an session instance.
func (initial *Initial) Session(name string) *session.Session {
	return &session.Session{
		Name:   name,
		Driver: initial.Store.LoadRedis("redis"),
		TTL:    initial.Config.Int64("cookie.expires", 300),
	}
}
