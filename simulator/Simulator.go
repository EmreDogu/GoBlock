package simulator

import (
	"fmt"
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

func arriveBlock(block *Block, this *Node) {
	if ContainsBlock(observedBlocks, block) {
		propagation := observedPropagations[IndexOf(observedBlocks, block)]
		propagation[this.nodeID] = GetCurrentTime() - block.time
	} else {
		if len(observedBlocks) > 10 {
			printPropagation(observedBlocks[0], observedPropagations[0])
			observedBlocks = observedBlocks[1:]
			observedPropagations = observedPropagations[1:]
		}

		propagation := make(map[int]int64)
		propagation[this.nodeID] = GetCurrentTime() - block.time
		observedBlocks = append(observedBlocks, block)
		observedPropagations = append(observedPropagations, propagation)
	}
}

func printPropagation(block *Block, propagation map[int]int64) {
	fmt.Print(block)
	fmt.Print(":" + strconv.Itoa(block.height))
	for k := range propagation {
		fmt.Println(strconv.Itoa(k) + "," + strconv.FormatInt(propagation[k], 10))
	}
	fmt.Println("")
}

func PrintAllPropagation() {
	for i := 0; i < len(observedBlocks); i++ {
		printPropagation(observedBlocks[i], observedPropagations[i])
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
