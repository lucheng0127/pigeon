package main

import (
	"pigeon/modules/sockets"
)

func main() {
	us := sockets.UnixSocket{SocketFile: "/var/run/pigeond.socket"}
	sockets.Send(&us, "test msg")
}
