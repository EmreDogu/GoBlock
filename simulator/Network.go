package simulator

import (
	"math"
	"math/rand"
	"os"
	settings "simblockgolang/settings"
	"strconv"
)

func GetRegionDistribution() []float64 {
	return settings.REGION_DISTRIBUTION_BITCOIN[:]
}

func GetDegreeDistribution() []float64 {
	return settings.DEGREE_DISTRIBUTION_BITCOIN[:]
}

func GetLatency(from int, to int) float64 {
	mean := float64(settings.LATENCY[from][to])
	shape := 0.2 * mean
	scale := mean - 5
	return math.Round(scale / math.Pow(rand.Float64(), 1/shape))
}

func GetBandwidth(from int, to int) int {
	if settings.UPLOAD_BANDWIDTH[from] <= settings.DOWNLOAD_BANDWIDTH[to] {
		return settings.UPLOAD_BANDWIDTH[from]
	} else {
		return settings.DOWNLOAD_BANDWIDTH[to]
	}
}

func PrintRegion() {
	f, err := os.Create("static.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{\"region\":[")

	if err2 != nil {
		panic(err2)
	}

	id := 0
	for ; id < len(settings.RegionList)-1; id++ {
		_, err3 := f.WriteString("{" + "\"id\":" + strconv.Itoa(id) + "," + "\"name\":\"" + settings.RegionList[id] + "\"" + "},")

		if err3 != nil {
			panic(err3)
		}
	}

	_, err3 := f.WriteString("{" + "\"id\":" + strconv.Itoa(id) + "," + "\"name\":\"" + settings.RegionList[id] + "\"" + "}" + "]}")

	if err3 != nil {
		panic(err3)
	}
}
