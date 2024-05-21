package simulator

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
	"github.com/EmreDogu/GoBlock/internal/blockchain/node"
)

type Simulator struct{}

func New() *Simulator {
	return &Simulator{}
}

var simulatedNodes []*node.Node
var observedBlocks []*block.Block
var observedPropagations []map[int]int64
var TargetInterval int

var currentBlockHeight int = 1

func ContainsBlock(s []*block.Block, e *block.Block) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func AddNode(node *node.Node) {
	simulatedNodes = append(simulatedNodes, node)
}

func GetSimulatedNodes() []*node.Node {
	return simulatedNodes
}

func ConstructNetworkWithAllNodes(numNodes int, numCon int, numHighCon int, cbrUsage float64) {
	id := 0

	var regionDistribution []float64 = GetRegionDistribution()
	var regionList []int = makeRandomListFollowDistribution(regionDistribution, false)

	var useCBRNodes []bool = makeRandomList(cbrUsage, numNodes)

	for id = 1; id <= numNodes; id++ {
		node := node.New(id, numCon, regionList[id-1], numHighCon, useCBRNodes[id-1])
		AddNode(node)

		f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err2 := f.WriteString("{" + "\"kind\":\"add-node\"," + "\"content\":{" + "\"timestamp\":0," + "\"node-id\":" + strconv.Itoa(id) + "," + "\"region-id\":" + strconv.Itoa(regionList[id-1]) + "}" + "},")
		if err2 != nil {
			panic(err2)
		}
	}

	node.ReceiveSimulatedNodes(simulatedNodes)

	for i := range simulatedNodes {
		simulatedNodes[i].JoinNetwork()
	}

	simulatedNodes[0].GenesisBlock()
}

func makeRandomListFollowDistribution(distribution []float64, check bool) []int {
	list := []int{}
	index := 0

	if check {
		for ; index < len(distribution); index++ {
			for float64(len(list)) <= float64(configs.NUM_OF_NODES)*distribution[index] {
				list = append(list, index)
			}
		}
		for len(list) < configs.NUM_OF_NODES {
			list = append(list, index)
		}
	} else {
		var accumulative float64 = 0.0
		for ; index < len(distribution); index++ {
			accumulative += distribution[index]
			for float64(len(list)) <= float64(configs.NUM_OF_NODES)*accumulative {
				list = append(list, index)
			}
		}
		for len(list) < configs.NUM_OF_NODES {
			list = append(list, index)
		}
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list
}

func makeRandomList(rate float64, numNodes int) []bool {
	list := []bool{}

	for i := 0; i < numNodes; i++ {
		list = append(list, (float64(i) < float64(numNodes)*rate))
	}
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list
}

func (s *Simulator) InitializeSimulatorLink() {
	link := node.NewSimulatorLink(s)
	node.ReceiveSimulatorLink(link)
}

func (s *Simulator) ArriveBlock(newBlock *block.Block, node *node.Node) {
	if ContainsBlock(observedBlocks, newBlock) {
		propagation := observedPropagations[IndexOf(observedBlocks, newBlock)]
		propagation[node.GetID()] = configs.GetCurrentTime() - newBlock.GetTime()
	} else {
		if len(observedBlocks) > 10 {
			printPropagation(observedBlocks[0], observedPropagations[0])
			observedBlocks = observedBlocks[1:]
			observedPropagations = observedPropagations[1:]
		}

		propagation := make(map[int]int64)
		propagation[node.GetID()] = configs.GetCurrentTime() - newBlock.GetTime()
		observedBlocks = append(observedBlocks, newBlock)
		observedPropagations = append(observedPropagations, propagation)
	}
}

func IndexOf(haystack []*block.Block, needle *block.Block) int {
	for i, v := range haystack {
		if v == needle {
			return i
		}
	}
	return -1
}

func Simulation(END_BLOCK_HEIGHT int) {
	for GetTask() != nil {
		if GetTask().task.GetType() == "mining" {
			task := GetTask().task

			if task.GetParent().GetHeight() == currentBlockHeight {
				currentBlockHeight++
			}
			if currentBlockHeight > END_BLOCK_HEIGHT {
				break
			}
		}

		RunTask()
	}
}

func printPropagation(block *block.Block, propagation map[int]int64) {
	fmt.Printf("%+v:%d\n", block, block.GetHeight())
	for key, value := range propagation {
		fmt.Printf("%d,%d\n", key, value)
	}
	fmt.Println()
}

func PrintAllPropagation() {
	for i := 0; i < len(observedBlocks); i++ {
		printPropagation(observedBlocks[i], observedPropagations[i])
	}
}
