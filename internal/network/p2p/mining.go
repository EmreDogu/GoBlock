package p2p

import (
	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
)

type Mining struct {
	miner    Node
	parent   *block.Block
	interval int64
}

func CreateMiningTask(node Node, block *block.Block, value int64) *Task {
	mining := &Mining{miner: node, parent: block, interval: value}
	return &Task{Message: mining}
}

func (m *Mining) GetType() string {
	return "mining"
}

func (m *Mining) GetInterval() int64 {
	return m.interval
}

func (m *Mining) GetParent() *block.Block {
	return m.parent
}

func (m *Mining) Run() {
	block := block.NewBlock(m.parent, m.miner, configs.GetCurrentTime())
	m.miner.ReceiveBlock(block)
}
