package simulator

import (
	"os"
	"strconv"
)

type MessageTask struct {
	messageType string
	from        *Node
	to          *Node
	block       *Block
	interval    int64
}

func BlockMessageTask(from *Node, to *Node, block *Block, interval int64) {
	task := &MessageTask{"BlockMessageTask", from, to, block, interval}
	putMessageTask(task)
}

func InvMessageTask(from *Node, to *Node, block *Block, interval int64) {
	task := &MessageTask{"InvMessageTask", from, to, block, interval}
	putMessageTask(task)
}

func RecMessageTask(from *Node, to *Node, block *Block, interval int64) {
	task := &MessageTask{"RecMessageTask", from, to, block, interval}
	putMessageTask(task)
}

func (task *MessageTask) Run() {
	if task.messageType == "BlockMessageTask" {
		task.from.SendNextBlockMessage()

		f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err2 := f.WriteString("{" + "\"kind\":\"flow-block\"," + "\"content\":{" + "\"transmission-timestamp\":" + strconv.FormatInt(GetCurrentTime()-task.interval, 10) + "," + "\"reception-timestamp\":" + strconv.FormatInt(GetCurrentTime(), 10) + "," + "\"begin-node-id\":" + strconv.Itoa(task.from.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(task.to.nodeID) + "," + "\"block-id\":" + strconv.Itoa(task.block.id) + "}" + "},")

		if err2 != nil {
			panic(err2)
		}
	}
	task.to.ReceiveMessage(task)
}
