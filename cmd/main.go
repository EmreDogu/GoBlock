package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/EmreDogu/GoBlock/cmd/simulator"
	"github.com/EmreDogu/GoBlock/configs"
	"github.com/EmreDogu/GoBlock/internal/blockchain/block"
)

func main() {
	var simulationTime int64 = 0

	start := time.Now().UnixMilli()

	s := simulator.Simulator{}
	s.InitializeSimulatorLink()

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

	simulator.ConstructNetworkWithAllNodes(configs.NUM_OF_NODES, configs.NUM_OF_CONNECTIONS, configs.NUM_OF_HIGH_BANDWIDTH_CONNECTIONS, configs.CBR_USAGE_RATE)

	simulator.Simulation(configs.END_BLOCK_HEIGHT)

	simulator.PrintAllPropagation()

	fmt.Println("")

	blocks := make([]*block.Block, 0)
	orphans := make([]*block.Block, 0)
	blockList := make([]*block.Block, 0)

	block := simulator.GetSimulatedNodes()[0].GetBlock()

	for block.GetParent() != nil {
		blocks = append(blocks, block)
		block = block.GetParent()
	}

	averageOrphansSize := 0
	for _, a := range simulator.GetSimulatedNodes() {
		for b := range a.GetOrphans() {
			orphans = append(orphans, b)
			averageOrphansSize += len(a.GetOrphans())
		}
	}
	averageOrphansSize = averageOrphansSize / len(simulator.GetSimulatedNodes())

	blocks = append(blocks, orphans...)

	blockList = append(blockList, blocks...)

	sort.Slice(blockList, func(i, j int) bool {
		return blockList[i].GetTime() < blockList[j].GetTime()
	})

	for _, orphan := range orphans {
		fmt.Printf("%v:%d\n", orphan, orphan.GetHeight())
	}
	fmt.Println(averageOrphansSize)

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

	_, err3 := f3.WriteString("{" + "\"kind\":\"simulation-end\"," + "\"content\":{" + "\"timestamp\":" + strconv.FormatInt(configs.GetCurrentTime(), 10) + "}" + "}" + "]")

	if err3 != nil {
		panic(err3)
	}

	end := time.Now().UnixMilli()
	simulationTime += end - start
	fmt.Println("Simulation time: " + strconv.FormatInt(simulationTime, 10))
}
