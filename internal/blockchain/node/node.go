package node

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
	"github.com/EmreDogu/GoBlock/internal/network/p2p"
)

type Node struct {
	nodeID            int
	region            int
	useCBR            bool
	routingTable      *RoutingTable
	consensus         *Consensus
	block             *block.Block
	sendingBlock      bool
	orphans           map[*block.Block]void
	downloadingBlocks map[*block.Block]void
	messageQue        []*p2p.Messages
	miningTask        *p2p.Task
}

type SimulatorLink struct {
	simulator Simulator
}

type Simulator interface {
	ArriveBlock(*block.Block, *Node)
	RemoveTask(*p2p.Task)
	PutTask(*p2p.Task)
}

type void struct{}

var simulatedNodes []*Node
var member void
var sl *SimulatorLink

var processingTime int = 2

func New(nodeID int, numCon int, region int, numHighCon int, useCBR bool) *Node {
	node := &Node{
		nodeID:            nodeID,
		region:            region,
		useCBR:            useCBR,
		routingTable:      nil,
		consensus:         nil,
		block:             nil,
		sendingBlock:      false,
		orphans:           map[*block.Block]void{},
		downloadingBlocks: map[*block.Block]void{},
		messageQue:        []*p2p.Messages{},
		miningTask:        nil}
	node.routingTable = NewRoutingTable(node, numCon, numHighCon)
	node.consensus = NewConsensus(node)
	return node
}

func (n *Node) GetID() int {
	return n.nodeID
}

func (n *Node) GetRegion() int {
	return n.region
}

func (n *Node) GetUseCBR() bool {
	return n.useCBR
}

func (n *Node) GetBlock() *block.Block {
	return n.block
}

func (n *Node) GetOrphans() map[*block.Block]void {
	return n.orphans
}

func ContainsNode(s []*Node, e *Node) bool {
	for _, a := range s {
		if a.nodeID == e.nodeID {
			return true
		}
	}
	return false
}

func ContainsMapBlock(s map[*block.Block]void, e *block.Block) bool {
	_, ok := s[e]
	if ok {
		return true
	} else {
		return false
	}
}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func NewSimulatorLink(s Simulator) *SimulatorLink {
	return &SimulatorLink{simulator: s}
}

func ReceiveSimulatorLink(link *SimulatorLink) {
	sl = link
}

func ReceiveSimulatedNodes(nodes []*Node) {
	simulatedNodes = nodes
}

func (n *Node) JoinNetwork() {
	n.routingTable.InitTable(simulatedNodes)
}

func (n *Node) GenesisBlock() {
	genesis := n.consensus.GenesisBlock()
	n.ReceiveBlock(genesis)
}

func (n *Node) addToChain(newBlock *block.Block) {
	if n.miningTask != nil {
		sl.simulator.RemoveTask(n.miningTask)
		n.miningTask = nil
	}

	n.block = newBlock
	n.printAddBlock(newBlock)
	sl.simulator.ArriveBlock(newBlock, n)
}

func (n *Node) addOrphans(orphanBlock *block.Block, validBlock *block.Block) {
	if orphanBlock != validBlock {
		n.orphans[orphanBlock] = member
		delete(n.orphans, validBlock)

		if validBlock == nil || orphanBlock.GetHeight() > validBlock.GetHeight() {
			n.addOrphans(orphanBlock.GetParent(), validBlock)
		} else if orphanBlock.GetHeight() == validBlock.GetHeight() {
			n.addOrphans(orphanBlock.GetParent(), validBlock.GetParent())
		} else {
			n.addOrphans(orphanBlock, validBlock.GetParent())
		}
	}
}

func (n *Node) mining() {
	task := n.consensus.Mining()
	n.miningTask = task

	if task != nil {
		sl.simulator.PutTask(task)
	}
}

func (n *Node) sendBlock(newBlock *block.Block) {
	for j := range n.routingTable.highoutbound {
		bandwidth := configs.GetBandwidth(n.region, n.routingTable.highoutbound[j].GetRegion())
		delay := int64(configs.COMPACT_BLOCK_SIZE*8/(bandwidth/1000) + processingTime)
		task := p2p.CreateCmpctBlockTask(n, n.routingTable.highoutbound[j], newBlock, delay)
		sl.simulator.PutTask(task)
	}

	for i := range n.routingTable.outbound {
		task := p2p.CreateInvTask(n, n.routingTable.outbound[i], newBlock)
		sl.simulator.PutTask(task)
	}
}

func (n *Node) ReceiveBlock(block *block.Block) {
	if n.consensus.IsReceivedBlockValid(block, n.block) {
		if n.block != nil && n.block.IsOnSameChainAs(block) {
			n.addOrphans(n.block, block)
		}

		n.addToChain(block)
		n.mining()
		n.sendBlock(block)
	} else if !ContainsMapBlock(n.orphans, block) && !block.IsOnSameChainAs(n.block) {
		n.addOrphans(block, n.block)
		sl.simulator.ArriveBlock(block, n)
	}
}

func (n *Node) ReceiveMessage(message *p2p.Messages) {
	from := message.GetFrom()

	if message.GetType() == "inv" {
		block := message.GetBlock()
		if !ContainsMapBlock(n.orphans, block) && !ContainsMapBlock(n.downloadingBlocks, block) {
			if n.consensus.IsReceivedBlockValid(block, n.block) {
				task := p2p.CreateGetDataTask(n, from, block)
				sl.simulator.PutTask(task)
				n.downloadingBlocks[block] = member
			} else if !block.IsOnSameChainAs(n.block) {
				task := p2p.CreateGetDataTask(n, from, block)
				sl.simulator.PutTask(task)
				n.downloadingBlocks[block] = member
			}
		}
	}

	if message.GetType() == "getdata" {
		n.messageQue = append(n.messageQue, message)
		if !n.sendingBlock {
			n.SendNextBlockMessage()
		}
	}

	if message.GetType() == "getblocktxn" {
		n.messageQue = append(n.messageQue, message)
		if !n.sendingBlock {
			n.SendNextBlockMessage()
		}
	}

	if message.GetType() == "cmpctblock" {
		block := message.GetBlock()

		randomNumber := rand.Intn(100)

		if randomNumber < 80 {
			delete(n.downloadingBlocks, block)

			n.ReceiveBlock(block)
		} else {
			task := p2p.CreateGetBlockTxnTask(n, from, block)
			sl.simulator.PutTask(task)
		}
	}

	if message.GetType() == "block" {
		block := message.GetBlock()
		delete(n.downloadingBlocks, block)
		n.ReceiveBlock(block)
	}
}

func (n *Node) SendNextBlockMessage() {
	if len(n.messageQue) > 0 {
		to := n.messageQue[0].GetFrom()
		messageTask := &p2p.Task{}
		bandwidth := configs.GetBandwidth(n.region, to.GetRegion())

		if n.messageQue[0].GetType() == "getdata" {
			block := n.messageQue[0].GetBlock()

			if n.messageQue[0].GetFrom().GetUseCBR() && n.useCBR {
				delay := int64(configs.COMPACT_BLOCK_SIZE*8/(bandwidth/1000) + processingTime)
				messageTask = p2p.CreateCmpctBlockTask(n, to, block, delay)
			} else {
				delay := int64(configs.BLOCK_SIZE*8/(bandwidth/1000) + processingTime)
				messageTask = p2p.CreateBlockTask(n, to, block, delay)
			}
		} else if n.messageQue[0].GetType() == "getblocktxn" {
			block := n.messageQue[0].GetBlock()
			delay := int64(configs.COMPACT_BLOCK_SIZE*8/(bandwidth/1000) + processingTime)
			messageTask = p2p.CreateBlockTask(n, to, block, delay)
		}

		n.sendingBlock = true
		n.messageQue = n.messageQue[1:]
		sl.simulator.PutTask(messageTask)
	} else {
		n.sendingBlock = false
	}
}

func (n *Node) printAddBlock(block *block.Block) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-block\"," + "\"content\":{" + "\"timestamp\":" + strconv.FormatInt(configs.GetCurrentTime(), 10) + "," + "\"node-id\":" + strconv.Itoa(n.nodeID) + "," + "\"block-id\":" + strconv.Itoa(block.GetID()) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
