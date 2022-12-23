package simulator

import (
	"math/rand"
)

type RoutingTable struct {
	selfNode *Node
	outbound []*Node
	inbound  []*Node
}

func (rt *RoutingTable) GetOutbound() []*Node {
	return rt.outbound
}

func (rt *RoutingTable) GetInbound() []*Node {
	return rt.inbound
}

func GetNumConnection(this *Node) int {
	return this.numConnection
}

func (rt *RoutingTable) GetNeighbors() []*Node {
	var neighbours []*Node
	neighbours = append(neighbours, rt.GetOutbound()...)
	neighbours = append(neighbours, rt.GetInbound()...)
	return neighbours
}

func (rt *RoutingTable) AddNeighbor(to *Node) bool {
	if to == rt.selfNode || contains(rt.outbound, to) || contains(rt.inbound, to) || len(rt.outbound) >= rt.selfNode.numConnection {
		return false
	} else {
		rt.outbound = append(rt.outbound, to)
		to.routingTable.inbound = append(to.routingTable.inbound, rt.selfNode)
		printAddLink(to, rt.selfNode)
		printAddLink(rt.selfNode, to)
		return true
	}
}

func (rt *RoutingTable) initTable() {
	candidates := []int{}
	for i := 0; i < len(GetSimulatedNodes()); i++ {
		candidates = append(candidates, i)
	}
	//rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	for i := 0; i < len(candidates); i++ {
		if len(rt.outbound) < rt.selfNode.numConnection {
			rt.AddNeighbor(GetSimulatedNodes()[i])
		} else {
			break
		}
	}
}
