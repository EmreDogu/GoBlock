package simulator

import (
	"fmt"
	"reflect"
	"strconv"
)

var TargetInterval int
var simulatedNodes []*Node
var observedBlocks []*Block
var observedPropagations []map[int]int64

func SetTargetInterval(targetInterval int) {
	TargetInterval = targetInterval
}

func GetTargetInterval() int {
	return TargetInterval
}

func AddNode(node *Node) {
	simulatedNodes = append(simulatedNodes, node)
}

func GetSimulatedNodes() []*Node {
	return simulatedNodes
}

func arriveBlock(block *Block, this *Node, numNodes int) {
	if ContainsBlock(observedBlocks, block) {
		propagation := observedPropagations[IndexOf(observedBlocks, block)]
		propagation[this.nodeID] = GetCurrentTime() - block.time
	} else {
		if len(observedBlocks) > 10 {
			printPropagation(observedBlocks[0], observedPropagations[0], numNodes)
			observedBlocks = observedBlocks[1:]
			observedPropagations = observedPropagations[1:]
		}

		propagation := make(map[int]int64)
		propagation[this.nodeID] = GetCurrentTime() - block.time
		observedBlocks = append(observedBlocks, block)
		observedPropagations = append(observedPropagations, propagation)
	}
}

func printPropagation(block *Block, propagation map[int]int64, numNodes int) {
	fmt.Print(block)
	fmt.Print(":" + strconv.Itoa(block.height))
	fmt.Println("Network routes for every node:")

	for i := 0; i < numNodes+1; i++ {
		if !reflect.ValueOf(block.route[i]).IsNil() {
			fmt.Println("Network route for node ID " + strconv.Itoa(i) + ":")
			for j := 0; j < len(block.route[i]); j++ {
				fmt.Print((block.route[i])[j])
				if j+1 < len(block.route[i]) {
					fmt.Print(", ")
				}
			}
			fmt.Println("")
			fmt.Println("")
		}
	}

	fmt.Println("")
	for k := range propagation {
		fmt.Println(strconv.Itoa(k) + "," + strconv.FormatInt(propagation[k], 10))
	}
	fmt.Println("")
}

func PrintAllPropagation(numNodes int) {
	for i := 0; i < len(observedBlocks); i++ {
		printPropagation(observedBlocks[i], observedPropagations[i], numNodes)
	}
}

func IndexOf(haystack []*Block, needle *Block) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}
