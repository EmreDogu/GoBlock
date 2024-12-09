package p2p

import (
	"os"
	"strconv"

	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
)

type Messages struct {
	messageType string
	from        Node
	to          Node
	block       *block.Block
	interval    int64
}

func CreateInvTask(fromNode Node, toNode Node, sentBlock *block.Block) *Task {
	messages := &Messages{messageType: "inv", from: fromNode, to: toNode, block: sentBlock}
	return &Task{Message: messages}
}

func CreateCmpctBlockTask(fromNode Node, toNode Node, sentBlock *block.Block, delay int64) *Task {
	messages := &Messages{messageType: "cmpctblock", from: fromNode, to: toNode, block: sentBlock, interval: configs.GetLatency(fromNode.GetRegion(), toNode.GetRegion()) + delay}
	return &Task{Message: messages}
}

func CreateBlockTask(fromNode Node, toNode Node, sentBlock *block.Block, delay int64) *Task {
	messages := &Messages{messageType: "block", from: fromNode, to: toNode, block: sentBlock, interval: configs.GetLatency(fromNode.GetRegion(), toNode.GetRegion()) + delay}
	return &Task{Message: messages}
}

func CreateGetDataTask(fromNode Node, toNode Node, sentBlock *block.Block) *Task {
	messages := &Messages{messageType: "getdata", from: fromNode, to: toNode, block: sentBlock}
	return &Task{Message: messages}
}

func CreateGetBlockTxnTask(fromNode Node, toNode Node, sentBlock *block.Block) *Task {
	messages := &Messages{messageType: "getblocktxn", from: fromNode, to: toNode, block: sentBlock}
	return &Task{Message: messages}
}

func CreateBlockTxnTask(fromNode Node, toNode Node, sentBlock *block.Block, delay int64) *Task {
	messages := &Messages{messageType: "blocktxn", from: fromNode, to: toNode, block: sentBlock, interval: configs.GetLatency(fromNode.GetRegion(), toNode.GetRegion()) + delay}
	return &Task{Message: messages}
}

func (m *Messages) GetType() string {
	return m.messageType
}

func (m *Messages) GetInterval() int64 {
	if m.messageType == "getblocktxn" || m.messageType == "inv" || m.messageType == "getdata" {
		var latency int64 = configs.GetLatency(m.from.GetRegion(), m.to.GetRegion())
		return latency + 10
	} else {
		return m.interval
	}
}

func (m *Messages) GetFrom() Node {
	return m.from
}

func (m *Messages) GetBlock() *block.Block {
	return m.block
}

func (m *Messages) GetParent() *block.Block {
	return m.block.GetParent()
}

func (task *Messages) Run() {
	if task.messageType == "block" || task.messageType == "cmpctblock" {
		task.from.SendNextBlockMessage()

		f, err := os.OpenFile("data/output/output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err2 := f.WriteString("{" + "\"kind\":\"flow-block\"," + "\"content\":{" + "\"transmission-timestamp\":" + strconv.FormatInt(configs.GetCurrentTime()-task.GetInterval(), 10) + "," + "\"reception-timestamp\":" + strconv.FormatInt(configs.GetCurrentTime(), 10) + "," + "\"begin-node-id\":" + strconv.Itoa(task.from.GetID()) + "," + "\"end-node-id\":" + strconv.Itoa(task.to.GetID()) + "," + "\"block-id\":" + strconv.Itoa(task.block.GetID()) + "}" + "},")

		if err2 != nil {
			panic(err2)
		}
	}

	task.to.ReceiveMessage(task)
}
