package simulator

import (
	"container/heap"
	"reflect"
	"time"
)

type ScheduledTask struct {
	taskType      string
	mintingTask   *MintingTask
	messageTask   *MessageTask
	scheduledTime float64
	index         int
}
type PriorityQueue []*ScheduledTask

var TaskMap = make(map[any]*ScheduledTask)
var Pq PriorityQueue = make(PriorityQueue, 0)
var CurrentTime int = 0

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].scheduledTime < pq[j].scheduledTime
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
	if reflect.ValueOf(pq[i].messageTask).IsNil() {
		pq[i].mintingTask.index = i
		pq[j].mintingTask.index = j
	}
}

func (pq PriorityQueue) GetHeight(i int) int {
	return pq[i].mintingTask.minter.block.height
}

func (pq PriorityQueue) GetMintingTask(i int) *MintingTask {
	return pq[i].mintingTask
}

func (pq PriorityQueue) GetMessageTask(i int) *MessageTask {
	return pq[i].messageTask
}

func (pq PriorityQueue) GetScheduledTime(i int) int {
	return int(pq[i].scheduledTime)
}

func (pq PriorityQueue) GetTaskType(i int) string {
	return pq[i].taskType
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

func (pq *PriorityQueue) updateMintingTask(item *ScheduledTask, task *MintingTask, scheduledTime float64) {
	item.mintingTask = task
	item.scheduledTime = scheduledTime
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) updateMessageTask(item *ScheduledTask, task *MessageTask, scheduledTime float64) {
	item.messageTask = task
	item.scheduledTime = scheduledTime
	heap.Fix(pq, item.index)
}

func GetCurrentTime() int64 {
	return time.Now().UnixMilli()
}

func putMintingTask(task *MintingTask) {
	ScheduledTask := &ScheduledTask{"mintingTask", task, nil, float64(GetCurrentTime()) + task.interval, 0}
	TaskMap[task] = ScheduledTask
	heap.Push(&Pq, ScheduledTask)
	Pq.updateMintingTask(ScheduledTask, ScheduledTask.mintingTask, ScheduledTask.scheduledTime)
}

func putMessageTask(task *MessageTask) {
	ScheduledTask := &ScheduledTask{"messageTask", nil, task, float64(GetCurrentTime()) + task.interval, 0}
	TaskMap[task] = ScheduledTask
	heap.Push(&Pq, ScheduledTask)
	Pq.updateMessageTask(ScheduledTask, ScheduledTask.messageTask, ScheduledTask.scheduledTime)
}

func removeTask(this *MintingTask) {

	for i := 0; i < len(Pq); i++ {
		n := Pq[i]
		if n.mintingTask == this {
			Pq[i], Pq[len(Pq)-1] = Pq[len(Pq)-1], Pq[i]
			Pq = Pq[:len(Pq)-1]
			i--
		}
	}

	heap.Init(&Pq)
	delete(TaskMap, this)
}
