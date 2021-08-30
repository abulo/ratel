package signals

import (
	"os"
	"os/signal"
)

//Shutdown suport twice signal must exit
func Listen(ln func(grace os.Signal)) {
	sig := make(chan os.Signal, 3)
	signal.Notify(
		sig,
		shutdownSignals...,
	)
	go func() {
		s := <-sig
		go ln(s)
		<-sig
		// os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}
