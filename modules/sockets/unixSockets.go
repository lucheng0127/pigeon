package sockets

import (
	"net"
	"os"
	"pigeon/modules/tasks"
	"pigeon/pigeond/log"
)

// UnixSocket define the location of socket file
type UnixSocket struct {
	SocketFile string
}

// Launch a unix socket server
func Launch(s sockets, done chan bool) {
	s.listen()
	done <- true
}

// Send msg to pigeond
func Send(s sockets, msg string) string {
	rst := s.send(msg)
	return rst
}

type sockets interface {
	listen()
	send(string) string
}

func checkError(err error, exit bool) {
	if err != nil {
		log.Log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

func handleUnixConn(conn *net.UnixConn, msg chan string) {
	// Get msg from conn and send to msg
	buf := make([]byte, 64)
	var dataString string

	for {
		readLen, _, err := conn.ReadFromUnix(buf)
		if err != nil {
			checkError(err, false)
			break
		}

		if readLen == 0 {
			log.Log.Debug("All data received")
			break // All data received
		} else {
			dataString += string(buf[:readLen])
			if readLen < 64 {
				log.Log.Debug("All data received")
				break // All data received
			}
		}

		buf = make([]byte, 64)
	}
	msg <- dataString
}

func (us *UnixSocket) listen() {
	// Start unix socket server
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
		msg := make(chan string) // Use msg channel transfer socket data
		conn, err := l.AcceptUnix()
		checkError(err, false)
		log.Log.Info("Get new conn")

		// Get msg from conn and send to msg
		go handleUnixConn(conn, msg)
		taskRst := make(chan string) // Use taskRst channel to get task result
		select {
		case taskInfo := <-msg:
			log.Log.Infof("Get %s from conn", taskInfo)
			go tasks.TaskManage(taskInfo, taskRst)
		}
		select {
		case ts := <-taskRst: // If taskRst send back to socket connection
			log.Log.Infof("Send %s to conn", ts)
			_, err = conn.Write([]byte(ts))
			checkError(err, false)
		}
	}

}

func handClientUnixConn(conn *net.UnixConn, rst chan string) {
	var msg string
	buf := make([]byte, 64)
	for {
		readLen, _, err := conn.ReadFromUnix(buf)
		if err != nil {
			checkError(err, false)
			break
		}

		if readLen == 0 {
			log.Log.Debug("All data received")
			break // All data received
		} else {
			msg += string(buf[:readLen])
			if readLen < 64 {
				log.Log.Debug("All data received")
				break
			}
		}

		buf = make([]byte, 64)
	}
	log.Log.Debug("Client get result", msg)
	rst <- msg
}

func (us *UnixSocket) send(msg string) string {
	addr, err := net.ResolveUnixAddr("unix", us.SocketFile)
	checkError(err, true)

	conn, err := net.DialUnix("unix", nil, addr)
	checkError(err, true)
	defer conn.Close()
	log.Log.Debug("Client connect to", conn)

	_, err = conn.Write([]byte(msg))
	checkError(err, true)
	log.Log.Debugf("Client send %s to conn", msg)

	rst := make(chan string)
	go handClientUnixConn(conn, rst)

	select {
	case rstStr := <-rst:
		return rstStr
	}
}
