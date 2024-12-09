package node

import (
	"math/rand"

	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
	"github.com/EmreDogu/GoBlock/internal/network/p2p"
)

type Consensus struct {
	selfNode *Node
}

func NewConsensus(node *Node) *Consensus {
	return &Consensus{
		selfNode: node}
}

func (c *Consensus) Mining() *p2p.Task {
	var random int64
	if configs.GetCurrentTime() == 0 {
		random = int64(rand.Int())
	} else {
		random = rand.Int63n(configs.GetCurrentTime()*100-configs.GetCurrentTime()) + configs.GetCurrentTime()
	}
	return p2p.CreateMiningTask(c.selfNode, c.selfNode.block, random)
}

func (c *Consensus) IsReceivedBlockValid(receivedBlock *block.Block, currentBlock *block.Block) bool {
	receivedBlockHeight := receivedBlock.GetHeight()
	var receivedBlockParent *block.Block
	if receivedBlockHeight != 0 {
		receivedBlockParent = receivedBlock.GetBlockWithHeight(receivedBlockHeight - 1)
	} else {
		receivedBlockParent = nil
	}
	return (receivedBlockHeight == 0 || (receivedBlockHeight == receivedBlockParent.GetHeight()+1)) && (currentBlock == nil || receivedBlockHeight > currentBlock.GetHeight())
}

func (c *Consensus) GenesisBlock() *block.Block {
	return block.BlockGenesisBlock(c.selfNode)
}
