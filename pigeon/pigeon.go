package main

import (
	"pigeon/pigeon/cmd"
)

const pigeondSocketFile = "/var/run/pigeond.socket"

func main() {
	cmd.Execute()
}
