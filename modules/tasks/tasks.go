package tasks

import (
	"pigeon/pigeond/log"
)

// Task struct
// Name -> task name
// Args -> task args
// AutoAck -> if auto ack, no matter task succeed
// or failed, it will return 0 to task result
// channel and send back to socket conn, but ExitCode
// will recode the result of task, if failed it will
// recode 1, but it will return 0 if auto ack
// ExitCode -> 1 task failed, 0 task succeed
type Task struct {
	Name     string
	Args     []string
	AutoAck  bool
	ExitCode int
}

// TaskManage use to run task with msg and return result to taskRst
func TaskManage(msg string, taskRst chan string) {
	log.Log.Info(msg)
	taskRst <- "0 Add script succeed"
}
