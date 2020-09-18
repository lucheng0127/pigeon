package main

import (
	"os"
	"pigeon/modules/sockets"
)

func main() {
	done := make(chan bool)
	if _, err := os.Stat("/var/run/pigeon"); os.IsNotExist(err) {
		os.Mkdir("/var/run/pigeon", 0755)
	}
	us := sockets.UnixSocket{SocketFile: "/var/run/pigeon/pigeond.socket"}
	sockets.Launch(&us, done)
	<-done
}
