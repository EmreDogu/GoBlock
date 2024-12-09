package node

import (
	"bufio"
	"os"
	"strconv"
	"strings"

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

func (r *RoutingTable) GetNeighbors() []*Node {
	neighbors := make([]*Node, 0)
	neighbors = append(neighbors, r.outbound...)
	neighbors = append(neighbors, r.inbound...)
	return neighbors
}

func (rt *RoutingTable) InitTable(simulatedNodes []*Node) {
	file, err := os.Open("data/input/1000_20.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)

	skipUntil := rt.selfNode.nodeID
	currentLine := 0

	// Read line by line
	for scanner.Scan() {
		currentLine++

		if currentLine < skipUntil {
			continue
		}

		line := scanner.Text()
		info := strings.Split(line, ",")

		for i := 2; i < len(info); i++ {
			i, err := strconv.Atoi(info[i])

			if err != nil {
				// ... handle error
				panic(err)
			}
			rt.AddNeighbor(simulatedNodes[i])
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
	f, err := os.OpenFile("data/output/output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-link\"," + "\"content\":{" + "\"timestamp\":" + strconv.FormatInt(configs.GetCurrentTime(), 10) + "," + "\"begin-node-id\":" + strconv.Itoa(rt.selfNode.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(endNode.nodeID) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
