package simulator

import "math/big"

type MintingTask struct {
	minter     *Node
	parent     *Block
	interval   int64
	difficulty *big.Int
	index      int
}

func (task *MintingTask) Run() {
	block := MakeBlock(task.parent, task.minter, GetCurrentTime(), task.difficulty)
	task.minter.ReceiveBlock(block)
}

func (task *MintingTask) GetParent() *Block {
	return task.parent
}
