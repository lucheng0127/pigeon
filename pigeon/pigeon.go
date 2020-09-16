package main

import (
	"fmt"
	"pigeon/modules/sockets"
)

func main() {
	us := sockets.UnixSocket{SocketFile: "/var/run/pigeond.socket"}
	rst := sockets.Send(&us, "F UPLOAD_SCRIPT /tmp/test_script.tar END")
	fmt.Println("Task result", rst)
}
