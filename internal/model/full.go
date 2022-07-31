package model

type BenchFullResult struct {
	Metrics []struct {
		Date    string `json:"date"`
		NsPerOp struct {
			Catalog string `json:"catalog"`
			Filters string `json:"filters"`
		} `json:"ns_per_op"`
		OffHeap              int `json:"offHeap"`
		InHeap               int `json:"inHeap"`
		InStack              int `json:"inStack"`
		TotalUsedMemory      int `json:"totalUsedMemory"`
		AllocationRates      int `json:"allocationRates"`
		NumberOfLiveObjects  int `json:"numberOfLiveObjects"`
		RateObjectsAllocated int `json:"rateObjectsAllocated"`
		Goroutines           int `json:"goroutines"`
	} `json:"metrics"`
}
