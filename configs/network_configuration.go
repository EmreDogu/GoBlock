package configs

import (
	"math"
	"math/rand"
)

var RegionList = [...]string{"NORTH_AMERICA", "EUROPE", "SOUTH_AMERICA", "ASIA_PACIFIC", "JAPAN",
	"AUSTRALIA"}

var REGION_DISTRIBUTION_BITCOIN = [...]float64{0.3316, 0.4998, 0.0090,
	0.1177, 0.0224, 0.0195}

var LATENCY = [6][6]int{
	{32, 124, 184, 198, 151, 189},
	{124, 11, 227, 237, 252, 294},
	{184, 227, 88, 325, 301, 322},
	{198, 237, 325, 85, 58, 198},
	{151, 252, 301, 58, 12, 126},
	{189, 294, 322, 198, 126, 16}}

var DOWNLOAD_BANDWIDTH = [7]int{
	52000000, 40000000, 18000000, 22800000,
	22800000, 29900000, 6 * 1000000}

var UPLOAD_BANDWIDTH = [7]int{
	19200000, 20700000, 5800000, 15700000,
	10200000, 11300000, 6 * 1000000}

func GetLatency(from int, to int) int64 {
	mean := LATENCY[from][to]
	shape := 0.2 * float64(mean)
	scale := mean - 5
	return int64(math.Round(float64(scale) / math.Pow(rand.Float64(), 1/shape)))
}

func GetBandwidth(from int, to int) int {
	if UPLOAD_BANDWIDTH[from] <= DOWNLOAD_BANDWIDTH[to] {
		return UPLOAD_BANDWIDTH[from]
	} else {
		return DOWNLOAD_BANDWIDTH[to]
	}
}
