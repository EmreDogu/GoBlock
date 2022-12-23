package simulator

import (
	"fmt"
	"strconv"
)

var TargetInterval int
var simulatedNodes []*Node
var observedBlocks []*Block
var observedPropagations = make(map[int]int64)

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
		propagation := make(map[int]int64)
		propagation[this.nodeID] = GetCurrentTime() - block.time
	} else {
		if len(observedBlocks) > 10 {
			keys := make([]int, len(observedPropagations))
			i := 0
			for k := range observedPropagations {
				keys[i] = k
				i++
			}
			printPropagation(observedBlocks[0], keys[0], observedPropagations[keys[0]])
			observedBlocks = observedBlocks[1:]
			delete(observedPropagations, keys[0])
		}
	}

	propagation := make(map[int]int64)
	propagation[this.nodeID] = GetCurrentTime() - block.time
	observedBlocks = append(observedBlocks, block)
	observedPropagations[this.nodeID] = GetCurrentTime() - block.time
}

func printPropagation(block *Block, id int, time int64) {
	fmt.Print(block)
	fmt.Print(":" + strconv.Itoa(block.height) + strconv.Itoa(id) + "," + strconv.FormatInt(time, 10))
	fmt.Println("")
}

func PrintAllPropagation() {
	keys := make([]int, len(observedPropagations))
	i := 0
	for k := range observedPropagations {
		keys[i] = k
		i++
	}
	for i := 0; i < len(observedPropagations); i++ {
		printPropagation(observedBlocks[i], keys[i], observedPropagations[keys[i]])
	}
}
