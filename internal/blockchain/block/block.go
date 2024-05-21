package block

import (
	"reflect"
)

type Node interface {
}

type Block struct {
	id     int
	height int
	parent *Block
	minter Node
	time   int64
}

var latestId int = 0

func NewBlock(parent *Block, minter Node, time int64) *Block {
	var height int

	if reflect.ValueOf(parent).IsNil() {
		height = 0
	} else {
		height = parent.height + 1
	}
	block := &Block{latestId, height, parent, minter, time}
	latestId++
	return block
}

func (b *Block) GetHeight() int {
	return b.height
}

func (b *Block) GetParent() *Block {
	return b.parent
}

func (b *Block) GetTime() int64 {
	return b.time
}

func (b *Block) GetBlock() *Block {
	return b
}

func (b *Block) GetID() int {
	return b.id
}

func (b *Block) GetBlockWithHeight(height int) *Block {
	if b.height == height {
		return b
	} else {
		return b.parent.GetBlockWithHeight(height)
	}
}

func (b *Block) IsOnSameChainAs(block *Block) bool {
	if block == nil {
		return false
	} else if b.height <= block.height {
		return b == block.GetBlockWithHeight(b.height)
	} else {
		return b.GetBlockWithHeight(block.height) == block
	}
}

func BlockGenesisBlock(minter Node) *Block {
	return NewBlock(nil, minter, 0)
}
