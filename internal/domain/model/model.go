package model

//все результаты с таблицы
type AllBenchmarks struct {
	AllMetrics []AllBenchmarkMetrics
}

type AllBenchmarkMetrics struct {
	Date    string
	Metrics []Metrics
}

//все результаты по конкретной дате
type FullResult struct {
	Date    string
	Metrics []*Metrics
}

//ns/op по конкретной дате
type BenchStat struct {
	Date string
	Stat []*Stat
}

type Stat struct {
	Query   string
	NsPerOp int
}

//метрики с базы данных по конкретному запросу
type Metrics struct {
	Query                string
	NsPerOp              int
	OffHeap              int
	InStack              int
	InHeap               int
	TotalUsedMemory      int
	AllocationRates      int
	NumberOfLiveObjects  int
	RateObjectsAllocated int
	Goroutines           int
}

//delta метрик между двумя датами
type FullDelta struct {
	Delta []DeltaMetrics
}

type DeltaMetrics struct {
	Query                     string
	NsPerOp                   string
	OffHeapDelta              string
	InHeapDelta               string
	InStackDelta              string
	TotalUsedMemoryDelta      string
	AllocationRatesDelta      string
	NumberOfLiveObjectsDelta  string
	RateObjectsAllocatedDelta string
	GoroutinesDelta           string
}
