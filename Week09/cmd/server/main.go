package main

import (
	"Week09/internal"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := internal.NewServer()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			sig := <- sigCh
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				s.GracefulStop()
			default:
			}
		}
	}()

	s.Start()
}