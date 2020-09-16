package sockets

import (
	"net"
	"os"
	"pigeon/pigeond/log"
)

// UnixSocket define the location of socket file
type UnixSocket struct {
	SocketFile string
}

// Launch a unix socket server
func Launch(s sockets, done chan bool) {
	log.Log.Info("Start pigeond server")
	s.listen()
	done <- true
	log.Log.Info("Pigeond server closed")
}

// Send msg to pigeond
func Send(s sockets, msg string) {
	s.send(msg)
}

type sockets interface {
	listen()
	send(string)
}

func checkError(err error, exit bool) {
	if err != nil {
		log.Log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

func chandleUnixConn(conn *net.UnixConn, msg chan string) {
	buf := make([]byte, 64)
	var dataString string

	for {
		readLen, _, err := conn.ReadFromUnix(buf)

		if err != nil {
			checkError(err, false)
			break
		}

		if readLen == 0 {
			break // All data received
		} else {
			dataString += string(buf[:readLen])
		}

		buf = make([]byte, 64)
	}
	msg <- dataString
}

func (us *UnixSocket) listen() {
	if _, err := os.Stat(us.SocketFile); err == nil {
		err := os.Remove(us.SocketFile)
		checkError(err, true)
	}

	addr, err := net.ResolveUnixAddr("unix", us.SocketFile)
	checkError(err, true)
	l, err := net.ListenUnix("unix", addr)
	checkError(err, true)
	log.Log.Info("Start to listen", addr)

	for {
		// This's only one conn each time
		msg := make(chan string)
		conn, err := l.AcceptUnix()
		checkError(err, false)
		go chandleUnixConn(conn, msg)
		select {
		case taskInfo := <-msg:
			log.Log.Info("Get msg:", taskInfo)
		}
	}

}

func (us *UnixSocket) send(msg string) {
	addr, err := net.ResolveUnixAddr("unix", us.SocketFile)
	checkError(err, true)

	conn, err := net.DialUnix("unix", nil, addr)
	checkError(err, true)
	log.Log.Info("Connect to pigeond", conn)

	_, err = conn.Write([]byte(msg))
	checkError(err, true)
	defer conn.Close()
}
