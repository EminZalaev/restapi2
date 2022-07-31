package service

import (
	"fmt"
	"log"
	"restapi2/internal/model"
	"strconv"
	"strings"
)

const dataFileName = "./data.json"

func CountDifference(indexes string, cache *CacheBench) (*model.BenchDelta, error) {
	indexes = strings.Trim(indexes, "[]")
	indexArr := strings.Split(indexes, ",")
	log.Print(indexArr)

	var benchmarkStatGet model.BenchMemMetrics
	benchMemStat, err := getBenchMem(cache)
	if err != nil {
		return nil, err
	}
	for i, el := range benchMemStat.Metrics {
		for _, el1 := range indexArr {
			j, err := strconv.Atoi(el1)
			if err != nil {
				return nil, fmt.Errorf("error parse data: %w", err)
			}
			if i == j {
				benchmarkStatGet.Metrics = append(benchmarkStatGet.Metrics, el)
			}
		}
	}
	delta, err := compareMetrics(&benchmarkStatGet)
	if err != nil {
		return nil, err
	}
	return delta, nil
}

func compareMetrics(m *model.BenchMemMetrics) (*model.BenchDelta, error) {
	var d model.BenchDelta
	if len(m.Metrics) < 2 {
		return nil, fmt.Errorf("not enogh data")
	}
	//compare current with previous
	for i := 0; i <= len(m.Metrics)-1; i++ {
		if i == len(m.Metrics)-1 {
			d.Delta = append(d.Delta, model.Delta{
				Date:                      m.Metrics[i].Date,
				OffHeapDelta:              getDeltaPercentage(m.Metrics[i].OffHeap, m.Metrics[i].OffHeap),
				InHeapDelta:               getDeltaPercentage(m.Metrics[i].InHeap, m.Metrics[i].InHeap),
				InStackDelta:              getDeltaPercentage(m.Metrics[i].InStack, m.Metrics[i].InStack),
				TotalUsedMemoryDelta:      getDeltaPercentage(m.Metrics[i].TotalUsedMemory, m.Metrics[i].TotalUsedMemory),
				AllocationRatesDelta:      getDeltaPercentage(m.Metrics[i].AllocationRates, m.Metrics[i].AllocationRates),
				NumberOfLiveObjectsDelta:  getDeltaPercentage(m.Metrics[i].NumberOfLiveObjects, m.Metrics[i].NumberOfLiveObjects),
				RateObjectsAllocatedDelta: getDeltaPercentage(m.Metrics[i].RateObjectsAllocated, m.Metrics[i].RateObjectsAllocated),
				GoroutinesDelta:           getDeltaPercentage(m.Metrics[i].Goroutines, m.Metrics[i].Goroutines),
			})
			continue
		}
		d.Delta = append(d.Delta, model.Delta{
			Date:                      m.Metrics[i].Date,
			OffHeapDelta:              getDeltaPercentage(m.Metrics[i].OffHeap, m.Metrics[i+1].OffHeap),
			InHeapDelta:               getDeltaPercentage(m.Metrics[i].InHeap, m.Metrics[i+1].InHeap),
			InStackDelta:              getDeltaPercentage(m.Metrics[i].InStack, m.Metrics[i+1].InStack),
			TotalUsedMemoryDelta:      getDeltaPercentage(m.Metrics[i].TotalUsedMemory, m.Metrics[i+1].TotalUsedMemory),
			AllocationRatesDelta:      getDeltaPercentage(m.Metrics[i].AllocationRates, m.Metrics[i+1].AllocationRates),
			NumberOfLiveObjectsDelta:  getDeltaPercentage(m.Metrics[i].NumberOfLiveObjects, m.Metrics[i+1].NumberOfLiveObjects),
			RateObjectsAllocatedDelta: getDeltaPercentage(m.Metrics[i].RateObjectsAllocated, m.Metrics[i+1].RateObjectsAllocated),
			GoroutinesDelta:           getDeltaPercentage(m.Metrics[i].Goroutines, m.Metrics[i+1].Goroutines),
		})

	}

	return &d, nil
}

func getDeltaPercentage(first, second int) string {
	val := strconv.Itoa(first)
	d := float64(second-first) / float64(first)
	if d > 0 {
		return fmt.Sprintf("%s(-%.2f)", val, d*-100)
	}
	return fmt.Sprintf("%s(+%.2f)", val, d*100)
}
