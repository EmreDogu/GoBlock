package simulator

import (
	"os"
	"strconv"

	"github.com/EmreDogu/GoBlock/configs"
)

func GetRegionDistribution() []float64 {
	return configs.REGION_DISTRIBUTION_BITCOIN[:]
}

func PrintRegion() {
	f, err := os.Create("data/output/static.json")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("{\"region\":[")

	if err2 != nil {
		panic(err2)
	}

	id := 0
	for ; id < len(configs.RegionList)-1; id++ {
		_, err3 := f.WriteString("{" + "\"id\":" + strconv.Itoa(id) + "," + "\"name\":\"" + configs.RegionList[id] + "\"" + "},")

		if err3 != nil {
			panic(err3)
		}
	}

	_, err3 := f.WriteString("{" + "\"id\":" + strconv.Itoa(id) + "," + "\"name\":\"" + configs.RegionList[id] + "\"" + "}" + "]}")

	if err3 != nil {
		panic(err3)
	}
}
