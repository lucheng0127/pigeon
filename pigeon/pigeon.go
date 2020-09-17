package main

import (
	"fmt"
	"pigeon/modules/sockets"
)

func main() {
	us := sockets.UnixSocket{SocketFile: "/var/run/pigeond.socket"}
	rst := sockets.Send(&us, "F ADD_SCRIPT /tmp/test_script.tar END")
	fmt.Println("Add script", rst)
	rst = sockets.Send(&us, "F LIST_SCRIPTS END")
	fmt.Println("List scripts", rst)
}
