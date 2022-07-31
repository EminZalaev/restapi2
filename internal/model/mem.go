package model

type BenchMemMetrics struct {
	Metrics []MemMetrics `json:"metrics"`
}

type MemMetrics struct {
	Date                 string `json:"date"`
	OffHeap              int    `json:"offHeap"`
	InHeap               int    `json:"inHeap"`
	InStack              int    `json:"inStack"`
	TotalUsedMemory      int    `json:"totalUsedMemory"`
	AllocationRates      int    `json:"allocationRates"`
	NumberOfLiveObjects  int    `json:"numberOfLiveObjects"`
	RateObjectsAllocated int    `json:"rateObjectsAllocated"`
	Goroutines           int    `json:"goroutines"`
}
