package simulator

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type RoutingTable struct {
	selfNode      *Node
	numConnection int
	outbound      []*Node
	inbound       []*Node
}

func (rt *RoutingTable) GetOutbound() []*Node {
	return rt.outbound
}

func (rt *RoutingTable) GetInbound() []*Node {
	return rt.inbound
}

func (rt *RoutingTable) GetNeighbors() []*Node {
	var neighbours []*Node
	neighbours = append(neighbours, rt.GetOutbound()...)
	neighbours = append(neighbours, rt.GetInbound()...)
	return neighbours
}

func (rt *RoutingTable) AddNeighbor(to *Node) bool {
	if to == rt.selfNode || contains(rt.outbound, to) || contains(rt.inbound, to) || len(rt.outbound) >= rt.numConnection {
		return false
	} else {
		rt.outbound = append(rt.outbound, to)
		to.routingTable.inbound = append(to.routingTable.inbound, rt.selfNode)
		rt.printAddLink(to)
		return true
	}
}

func (rt *RoutingTable) removeInbounds() {
	for a := range rt.inbound {
		for i := range rt.inbound[a].routingTable.outbound {
			if rt.inbound[a].routingTable.outbound[i] == rt.selfNode {
				rt.inbound[a].routingTable.outbound = append(rt.inbound[a].routingTable.outbound[:i], rt.inbound[a].routingTable.outbound[i+1:]...)
			}
		}
	}
	rt.inbound = nil
}

func (rt *RoutingTable) removeOutbound(node *Node) {
	for i := range rt.outbound {
		if rt.outbound[i] == node {
			rt.outbound = append(rt.outbound[:i], rt.outbound[i+1:]...)
			rt.printRemoveLink(node)
		}
	}
}

func (rt *RoutingTable) removeInbound(node *Node) {
	for i := range rt.inbound {
		if rt.inbound[i] == node {
			rt.inbound = append(rt.inbound[:i], rt.inbound[i+1:]...)
			rt.printRemoveLink(node)
		}
	}
}

func (rt *RoutingTable) ReconnectAll(matrix [][]int64, numNodes int) {
	flag := false
	nodemap := make(map[int]*Node, 0)
	j := rt.selfNode.nodeID

	for i := 0; i < numNodes; i++ {
		if i == j || matrix[i][j] == 0 {
			continue
		} else {
			nodemap[int(matrix[i][j])] = GetSimulatedNodes()[i]
			matrix[i][j] = 0
		}
	}

	keys := make([]int, 0, len(nodemap))

	for k := range nodemap {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	if !(len(keys) == len(rt.inbound)) {
		difference := len(rt.inbound) - len(keys)
		GetSimulatedNodes()[j].routingTable.removeInbounds()

		for i := 0; i < len(keys); i++ {
			nodemap[keys[i]].routingTable.AddNeighbor(GetSimulatedNodes()[j])
		}

		candidates := []int{}
		for i := 0; i < len(GetSimulatedNodes()); i++ {
			if len(GetSimulatedNodes()[i].routingTable.outbound) < 8 {
				candidates = append(candidates, i)
				flag = true
			}
		}

		if flag {
			rand.Shuffle(len(candidates), func(i, j int) {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			})
			for k := 0; k < difference; k++ {
				if k < len(candidates) {
					GetSimulatedNodes()[candidates[k]].routingTable.AddNeighbor(GetSimulatedNodes()[j])
				}
			}
		}
	} else {
		nodemap[keys[len(keys)-1]].routingTable.removeOutbound(GetSimulatedNodes()[j])
		GetSimulatedNodes()[j].routingTable.removeInbound(nodemap[keys[len(keys)-1]])

		candidates := []int{}
		for i := 0; i < len(GetSimulatedNodes()); i++ {
			if len(GetSimulatedNodes()[i].routingTable.outbound) < 8 {
				candidates = append(candidates, i)
				flag = true
			}
		}

		if flag {
			rand.Shuffle(len(candidates), func(i, j int) {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			})
			GetSimulatedNodes()[candidates[0]].routingTable.AddNeighbor(GetSimulatedNodes()[j])
		}
	}
}

func (rt *RoutingTable) initTable(CON_ALG int) {
	if CON_ALG == 1 {
		//parsetan sonra
	} else if CON_ALG == 2 {
		var holderlist []string
		var nodelist []*Node = GetSimulatedNodes()

		for len(rt.selfNode.routingTable.outbound) < rt.numConnection {
			if len(holderlist) == len(nodelist)-1 {
				fmt.Println("Can't complete node id: " + strconv.Itoa(rt.selfNode.nodeID) + "'s outgoing connection to " + strconv.Itoa(rt.numConnection) + " Final count was: " + strconv.Itoa(len(rt.selfNode.routingTable.outbound)) + ". Moving on to next node\n")
				holderlist = nil
				break
			}

			var randomValue int

			for {
				randomValue = rand.Intn(len(nodelist))
				if !(nodelist[randomValue] == rt.selfNode && !stringcontains(holderlist, nodelist[randomValue].ip)) {
					break
				}
			}

			rt.AddNeighbor(nodelist[randomValue])
			holderlist = append(holderlist, nodelist[randomValue].ip)

			if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
				holderlist = nil
				break
			}
		}
	} else if CON_ALG == 3 {
		var holderlist []string
		var nodelist []*Node = GetSimulatedNodes()

		for len(rt.selfNode.routingTable.outbound) < rt.numConnection {
			if len(holderlist) == len(nodelist)-1 {
				fmt.Println("Can't complete node id: " + strconv.Itoa(rt.selfNode.nodeID) + "'s outgoing connection to " + strconv.Itoa(rt.numConnection) + " Final count was: " + strconv.Itoa(len(rt.selfNode.routingTable.outbound)) + ". Moving on to next node\n")
				holderlist = nil
				break
			}

			var randomValue int

			for {
				randomValue = rand.Intn(len(nodelist))
				if !(nodelist[randomValue] == rt.selfNode && !stringcontains(holderlist, nodelist[randomValue].ip)) {
					break
				}
			}

			if calculateDistance(rt.selfNode, nodelist[randomValue]) <= 2500 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 0 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					holderlist = nil
					break
				}

			} else {
				holderlist = append(holderlist, nodelist[randomValue].ip)
			}
		}
	} else if CON_ALG == 4 {
		var holderlist []string
		var nodelist []*Node = GetSimulatedNodes()
		var x int
		var z int

		for len(rt.selfNode.routingTable.outbound) < rt.numConnection {
			if len(holderlist) == len(nodelist)-1 {
				fmt.Println("Can't complete node id: " + strconv.Itoa(rt.selfNode.nodeID) + "'s outgoing connection to " + strconv.Itoa(rt.numConnection) + " Final count was: " + strconv.Itoa(len(rt.selfNode.routingTable.outbound)) + ". Moving on to next node\n")
				holderlist = nil
				break
			}

			var randomValue int

			for {
				randomValue = rand.Intn(len(nodelist))
				if !(nodelist[randomValue] == rt.selfNode && !stringcontains(holderlist, nodelist[randomValue].ip)) {
					break
				}
			}

			if x < 7 && calculateDistance(rt.selfNode, nodelist[randomValue]) <= 5500 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 0 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				x++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					z = 0
					holderlist = nil
					break
				}

			} else if z < 1 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 5500 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				z++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					z = 0
					holderlist = nil
					break
				}

			} else {
				holderlist = append(holderlist, nodelist[randomValue].ip)
			}
		}
	} else if CON_ALG == 5 {
		var holderlist []string
		var nodelist []*Node = GetSimulatedNodes()
		var x int
		var z int

		for len(rt.selfNode.routingTable.outbound) < rt.numConnection {
			if len(holderlist) == len(nodelist)-1 {
				fmt.Println("Can't complete node id: " + strconv.Itoa(rt.selfNode.nodeID) + "'s outgoing connection to " + strconv.Itoa(rt.numConnection) + " Final count was: " + strconv.Itoa(len(rt.selfNode.routingTable.outbound)) + ". Moving on to next node\n")
				holderlist = nil
				break
			}

			var randomValue int

			for {
				randomValue = rand.Intn(len(nodelist))
				if !(nodelist[randomValue] == rt.selfNode && !stringcontains(holderlist, nodelist[randomValue].ip)) {
					break
				}
			}

			if x < 4 && calculateDistance(rt.selfNode, nodelist[randomValue]) <= 5500 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 0 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				x++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					z = 0
					holderlist = nil
					break
				}

			} else if z < 4 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 5500 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				z++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					z = 0
					holderlist = nil
					break
				}

			} else {
				holderlist = append(holderlist, nodelist[randomValue].ip)
			}
		}
	} else if CON_ALG == 6 {
		var holderlist []string
		var nodelist []*Node = GetSimulatedNodes()
		var x int
		var y int
		var z int

		for len(rt.selfNode.routingTable.outbound) < rt.numConnection {
			if len(holderlist) == len(nodelist)-1 {
				fmt.Println("Can't complete node id: " + strconv.Itoa(rt.selfNode.nodeID) + "'s outgoing connection to " + strconv.Itoa(rt.numConnection) + " Final count was: " + strconv.Itoa(len(rt.selfNode.routingTable.outbound)) + ". Moving on to next node\n")
				holderlist = nil
				break
			}

			var randomValue int

			for {
				randomValue = rand.Intn(len(nodelist))
				if !(nodelist[randomValue] == rt.selfNode && !stringcontains(holderlist, nodelist[randomValue].ip)) {
					break
				}
			}

			if x < 3 && calculateDistance(rt.selfNode, nodelist[randomValue]) <= 2500 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 0 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				x++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					y = 0
					z = 0
					holderlist = nil
					break
				}

			} else if y < 3 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 2500 && calculateDistance(rt.selfNode, nodelist[randomValue]) < 5500 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				y++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					y = 0
					z = 0
					holderlist = nil
					break
				}

			} else if z < 2 && calculateDistance(rt.selfNode, nodelist[randomValue]) > 5500 {
				rt.AddNeighbor(nodelist[randomValue])
				holderlist = append(holderlist, nodelist[randomValue].ip)
				z++

				if len(rt.selfNode.routingTable.outbound) == rt.numConnection {
					x = 0
					y = 0
					z = 0
					holderlist = nil
					break
				}

			} else {
				holderlist = append(holderlist, nodelist[randomValue].ip)
			}
		}
	}
}

func (rt *RoutingTable) initTableBCBSN(CON_ALG int, nodelist []*Node) {
	if CON_ALG == 7 {
		flag1 := true
		flag2 := true
		var supernodelist []*Node
		var holderlist []string

		for i := 0; i < len(nodelist); i++ {
			if nodelist[i].region_id == 0 && flag1 {
				flag1 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 1 && flag2 {
				flag2 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if len(supernodelist) == 2 {
				supernodelist[0].routingTable.AddNeighbor(supernodelist[1])
				supernodelist[1].routingTable.AddNeighbor(supernodelist[0])
				break
			}
		}

		for j := 0; j < len(nodelist); j++ {
			if calculateDistance(supernodelist[0], nodelist[j]) <= calculateDistance(supernodelist[1], nodelist[j]) {
				nodelist[j].region_id = 0
				nodelist[j].routingTable.AddNeighbor(supernodelist[0])
			} else {
				nodelist[j].region_id = 1
				nodelist[j].routingTable.AddNeighbor(supernodelist[1])
			}
		}

		var i int

		for i < len(nodelist)-1 {
			for len(nodelist[i].routingTable.outbound) < nodelist[i].routingTable.numConnection {

				var randomValue int

				for {
					randomValue = rand.Intn(len(nodelist))
					if !(nodelist[randomValue] == nodelist[i] && !stringcontains(holderlist, nodelist[randomValue].ip)) {
						break
					}
				}

				if len(holderlist) == len(nodelist)-1 {
					fmt.Println("Can't complete node id: " + strconv.Itoa(nodelist[i].nodeID) + "'s outgoing connection to " + strconv.Itoa(nodelist[i].routingTable.numConnection) + " Final count was: " + strconv.Itoa(len(nodelist[i].routingTable.outbound)) + ". Moving on to next node\n")
					i++
					holderlist = nil
					if i == len(nodelist) {
						continue
					}
				}

				if calculateDistance(nodelist[i], nodelist[randomValue]) > 0 && nodelist[i].region_id == nodelist[randomValue].region_id {
					nodelist[i].routingTable.AddNeighbor(nodelist[randomValue])
					holderlist = append(holderlist, nodelist[randomValue].ip)

					if len(nodelist[i].routingTable.outbound) == nodelist[i].routingTable.numConnection {
						i++
						holderlist = nil
						break
					}

				} else {
					holderlist = append(holderlist, nodelist[randomValue].ip)
				}
			}
		}
	} else if CON_ALG == 8 {
		flag1 := true
		flag2 := true
		flag3 := true
		flag4 := true
		flag5 := true
		flag6 := true
		flag7 := true
		var supernodelist []*Node
		var holderlist []string

		for i := 0; i < len(nodelist); i++ {
			if nodelist[i].region_id == 0 && flag1 {
				flag1 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 1 && flag2 {
				flag2 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 2 && flag3 {
				flag3 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 3 && flag4 {
				flag4 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 4 && flag5 {
				flag5 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 5 && flag6 {
				flag6 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if nodelist[i].region_id == 6 && flag7 {
				flag7 = false
				supernodelist = append(supernodelist, nodelist[i])
				nodelist = append(nodelist[:i], nodelist[i+1:]...)
			}
			if len(supernodelist) == 7 {
				break
			}
		}

		for j := 0; j < len(supernodelist); j++ {
			for k := 0; k < len(supernodelist); k++ {
				if contains(supernodelist[j].routingTable.outbound, supernodelist[k]) {
					continue
				} else {
					supernodelist[j].routingTable.AddNeighbor(supernodelist[k])
					supernodelist[k].routingTable.AddNeighbor(supernodelist[j])
				}
			}
		}

		for j := 0; j < len(nodelist); j++ {
			distance := map[int]*Node{}
			for a := range supernodelist {
				distance[calculateDistance(supernodelist[a], nodelist[j])] = supernodelist[a]
			}
			keys := make([]int, 0, len(distance))

			for k := range distance {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			nodelist[j].region_id = keys[0]
			nodelist[j].routingTable.AddNeighbor(distance[keys[0]])
		}

		var i int

		for i < len(nodelist)-1 {
			for len(nodelist[i].routingTable.outbound) < nodelist[i].routingTable.numConnection {

				var randomValue int

				for {
					randomValue = rand.Intn(len(nodelist))
					if !(nodelist[randomValue] == nodelist[i] && !stringcontains(holderlist, nodelist[randomValue].ip)) {
						break
					}
				}

				if len(holderlist) == len(nodelist)-1 {
					fmt.Println("Can't complete node id: " + strconv.Itoa(nodelist[i].nodeID) + "'s outgoing connection to " + strconv.Itoa(nodelist[i].routingTable.numConnection) + " Final count was: " + strconv.Itoa(len(nodelist[i].routingTable.outbound)) + ". Moving on to next node\n")
					i++
					holderlist = nil
					if i == len(nodelist) {
						continue
					}
				}

				if calculateDistance(nodelist[i], nodelist[randomValue]) > 0 && nodelist[i].region_id == nodelist[randomValue].region_id {
					nodelist[i].routingTable.AddNeighbor(nodelist[randomValue])
					holderlist = append(holderlist, nodelist[randomValue].ip)

					if len(nodelist[i].routingTable.outbound) == nodelist[i].routingTable.numConnection {
						i++
						holderlist = nil
						break
					}

				} else {
					holderlist = append(holderlist, nodelist[randomValue].ip)
				}
			}
		}
	}
}

func stringcontains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func calculateDistance(from *Node, to *Node) int {
	var AVERAGE_RADIUS_OF_EARTH float64 = 6371

	var latDistance float64 = (from.latitude - to.latitude) * (math.Pi / 180)
	var lngDistance float64 = (from.longitude - to.longitude) * (math.Pi / 180)

	var a float64 = (math.Sin(latDistance/2) * math.Sin(latDistance/2)) +
		(math.Cos(from.latitude*(math.Pi/180)))*
			(math.Cos(to.longitude*(math.Pi/180)))*
			(math.Sin(lngDistance/2))*
			(math.Sin(lngDistance/2))

	var c float64 = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return (int)(math.Round(AVERAGE_RADIUS_OF_EARTH * c))
}

func (rt *RoutingTable) printAddLink(to *Node) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"add-link\"," + "\"content\":{" + "\"timestamp\":" + strconv.Itoa(int(GetCurrentTime())) + "," + "\"begin-node-id\":" + strconv.Itoa(rt.selfNode.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(to.nodeID) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}

func (rt *RoutingTable) printRemoveLink(to *Node) {
	f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{" + "\"kind\":\"remove-link\"," + "\"content\":{" + "\"timestamp\":" + strconv.Itoa(int(GetCurrentTime())) + "," + "\"begin-node-id\":" + strconv.Itoa(rt.selfNode.nodeID) + "," + "\"end-node-id\":" + strconv.Itoa(to.nodeID) + "}" + "},")

	if err2 != nil {
		panic(err2)
	}
}
