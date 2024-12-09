package simulator

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"

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

	var regionDistribution []float64 = GetRegionDistribution()
	var regionList []int = makeRandomListFollowDistribution(regionDistribution, false)

	var useCBRNodes []bool = makeRandomList(cbrUsage, numNodes)

	file, err := os.Open("data/input/1000_2.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new scanner for the file
	scanner := bufio.NewScanner(file)

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()
		info := strings.Split(line, ",")

		i, err := strconv.Atoi(info[0])
		if err != nil {
			// ... handle error
			panic(err)
		}

		j, err := strconv.Atoi(info[1])
		if err != nil {
			// ... handle error
			panic(err)
		}

		node := node.New(i+1, numCon, j, numHighCon, useCBRNodes[i])
		AddNode(node)

		f, err := os.OpenFile("data/output/output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		_, err2 := f.WriteString("{" + "\"kind\":\"add-node\"," + "\"content\":{" + "\"timestamp\":0," + "\"node-id\":" + strconv.Itoa(i+1) + "," + "\"region-id\":" + strconv.Itoa(regionList[i]) + "}" + "},")
		if err2 != nil {
			panic(err2)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		panic(err)
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
			//printPropagation(observedBlocks[0], observedPropagations[0])
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

func Statistics() {
	f1, err := os.Create("data/output/output.txt")

	if err != nil {
		panic(err)
	}

	defer f1.Close()

	var count int
	var meanBlock int
	var medianBlock int
	percentilePropagation := []int{0, 0, 0, 0, 0, 0, 0, 0, 0}
	percentages := []float64{10, 20, 30, 40, 50, 60, 70, 80, 90}
	var combinedPropagation []string
	for _, propagationMap := range observedPropagations {
		values := make([]int64, 0, len(propagationMap))
		for _, v := range propagationMap {
			if v > 0 {
				values = append(values, v)
			}
		}

		sort.Slice(values, func(i, j int) bool {
			return values[i] < values[j]
		})

		if len(values)%2 != 0 {
			medianBlock += int(values[len(values)/2])
		} else {
			medianBlock += (int(values[(len(values)-1)/2]) + int(values[len(values)/2])) / 2
		}

		for value := range values {
			count++
			meanBlock += int(values[value])
			combinedPropagation = append(combinedPropagation, (strconv.FormatInt(values[value], 10) + ","))
		}
		combinedPropagation = append(combinedPropagation, "\n")

		numElements := len(values)
		for key, percentage := range percentages {
			index := int(math.Ceil(float64(numElements) * (percentage / 100.0)))
			if index <= numElements {
				percentilePropagation[key] += int(values[index-1])
			}
		}
	}
	meanBlock = meanBlock / count
	medianBlock = medianBlock / len(observedPropagations)

	_, err2 := f1.WriteString("Mean Block Propagation Time = " + strconv.Itoa(meanBlock) + "\n")
	_, err3 := f1.WriteString("Median Block Propagation Time = " + strconv.Itoa(medianBlock) + "\n")

	if err2 != nil || err3 != nil {
		panic(err2)
	}

	for i := range percentilePropagation {
		percentilePropagation[i] = percentilePropagation[i] / len(observedPropagations)
		_, err2 := f1.WriteString(strconv.FormatFloat(percentages[i], 'f', -1, 64) + "%" + " percentile of Block Propagation Time = " + strconv.Itoa(percentilePropagation[i]) + "\n")

		if err2 != nil {
			panic(err2)
		}
	}

	_, err4 := f1.WriteString("Block-based Propagation Times = ")

	for _, str := range combinedPropagation {
		_, err := f1.WriteString(str)
		if err != nil {
			panic(err)
		}
	}

	if err4 != nil {
		panic(err4)
	}

	for i := range observedBlocks {
		_, err := f1.WriteString("\n" + strconv.Itoa(observedBlocks[i].GetID()) + " " + strconv.Itoa(observedBlocks[i].GetHeight()) + " " + strconv.Itoa(int(observedBlocks[i].GetTime())))
		if err != nil {
			panic(err)
		}
	}
}
