package simulator

import "reflect"

type Block struct {
	id              int
	height          int
	parent          *Block
	minter          *Node
	time            int64
	difficulty      int
	totalDifficulty int
	nextDifficulty  int
}

var genesisNextDifficulty int
var latestid int

func MakeBlock(parent *Block, minter *Node, time int64, difficulty int) *Block {
	var height int
	var totalDifficulty int
	var nextDifficulty int

	if reflect.ValueOf(parent).IsNil() {
		height = 0
		totalDifficulty = difficulty
		nextDifficulty = genesisNextDifficulty
	} else {
		height = parent.height + 1
		totalDifficulty = parent.difficulty + difficulty
		nextDifficulty = parent.nextDifficulty
	}
	block := &Block{latestid, height, parent, minter, time, difficulty, totalDifficulty, nextDifficulty}
	latestid++
	return block
}

func BlockGenesisBlock(minter *Node) *Block {
	totalMiningPower := 0
	for i := 0; i < len(GetSimulatedNodes()); i++ {
		totalMiningPower += GetSimulatedNodes()[i].miningPower
	}
	genesisNextDifficulty = totalMiningPower * GetTargetInterval()
	return MakeBlock(nil, minter, 0, 0)
}

func (b *Block) GetBlockWithHeight(height int) *Block {
	if b.height == height {
		return b
	} else {
		return b.parent.GetBlockWithHeight(height)
	}
}

func (b *Block) IsOnSameChainAs(block *Block) bool {
	if reflect.ValueOf(block).IsNil() {
		return false
	} else if b.height <= block.height {
		return (b == block.GetBlockWithHeight(b.height))
	} else {
		return (b.GetBlockWithHeight(block.height) == block)
	}
}

func (b *Block) GetParent() *Block {
	return b.parent
}

func GetTime(this *Block) int64 {
	return this.time
}

func GetHeight(this *Block) int {
	return this.height
}

func GetID(this *Block) int {
	return this.id
}
