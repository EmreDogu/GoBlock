package simulator

import (
	"container/heap"
	"reflect"
)

type ScheduledTask struct {
	taskType      string
	mintingTask   *MintingTask
	messageTask   *MessageTask
	scheduledTime int64
	index         int
}
type PriorityQueue []*ScheduledTask

var TaskMap = make(map[any]*ScheduledTask)
var Pq PriorityQueue = make(PriorityQueue, 0)
var CurrentTime int64 = 0

func GetPriorityQueue() *PriorityQueue {
	return &Pq
}

func GetTask() any {
	if len(Pq) > 0 {
		if Pq[0].taskType == "mintingTask" {
			return Pq[0].mintingTask
		} else {
			return Pq[0].messageTask
		}
	} else {
		return nil
	}
}

func GetMintingTask() *MintingTask {
	return Pq[0].mintingTask
}

func GetScheduledTask() *ScheduledTask {
	return Pq[0]
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].scheduledTime < pq[j].scheduledTime
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*ScheduledTask)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[0]
	old[0] = nil    // avoid memory leak
	item.index = -1 // for safety
	*pq = old[1:n]
	return item
}

func (pq *PriorityQueue) Peek() *ScheduledTask {
	old := *pq
	item := old[0]
	return item
}

func (pq *PriorityQueue) updateMintingTask(item *ScheduledTask, task *MintingTask, scheduledTime int64) {
	item.mintingTask = task
	item.scheduledTime = scheduledTime
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) updateMessageTask(item *ScheduledTask, task *MessageTask, scheduledTime int64) {
	item.messageTask = task
	item.scheduledTime = scheduledTime
	heap.Fix(pq, item.index)
}

func GetCurrentTime() int64 {
	return CurrentTime
}

func putMintingTask(task *MintingTask) {
	ScheduledTask := &ScheduledTask{"mintingTask", task, nil, GetCurrentTime() + task.interval, 0}
	TaskMap[task] = ScheduledTask
	heap.Push(&Pq, ScheduledTask)
	Pq.updateMintingTask(ScheduledTask, ScheduledTask.mintingTask, ScheduledTask.scheduledTime)
}

func putMessageTask(task *MessageTask) {
	ScheduledTask := &ScheduledTask{"messageTask", nil, task, GetCurrentTime() + task.interval, 0}
	TaskMap[task] = ScheduledTask
	heap.Push(&Pq, ScheduledTask)
	Pq.updateMessageTask(ScheduledTask, ScheduledTask.messageTask, ScheduledTask.scheduledTime)
}

func removeTask(this *MintingTask) {
	_, ok := TaskMap[this]
	if ok {
		for i := 0; i < len(Pq); i++ {
			n := Pq[i]
			if n.mintingTask == this {
				Pq[i], Pq[len(Pq)-1] = Pq[len(Pq)-1], Pq[i]
				Pq = Pq[:len(Pq)-1]
				i--
			}
		}

		delete(TaskMap, this)
	}
}

func RunTask() {
	if Pq.Len() > 0 {
		currentTask := GetScheduledTask()
		heap.Pop(&Pq)
		CurrentTime = currentTask.scheduledTime
		delete(TaskMap, currentTask)
		if reflect.ValueOf(currentTask.messageTask).IsNil() {
			currentTask.mintingTask.Run()
		} else {
			currentTask.messageTask.Run()
		}
	}
}
