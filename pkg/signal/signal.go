package signal

import (
	"os"
	"os/signal"
	"syscall"
)

func WaitOSSignal() os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	return s
}
