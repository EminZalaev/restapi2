package model

type BenchDelta struct {
	Delta []Delta
}

type Delta struct {
	Date                      string
	OffHeapDelta              string
	InHeapDelta               string
	InStackDelta              string
	TotalUsedMemoryDelta      string
	AllocationRatesDelta      string
	NumberOfLiveObjectsDelta  string
	RateObjectsAllocatedDelta string
	GoroutinesDelta           string
}
