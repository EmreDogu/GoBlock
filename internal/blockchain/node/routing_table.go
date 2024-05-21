package node

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/EmreDogu/GoBlock/configs"
)

type RoutingTable struct {
	selfNode      *Node
	numConnection int
	numHighCon    int
	outbound      []*Node
	highoutbound  []*Node
	inbound       []*Node
	highinbound   []*Node
}

func NewRoutingTable(node *Node, numCon int, numHighCon int) *RoutingTable {
	return &RoutingTable{
		selfNode:      node,
		numConnection: numCon,
		numHighCon:    numHighCon,
		outbound:      []*Node{},
		highoutbound:  []*Node{},
		inbound:       []*Node{},
		highinbound:   []*Node{}}
}

func (rt *RoutingTable) InitTable(simulatedNodes []*Node) {
	candidates := make([]int, len(simulatedNodes))
	for i := range candidates {
		candidates[i] = i
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, candidate := range candidates {
		if len(rt.outbound) < rt.numConnection {
			rt.AddNeighbor(simulatedNodes[candidate])
		} else {
			break
		}
	}
}

func (rt *RoutingTable) AddNeighbor(to *Node) bool {
	if to.nodeID == rt.selfNode.nodeID || ContainsNode(rt.outbound, to) || ContainsNode(rt.inbound, to) || len(rt.outbound) >= rt.numConnection {
		return false
	} else if to.useCBR && rt.selfNode.useCBR && len(rt.highoutbound) < rt.numHighCon && len(to.routingTable.inbound) < to.routingTable.numHighCon {
		rt.highoutbound = append(rt.highoutbound, to)
		to.routingTable.highinbound = append(to.routingTable.highinbound, rt.selfNode)
		rt.printAddLink(to)
		to.routingTable.printAddLink(rt.selfNode)
		return true
	} else {
		rt.outbound = append(rt.outbound, to)
		to.routingTable.inbound = append(to.routingTable.inbound, rt.selfNode)
		rt.printAddLink(to)
		to.routingTable.printAddLink(rt.selfNode)
		return true
	}
}

func (rt *RoutingTable) printAddLink(endNode *Node) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-link\"," + "\"content\":{" + "\"timestamp\":" + strconv.FormatInt(configs.GetCurrentTime(), 10) + "," + "\"begin-node-id\":" + strconv.Itoa(rt.selfNode.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(endNode.nodeID) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
