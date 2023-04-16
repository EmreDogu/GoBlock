package simulator

import (
	"math"
	"math/rand"
	"reflect"
)

type ConsensusAlgo struct {
	selfNode *Node
}

func (ca *ConsensusAlgo) IsReceivedBlockValid(receivedBlock *Block, currentBlock *Block) bool {
	if !(isTest(receivedBlock)) {
		return false
	}
	receivedBlockHeight := receivedBlock.height
	var receivedBlockParent *Block
	if receivedBlockHeight == 0 {
		receivedBlockParent = nil
	} else {
		receivedBlockParent = receivedBlock.GetBlockWithHeight(receivedBlockHeight - 1)
	}
	return (receivedBlockHeight == 0 || receivedBlock.difficulty.Cmp(receivedBlockParent.nextDifficulty) >= 0) && (reflect.ValueOf(currentBlock.parent).IsNil() || receivedBlock.totalDifficulty.Cmp(currentBlock.totalDifficulty) > 0)
}

func (ca *ConsensusAlgo) Minting() *MintingTask {
	var selfNode = ca.selfNode
	var parent = selfNode.block
	var difficulty = parent.nextDifficulty
	var p = 1.0 / float64(difficulty.Int64())
	var u = rand.Float64()
	if p <= math.Pow(2, -53) {
		return nil
	} else {
		return &MintingTask{selfNode, selfNode.block, int64(math.Round(math.Log(u) / math.Log(
			1.0/p) / float64(selfNode.miningPower))), difficulty, 0}
	}
}

func (ca *ConsensusAlgo) GenesisBlock() *Block {
	return BlockGenesisBlock(ca.selfNode)
}

func isTest(t interface{}) bool {
	switch t.(type) {
	case *Block:
		return true
	default:
		return false
	}
}
