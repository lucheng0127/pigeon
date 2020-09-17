package tasks

import (
	"fmt"
	"pigeon/pigeond/log"
	"strings"
	"sync"
	"time"
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
	Result   string
}

func taskProxy(task *Task, wg *sync.WaitGroup) {
	defer wg.Done()

	switch task.Name {
	case "LIST_SCRIPTS":
		time.Sleep(3 * time.Second)
		task.ExitCode, task.Result = ListScript()
	case "ADD_SCRIPT":
		time.Sleep(3 * time.Second)
		task.ExitCode = 0
		task.Result = "Add script succeed"
	default:
		log.Log.Infof("No task for command %s", task.Name)
		task.ExitCode = 1
		task.Result = "Error command"
	}
}

// TaskManage use to run task with msg and return result to taskRst
func TaskManage(msg string, taskRst chan string) {
	// msg format: [T/F] [Task Name] [Arg1] [Arg2] ... [END]
	log.Log.Info(msg)
	task := new(Task)
	task.ExitCode = 0
	task.Result = ""
	task.AutoAck = false
	msgList := strings.Split(msg, " ")
	finished := make(chan bool)

	// Block until task finish
	select {
	case <-finished:
		log.Log.Infof("Task %s finished", task.Name)
	default:
		if len(msgList) < 3 {
			taskRst <- "1 Command missing"
			finished <- true
		}
		if msgList[len(msgList)-1] != "END" {
			taskRst <- "1 Command missing"
			finished <- true
		}
		if msgList[0] == "T" {
			task.AutoAck = true
		} else {
			task.AutoAck = false
		}
		task.Name = msgList[1]
		task.Args = msgList[2 : len(msgList)-1]

		var wg sync.WaitGroup
		go taskProxy(task, &wg)
		wg.Add(1)
		if task.AutoAck == false {
			wg.Wait()
		}

		rst := fmt.Sprintf("%d %s", task.ExitCode, task.Result)
		taskRst <- rst
		finished <- true
	}
}
