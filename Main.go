package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"

	settings "simblockgolang/settings"
	simulator "simblockgolang/simulator"
)

var simulationTime int64 = 0

func main() {
	start := time.Now().UnixMilli()
	simulator.SetTargetInterval(settings.INTERVAL)

	f, err := os.Create("output.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("[")

	if err2 != nil {
		panic(err2)
	}

	simulator.PrintRegion()

	constructNetworkWithAllNodes(settings.NUM_OF_NODES)

	currentBlockHeight := 1

	for simulator.Pq.Len() > 0 {
		if isTest(simulator.GetTask()) {
			task := simulator.GetMintingTask()
			if task.GetParent().GetHeight() == currentBlockHeight {
				currentBlockHeight++
			}
			if currentBlockHeight > settings.END_BLOCK_HEIGHT {
				break
			}
			/*// Log every 100 blocks and at the second block
			// TODO use constants here
			if currentBlockHeight%100 == 0 || currentBlockHeight == 2 {
				writeGraph(currentBlockHeight)
			}*/
		}
		simulator.RunTask()
	}

	simulator.PrintAllPropagation()

	fmt.Println("")

	blocks := make([]*simulator.Block, 0)
	block := simulator.GetSimulatedNodes()[0].GetBlock()

	for !reflect.ValueOf(block.GetParent()).IsNil() {
		blocks = append(blocks, block)
		block = block.GetParent()
	}

	orphans := make([]*simulator.Block, 0)
	averageOrphansSize := 0

	for _, a := range simulator.GetSimulatedNodes() {
		for _, b := range a.GetOrphans() {
			orphans = append(orphans, b)
			averageOrphansSize += len(a.GetOrphans())
		}
	}

	averageOrphansSize = averageOrphansSize / len(simulator.GetSimulatedNodes())

	blocks = append(blocks, orphans...)

	sort.SliceStable(blocks, func(i, j int) bool {
		return blocks[i].GetTime() < blocks[j].GetTime()
	})

	for _, a := range orphans {
		fmt.Print(a)
		fmt.Println(":" + strconv.Itoa(a.GetHeight()))
	}

	fmt.Println("")
	fmt.Println("Average orphans size: " + strconv.Itoa(averageOrphansSize))

	f2, err := os.Create("blockList.txt")

	if err != nil {
		panic(err)
	}

	defer f2.Close()

	for _, b := range blocks {
		if !simulator.ContainsBlock(orphans, b) {

			_, err2 := f2.WriteString("OnChain : " + strconv.Itoa(b.GetHeight()) + " : " + strconv.Itoa(b.GetID()) + " ")

			if err2 != nil {
				panic(err2)
			}
		} else {
			_, err2 := f2.WriteString("Orphan : " + strconv.Itoa(b.GetHeight()) + " : " + strconv.Itoa(b.GetID()) + " ")

			if err2 != nil {
				panic(err2)
			}
		}
	}

	f3, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f3.Close()

	_, err3 := f3.WriteString("{" + "\"kind\":\"simulation-end\"," + "\"content\":{" + "\"timestamp\":" + strconv.FormatInt(simulator.GetCurrentTime(), 10) + "}" + "}" + "]")

	if err3 != nil {
		panic(err3)
	}

	end := time.Now().UnixMilli()
	simulationTime = end - start
	fmt.Println("")
	fmt.Println("Simulation time: " + strconv.FormatInt(simulationTime, 10))
}

func makeRandomListFollowDistribution(distribution []float64, facum bool) []int {
	list := []int{}
	index := 0

	if facum {
		for ; index < len(distribution); index++ {
			for float64(len(list)) <= float64(settings.NUM_OF_NODES)*distribution[index] {
				list = append(list, index)
			}
		}
		for len(list) < settings.NUM_OF_NODES {
			list = append(list, index)
		}
	} else {
		var acumulative float64 = 0.0
		for ; index < len(distribution); index++ {
			acumulative += distribution[index]
			for float64(len(list)) <= float64(settings.NUM_OF_NODES)*acumulative {
				list = append(list, index)
			}
		}
		for len(list) < settings.NUM_OF_NODES {
			list = append(list, index)
		}
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list
}

func makeRandomList(rate float64) []bool {
	list := []bool{}

	for i := 0; i < settings.NUM_OF_NODES; i++ {
		list = append(list, (float64(i) < float64(settings.NUM_OF_NODES)*rate))
	}

	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	return list
}

func genMiningPower() int {
	r := rand.Float64()
	return int(math.Max(r*float64(settings.STDEV_OF_MINING_POWER)+float64(settings.AVERAGE_MINING_POWER), 1))
}

func constructNetworkWithAllNodes(numNodes int) {
	var regionDistribution []float64 = simulator.GetRegionDistribution()
	var regionList []int = makeRandomListFollowDistribution(regionDistribution, false)

	var degreeDistribution []float64 = simulator.GetDegreeDistribution()
	var degreeList []int = makeRandomListFollowDistribution(degreeDistribution, true)

	var useCBRNodes []bool = makeRandomList(settings.CBR_USAGE_RATE)

	var churnNodes []bool = makeRandomList(settings.CHURN_NODE_RATE)

	for id := 1; id <= numNodes; id++ {
		node := simulator.MakeNode(id, degreeList[id-1]+1, regionList[id-1], genMiningPower(), useCBRNodes[id-1], churnNodes[id-1])
		simulator.AddNode(node)

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

	nodes := simulator.GetSimulatedNodes()
	for i := range nodes {
		nodes[i].JoinNetwork()
	}

	nodes[0].GenesisBlock()
}

func isTest(t interface{}) bool {
	switch t.(type) {
	case *simulator.MintingTask:
		return true
	default:
		return false
	}
}
