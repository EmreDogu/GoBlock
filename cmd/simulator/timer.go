package simulator

import (
	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/network/p2p"
	"github.com/rdleal/go-priorityq/kpq"
)

type ScheduledTask struct {
	task          *p2p.Task
	scheduledTime int64
}

var taskQueue = kpq.NewKeyedPriorityQueue[*ScheduledTask](func(a, b int64) bool {
	return a < b
})
var taskMap = make(map[*p2p.Task]*ScheduledTask)

func (s *Simulator) RemoveTask(task *p2p.Task) {
	if value, ok := taskMap[task]; ok {
		taskQueue.Remove(value)
		delete(taskMap, task)
	}
}

func (s *Simulator) PutTask(newTask *p2p.Task) {
	scheduledTask := &ScheduledTask{task: newTask, scheduledTime: configs.GetCurrentTime() + newTask.Message.GetInterval()}
	taskMap[newTask] = scheduledTask
	taskQueue.Push(scheduledTask, configs.GetCurrentTime()+newTask.Message.GetInterval())
}

func GetTask() *ScheduledTask {
	key, ok := taskQueue.PeekKey()
	if ok {
		return key
	} else {
		return nil
	}
}

func RunTask() {
	currentScheduledTask, _, ok := taskQueue.Pop()
	if ok {
		currentTask := currentScheduledTask.task
		configs.SetCurrentTime(currentScheduledTask.scheduledTime)
		delete(taskMap, currentTask)
		currentTask.Run()
	}
}
