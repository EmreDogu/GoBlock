package simulator

import (
	"math"
	"os"
	"reflect"
	settings "simblockgolang/settings"
	"strconv"
)

type Node struct {
	nodeID            int
	region_id         int
	ip                string
	latitude          float64
	longitude         float64
	location          string
	downloadSpeed     int
	uploadSpeed       int
	miningPower       int
	routingTable      *RoutingTable
	consensusAlgo     *ConsensusAlgo
	block             *Block
	orphans           map[*Block]void
	mintingTask       *MintingTask
	sendingBlock      bool
	messageQue        []*MessageTask
	downloadingBlocks map[*Block]void
}

type void struct{}

var member void
var processingTime = 2

func MakeNode(nodeID int, numConnection int, ip string, region_id int, latitude float64, longitude float64, location string, miningpower int, downloadSpeed int, uploadSpeed int) *Node {
	node := &Node{nodeID, region_id, ip, latitude, longitude, location, miningpower, downloadSpeed, uploadSpeed, nil, nil, nil, nil, nil, false, nil, nil}
	node.routingTable = &RoutingTable{node, numConnection, []*Node{}, []*Node{}}
	node.consensusAlgo = &ConsensusAlgo{node}
	node.block = &Block{}
	node.mintingTask = &MintingTask{}
	node.orphans = map[*Block]void{}
	node.messageQue = []*MessageTask{}
	node.downloadingBlocks = map[*Block]void{}
	return node
}

func (n *Node) JoinNetwork(CON_ALG int) {
	n.routingTable.initTable(CON_ALG)
}

func (n *Node) JoinNetworkBCBSN(CON_ALG int, nodelist []*Node) {
	n.routingTable.initTableBCBSN(CON_ALG, nodelist)
}

func (n *Node) GenesisBlock() {
	genesis := n.consensusAlgo.GenesisBlock()
	n.ReceiveBlock(genesis, nil, n)
}

func contains(s []*Node, e *Node) bool {
	for _, a := range s {
		if a.nodeID == e.nodeID {
			return true
		}
	}
	return false
}

func containsint(s []int, e int) bool {
	for _, a := range s {
		if a == e {
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
	for i := range n.routingTable.outbound {
		InvMessageTask(n, n.routingTable.outbound[i], block, GetLatency(n.region_id, n.routingTable.outbound[i].region_id)+10)
	}
}

func (n *Node) ReceiveBlock(block *Block, from *Node, to *Node) {
	if n.consensusAlgo.IsReceivedBlockValid(block, n.block) {
		if !reflect.ValueOf(n.block).IsNil() && !n.block.IsOnSameChainAs(block) {
			n.AddOrphans(n.block, block)
		}

		nodes := make([]int, 0)
		if !reflect.ValueOf(from).IsNil() && len(block.route[from.nodeID]) == 0 {
			for j := 0; j < len(block.route[from.nodeID]); j++ {
				nodes = append(nodes, block.route[from.nodeID][j])
			}
			block.route[to.nodeID] = nodes
			if containsint(block.route[to.nodeID], from.nodeID) {
				block.route[to.nodeID] = append(block.route[to.nodeID], from.nodeID)
			}
		} else {
			nodes = append(nodes, to.nodeID)
			block.route[to.nodeID] = nodes
		}

		n.AddToChain(block)
		n.Minting()
		n.SendInv(block)
	} else if !containsMapBlock(n.orphans, block) && !block.IsOnSameChainAs(n.block) {
		n.AddOrphans(block, n.block)
		arriveBlock(block, n, settings.NUM_OF_NODES)
	}
}

func (n *Node) ReceiveMessage(message *MessageTask) {
	from := message.from

	if message.messageType == "InvMessageTask" {
		block := message.block
		time := GetCurrentTime() - block.time
		if !reflect.ValueOf(from).IsNil() && (settings.NEIGH_SEL == "E" || settings.NEIGH_SEL == "e") {
			if from.nodeID != n.nodeID {
				if settings.Matrix[from.nodeID][n.nodeID] != 0 {
					settings.Matrix[from.nodeID][n.nodeID] = int64(math.Round(0.7*float64(settings.Matrix[from.nodeID][n.nodeID]) + 0.3*float64(time)))
				} else {
					settings.Matrix[from.nodeID][n.nodeID] = time
				}
			}
		}

		if !containsMapBlock(n.orphans, block) && !containsMapBlock(n.downloadingBlocks, block) {
			if n.consensusAlgo.IsReceivedBlockValid(block, n.block) {
				RecMessageTask(n, from, block, GetLatency(from.region_id, n.region_id)+10)
				n.downloadingBlocks[block] = member
			} else if !block.IsOnSameChainAs(n.block) {
				RecMessageTask(n, from, block, GetLatency(from.region_id, n.region_id)+10)
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

	if message.messageType == "BlockMessageTask" {
		block := message.block
		delete(n.downloadingBlocks, block)
		n.ReceiveBlock(block, from, n)
	}
}

func (n *Node) SendNextBlockMessage() {
	if len(n.messageQue) > 0 {
		to := n.messageQue[0].from
		bandwidth := GetBandwidth(n.region_id, to.region_id)

		if n.messageQue[0].messageType == "RecMessageTask" {
			block := n.messageQue[0].block
			delay := settings.BLOCK_SIZE*8/(bandwidth/1000) + processingTime
			BlockMessageTask(n, to, block, int64(delay))
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
	n.printAddBlock(block)
	arriveBlock(block, n, settings.NUM_OF_NODES)
}

func (n *Node) AddOrphans(orphanBlock *Block, validBlock *Block) {
	if orphanBlock != validBlock {
		n.orphans[orphanBlock] = member
		delete(n.orphans, validBlock)
		if reflect.ValueOf(validBlock).IsNil() || orphanBlock.height > validBlock.height {
			n.AddOrphans(orphanBlock.parent, validBlock)
		} else if orphanBlock.height == validBlock.height {
			n.AddOrphans(orphanBlock.parent, validBlock.parent)
		} else {
			n.AddOrphans(orphanBlock, validBlock.parent)
		}
	}
}

func (n *Node) GetBlock() *Block {
	return n.block
}

func (n *Node) GetRoutingTable() *RoutingTable {
	return n.routingTable
}

func (n *Node) GetOrphans() map[*Block]void {
	return n.orphans
}

func (n *Node) printAddBlock(block *Block) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-block\"," + "\"content\":{" + "\"timestamp\":" + strconv.Itoa(int(GetCurrentTime())) + "," + "\"node-id\":" + strconv.Itoa(n.nodeID) + "," + "\"block-id\":" + strconv.Itoa(block.id) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
