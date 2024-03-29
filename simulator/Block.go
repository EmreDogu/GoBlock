package simulator

import (
	"math/big"
	"reflect"
)

type Block struct {
	id              int
	height          int
	parent          *Block
	minter          *Node
	time            int64
	difficulty      *big.Int
	totalDifficulty *big.Int
	nextDifficulty  *big.Int
	route           map[int][]int
}

var genesisNextDifficulty *big.Int
var latestid int = 0

func MakeBlock(parent *Block, minter *Node, time int64, difficulty *big.Int) *Block {
	var height int
	var totalDifficulty *big.Int
	var nextDifficulty *big.Int

	if reflect.ValueOf(parent).IsNil() {
		height = 0
		totalDifficulty = difficulty
		nextDifficulty = genesisNextDifficulty
	} else {
		height = parent.GetHeight() + 1
		totalDifficulty = parent.GetDifficulty().Add(parent.difficulty, difficulty)
		nextDifficulty = parent.GetNextDifficulty()
	}
	block := &Block{latestid, height, parent, minter, time, difficulty, totalDifficulty, nextDifficulty, make(map[int][]int)}
	latestid++
	return block
}

func BlockGenesisBlock(minter *Node) *Block {
	totalMiningPower := 0
	for i := 0; i < len(GetSimulatedNodes()); i++ {
		totalMiningPower += GetSimulatedNodes()[i].miningPower
	}
	genesisNextDifficulty = big.NewInt(int64(totalMiningPower) * int64(GetTargetInterval()))
	return MakeBlock(nil, minter, 0, big.NewInt(0))
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

func (b *Block) GetTime() int64 {
	return b.time
}

func (b *Block) GetHeight() int {
	return b.height
}

func (b *Block) GetID() int {
	return b.id
}

func (b *Block) GetDifficulty() *big.Int {
	return b.difficulty
}

func (b *Block) GetTotalDifficulty() *big.Int {
	return b.totalDifficulty
}

func (b *Block) GetNextDifficulty() *big.Int {
	return b.nextDifficulty
}
