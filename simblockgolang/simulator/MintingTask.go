package simulator

type MintingTask struct {
	minter     *Node
	parent     *Block
	interval   float64
	difficulty int
	index      int
}

func (task *MintingTask) Run() {
	block := MakeBlock(task.parent, task.minter, GetCurrentTime(), task.difficulty)
	task.minter.ReceiveBlock(block)
}
