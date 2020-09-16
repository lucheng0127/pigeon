package main

import (
	"pigeon/modules/sockets"
)

func main() {
	done := make(chan bool)
	us := sockets.UnixSocket{SocketFile: "/var/run/pigeond.socket"}
	sockets.Launch(&us, done)
	<-done
}
