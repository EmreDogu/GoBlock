package simulator

import (
	"math/rand"
	"os"
	"reflect"
	settings "simblockgolang/settings"
	"strconv"
)

type Node struct {
	nodeID            int
	numConnection     int
	region            int
	miningPower       int
	routingTableName  string
	consensusAlgoName string
	useCBR            bool
	isChurnNode       bool
	routingTable      *RoutingTable
	consensusAlgo     *ConsensusAlgo
	block             *Block
	orphans           []*Block
	mintingTask       *MintingTask
	sendingBlock      bool
	messageQue        []*MessageTask
	downloadingBlocks map[*Block]void
}

type void struct{}

var member void
var processingTime = 2

func MakeNode(nodeID int, numConnection int, region int, miningPower int, routingTableName string, consensusAlgoName string, useCBR bool, isChurnNode bool) *Node {
	node := &Node{nodeID, numConnection, region, miningPower, routingTableName, consensusAlgoName, useCBR, isChurnNode, nil, nil, nil, nil, nil, false, nil, nil}
	node.routingTable = &RoutingTable{node, []*Node{}, []*Node{}}
	node.consensusAlgo = &ConsensusAlgo{node}
	node.block = &Block{}
	node.mintingTask = &MintingTask{}
	node.orphans = []*Block{}
	node.messageQue = []*MessageTask{}
	node.downloadingBlocks = map[*Block]void{}
	return node
}

func (n *Node) JoinNetwork() {
	n.routingTable.initTable()
}

func (n *Node) GenesisBlock() {
	genesis := n.consensusAlgo.GenesisBlock()
	n.ReceiveBlock(genesis)
}

func contains(s []*Node, e *Node) bool {
	for _, a := range s {
		if a.nodeID == e.nodeID {
			return true
		}
	}
	return false
}

func ContainsBlock(s []*Block, e *Block) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsMapBlock(s map[*Block]void, e *Block) bool {
	_, ok := s[e]
	if ok {
		return true
	} else {
		return false
	}
}

func (n *Node) Minting() {
	task := n.consensusAlgo.Minting()
	n.mintingTask = task
	if !reflect.ValueOf(task).IsNil() {
		putMintingTask(task)
	}
}

func (n *Node) SendInv(block *Block) {
	for i := range n.routingTable.GetNeighbors() {
		InvMessageTask(n, n.routingTable.GetNeighbors()[i], block, GetLatency(n.region, n.routingTable.GetNeighbors()[i].region)+10)
	}
}

func (n *Node) ReceiveBlock(block *Block) {
	if n.consensusAlgo.IsReceivedBlockValid(block, n.block) {
		if !reflect.ValueOf(n.block.parent).IsNil() && !n.block.IsOnSameChainAs(block) {
			n.AddOrphans(n.block, block)
		}
		n.AddToChain(block)
		n.Minting()
		n.SendInv(block)
	} else if !ContainsBlock(n.orphans, block) && !block.IsOnSameChainAs(n.block) {
		n.AddOrphans(block, n.block)
		arriveBlock(block, n)
	}
}

func (n *Node) ReceiveMessage(message *MessageTask) {
	from := message.from

	if message.messageType == "InvMessageTask" {
		block := message.block
		if !ContainsBlock(n.orphans, block) && !containsMapBlock(n.downloadingBlocks, block) {
			if n.consensusAlgo.IsReceivedBlockValid(block, n.block) {
				RecMessageTask(n, from, block, GetLatency(from.region, n.region)+10)
				n.downloadingBlocks[block] = member
			} else if block.IsOnSameChainAs(n.block) {
				// get new orphan block
				RecMessageTask(n, from, block, GetLatency(from.region, n.region)+10)
				n.downloadingBlocks[block] = member
			}
		}
	}

	if message.messageType == "RecMessageTask" {
		n.messageQue = append(n.messageQue, message)
		if !n.sendingBlock {
			n.SendNextBlockMessage()
		}
	}

	if message.messageType == "GetBlockTxnMessageTask" {
		n.messageQue = append(n.messageQue, message)
		if !n.sendingBlock {
			n.SendNextBlockMessage()
		}
	}

	if message.messageType == "CmpctBlockMessageTask" {
		block := message.block
		//rand.Seed(time.Now().UnixNano())
		random := rand.Float64()
		var CBRfailureRate float64
		if n.isChurnNode {
			CBRfailureRate = settings.CBR_FAILURE_RATE_FOR_CHURN_NODE
		} else {
			CBRfailureRate = settings.CBR_FAILURE_RATE_FOR_CONTROL_NODE
		}
		var success bool
		if random > CBRfailureRate {
			success = true
		} else {
			success = false
		}
		if success {
			delete(n.downloadingBlocks, block)
			n.ReceiveBlock(block)
		} else {
			GetBlockTxnMessageTask(n, from, block, GetLatency(from.region, n.region)+10)
		}
	}

	if message.messageType == "BlockMessageTask" {
		block := message.block
		delete(n.downloadingBlocks, block)
		n.ReceiveBlock(block)
	}
}

func GetFailedBlockSize(this *Node) int {
	//rand.Seed(time.Now().UnixNano())
	if this.isChurnNode {
		index := rand.Intn(len(settings.CBR_FAILURE_BLOCK_SIZE_DISTRIBUTION_FOR_CHURN_NODE))
		return int(float64(settings.BLOCK_SIZE) * settings.CBR_FAILURE_BLOCK_SIZE_DISTRIBUTION_FOR_CHURN_NODE[index])
	} else {
		index := rand.Intn(len(settings.CBR_FAILURE_BLOCK_SIZE_DISTRIBUTION_FOR_CONTROL_NODE))
		return int(float64(settings.BLOCK_SIZE) * settings.CBR_FAILURE_BLOCK_SIZE_DISTRIBUTION_FOR_CONTROL_NODE[index])
	}
}

func (n *Node) SendNextBlockMessage() {
	if len(n.messageQue) > 0 {
		to := n.messageQue[0].from
		bandwidth := GetBandwidth(n.region, to.region)

		if n.messageQue[0].messageType == "RecMessageTask" {
			block := n.messageQue[0].block
			// If use compact block relay.
			if n.messageQue[0].from.useCBR && n.useCBR {
				// Convert bytes to bits and divide by the bandwidth expressed as bit per millisecond, add
				// processing time.
				delay := settings.COMPACT_BLOCK_SIZE*8/(bandwidth/1000) + processingTime

				// Send compact block message.
				CmpctBlockMessageTask(n, to, block, float64(delay))
			} else {
				// Else use lagacy protocol.
				delay := settings.BLOCK_SIZE*8/(bandwidth/1000) + processingTime
				BlockMessageTask(n, to, block, float64(delay))
			}
		} else if n.messageQue[0].messageType == "GetBlockTxnMessageTask" {
			// Else from requests missing transactions.
			block := n.messageQue[0].block
			delay := GetFailedBlockSize(n)*8/(bandwidth/1000) + processingTime
			BlockMessageTask(n, to, block, float64(delay))
		}

		n.sendingBlock = true
		n.messageQue = n.messageQue[1:]
	} else {
		n.sendingBlock = false
	}
}

func (n *Node) AddToChain(block *Block) {
	if !reflect.ValueOf(n.mintingTask).IsNil() {
		removeTask(n.mintingTask)
		n.mintingTask = nil
	}
	n.block = block
	printAddBlock(n, block)
	arriveBlock(block, n)
}

func (n *Node) AddOrphans(orphanBlock *Block, validBlock *Block) {
	if orphanBlock != validBlock {
		n.orphans = append(n.orphans, orphanBlock)
		for i := 0; i < len(n.orphans); i++ {
			if n.orphans[i] == validBlock {
				n.orphans[i] = n.orphans[len(n.orphans)-1]
				n.orphans = n.orphans[:len(n.orphans)-1]
			}
		}
		if reflect.ValueOf(validBlock).IsNil() || orphanBlock.height > validBlock.height {
			n.AddOrphans(orphanBlock.parent, validBlock)
		} else if orphanBlock.height == validBlock.height {
			n.AddOrphans(orphanBlock.parent, validBlock.parent)
		} else {
			n.AddOrphans(orphanBlock, validBlock.parent)
		}
	}
}

func printAddLink(this *Node, to *Node) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-link\"," + "\"content\":{" + "\"timestamp\":" + strconv.Itoa(int(GetCurrentTime())) + "," + "\"begin-node-id\":" + strconv.Itoa(this.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(to.nodeID) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}

func (n *Node) GetBlock() *Block {
	return n.block
}

func (n *Node) GetOrphans() []*Block {
	return n.orphans
}

func printAddBlock(this *Node, block *Block) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-block\"," + "\"content\":{" + "\"timestamp\":" + strconv.Itoa(int(GetCurrentTime())) + "," + "\"node-id\":" + strconv.Itoa(this.nodeID) + "," + "\"block-id\":" + strconv.Itoa(block.id) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
