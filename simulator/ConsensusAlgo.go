package simulator

import (
	"math"
	"math/big"
	"math/rand"
	"reflect"
)

type ConsensusAlgo struct {
	selfNode *Node
}

var index int = -1

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
	var p = new(big.Float).Quo(big.NewFloat(1.0), big.NewFloat(0).SetInt(difficulty))
	var u = rand.Float64()
	if p.Cmp(big.NewFloat(math.Pow(2, -53))) <= 0 {
		return nil
	} else {
		index++
		float, _ := p.Float64()
		return &MintingTask{selfNode, selfNode.block, int64(math.Round(math.Log(u) / math.Log(
			1.0/float) / float64(selfNode.miningPower))), difficulty, index}
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
