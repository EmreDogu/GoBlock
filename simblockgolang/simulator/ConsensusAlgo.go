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
	var receivedBlockParent *Block = &Block{}
	if receivedBlockHeight >= 1 {
		receivedBlockParent = receivedBlock.GetBlockWithHeight(receivedBlockHeight - 1)
	} else {
		receivedBlockParent = nil
	}
	return (receivedBlockHeight == 0 || receivedBlock.difficulty >= receivedBlockParent.nextDifficulty) && (reflect.ValueOf(currentBlock.parent).IsNil() || receivedBlock.totalDifficulty > currentBlock.totalDifficulty)
}

func (ca *ConsensusAlgo) Minting() *MintingTask {
	var selfNode = ca.selfNode
	var parent = selfNode.block
	var difficulty = parent.nextDifficulty
	var p = 1.0 / float64(difficulty)
	var u = rand.Float64()
	if p <= math.Pow(2, -53) {
		return nil
	} else {
		return &MintingTask{selfNode, selfNode.block, math.Log(u) / math.Log(
			1.0-p) / float64(selfNode.miningPower), difficulty, 0}
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
