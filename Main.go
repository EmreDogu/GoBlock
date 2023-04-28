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

var NUM_OF_BLOCKS int

var CON_ALG int

var NODE_LIST int

var NUM_OF_CON int

var simulationTime int64 = 0

func main() {

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

	fmt.Print("Bu simülasyonda kaç tane düğüm oluşturulacağını girin: ")
	fmt.Scan(&settings.NUM_OF_NODES)
	fmt.Println("")

	fmt.Print("Bu simülasyonda kaç tane blok kazılacağını girin: ")
	fmt.Scan(&NUM_OF_BLOCKS)
	fmt.Println("")

	fmt.Print("Düğümlerin nasıl oluşturulacağını girin (1-rastgele, 2-özel, 3-bitnodes): ")
	fmt.Scan(&NODE_LIST)
	fmt.Println("")

	if NODE_LIST != 1 {
		fmt.Print("Düğümlerin hangi eş seçim yolu/algoritması ile bağlanacağını girin (1-özel, 2-randompair, 3-nearpair, 4-clusterpair, 5-halfpair, 6-nmfpair, 7-TwoContinentBCBSN, 8-BCBSN): ")
		fmt.Scan(&CON_ALG)
		fmt.Println("")
	} else {
		CON_ALG = 2
	}

	fmt.Print("Düğümlerin kaç tane dışarıya giden bağlantıya sahip olabileceğini girin: ")
	fmt.Scan(&NUM_OF_CON)
	fmt.Println("")

	fmt.Print("Yakın komşu seçimi (proximity neighbor selection) algoritması aktifleştirilsin mi? (E-H): ")
	fmt.Scan(&settings.NEIGH_SEL)
	fmt.Println("")

	if settings.NEIGH_SEL == "E" || settings.NEIGH_SEL == "e" {
		settings.Matrix = make([][]int64, settings.NUM_OF_NODES)
		for i := 0; i < settings.NUM_OF_NODES; i++ {
			settings.Matrix[i] = make([]int64, settings.NUM_OF_NODES)
		}
	}

	start := time.Now().UnixMilli()
	simulator.SetTargetInterval(settings.INTERVAL)

	constructNetworkWithAllNodes(settings.NUM_OF_NODES, NUM_OF_CON, NODE_LIST)

	currentBlockHeight := 1

	for simulator.Pq.Len() > 0 {
		if isTest(simulator.GetTask()) {
			task := simulator.GetMintingTask()
			if task.GetParent().GetHeight() == currentBlockHeight {
				currentBlockHeight++
			}
			if currentBlockHeight > NUM_OF_BLOCKS {
				break
			}
			if currentBlockHeight%4 == 0 && (settings.NEIGH_SEL == "E" || settings.NEIGH_SEL == "e") {
				for i := range simulator.GetSimulatedNodes() {
					simulator.GetSimulatedNodes()[i].GetRoutingTable().ReconnectAll(settings.Matrix, settings.NUM_OF_NODES)
				}
			}
			/*// Log every 100 blocks and at the second block
			// TODO use constants here
			if currentBlockHeight%100 == 0 || currentBlockHeight == 2 {
				writeGraph(currentBlockHeight)
			}*/
		}
		simulator.RunTask()
	}

	simulator.PrintAllPropagation(settings.NUM_OF_NODES)

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
		for b := range a.GetOrphans() {
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

func genMiningPower() int {
	r := 0 + rand.Float64()*(1-0)
	return int(math.Max(r*float64(settings.STDEV_OF_MINING_POWER)+float64(settings.AVERAGE_MINING_POWER), 1))
}

func constructNetworkWithAllNodes(numNodes int, numCon int, nodeList int) {
	id := 0

	if nodeList == 1 {
		var regionDistribution []float64 = simulator.GetRegionDistribution()
		var regionList []int = makeRandomListFollowDistribution(regionDistribution, false)

		for id = 0; id < numNodes-1; id++ {
			node := simulator.MakeNode(id, numCon, strconv.Itoa(id), regionList[id], 0.0, 0.0, settings.RegionList[regionList[id]], genMiningPower(), settings.DOWNLOAD_BANDWIDTH[regionList[id]], settings.UPLOAD_BANDWIDTH[regionList[id]])
			simulator.AddNode(node)

			f, err := os.OpenFile("output.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}

			defer f.Close()

			_, err2 := f.WriteString("{" + "\"kind\":\"add-node\"," + "\"content\":{" + "\"timestamp\":0," + "\"node-id\":" + strconv.Itoa(id) + "," + "\"region-id\":" + strconv.Itoa(regionList[id]) + "}" + "},")
			if err2 != nil {
				panic(err2)
			}
		}
	} else {
		// json parselama gelecek
	}

	nodes := simulator.GetSimulatedNodes()

	if CON_ALG < 7 && CON_ALG > 0 {
		for i := range nodes {
			nodes[i].JoinNetwork(CON_ALG)
		}
	} else if CON_ALG == 7 || CON_ALG == 8 {
		nodes[0].JoinNetworkBCBSN(CON_ALG, nodes)
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
