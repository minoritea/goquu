package main

import (
	"./job"
)
func main() {
	server, _ := job.NewServer()
	server.Run()
}
