package main

import (
	"./goquu"
	"./queue"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func trap() {
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		<-sig
		fmt.Println("server is shutting down...")
		queue.CloseAll()
		os.Exit(0)
	}()
}
func main() {
	trap()
	server, _ := goquu.NewServer()
	server.Run()
}
