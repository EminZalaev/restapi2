package service

import "restapi2/internal/model"

func getBenchMem(cache *CacheBench) (*model.BenchMemMetrics, error) {
	mem, err := readBenchMem(cache)
	if err != nil {
		return nil, err
	}
	reverseMem(mem)
	return mem, nil
}

func readBenchMem(cache *CacheBench) (*model.BenchMemMetrics, error) {
	tmp := cache.Get(0)
	mtrc := make([]model.MemMetrics, 0, len(tmp.Metrics))
	for i := 0; i < len(tmp.Metrics); i++ {
		mtrc = append(mtrc, model.MemMetrics{
			Date:                 tmp.Metrics[i].Date,
			OffHeap:              tmp.Metrics[i].OffHeap,
			InHeap:               tmp.Metrics[i].InHeap,
			InStack:              tmp.Metrics[i].InStack,
			TotalUsedMemory:      tmp.Metrics[i].TotalUsedMemory,
			AllocationRates:      tmp.Metrics[i].AllocationRates,
			NumberOfLiveObjects:  tmp.Metrics[i].NumberOfLiveObjects,
			RateObjectsAllocated: tmp.Metrics[i].RateObjectsAllocated,
			Goroutines:           tmp.Metrics[i].Goroutines,
		})
	}

	return &model.BenchMemMetrics{Metrics: mtrc}, nil
}

func reverseMem(m *model.BenchMemMetrics) {
	for i, j := 0, len(m.Metrics)-1; i < j; i, j = i+1, j-1 {
		m.Metrics[i], m.Metrics[j] = m.Metrics[j], m.Metrics[i]
	}
}
