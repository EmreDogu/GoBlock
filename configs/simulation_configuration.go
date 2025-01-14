package configs

var NUM_OF_NODES int = 1000
var NUM_OF_BLOCKS int = 5
var NUM_OF_CONNECTIONS int = 2
var NUM_OF_HIGH_BANDWIDTH_CONNECTIONS int = 0
var END_BLOCK_HEIGHT int = 5
var CBR_USAGE_RATE float64 = 0.964
var BLOCK_SIZE int = 535000
var COMPACT_BLOCK_SIZE = 18 * 1000
var currentTime int64

func GetCurrentTime() int64 {
	return currentTime
}

func SetCurrentTime(time int64) {
	currentTime = time
}
