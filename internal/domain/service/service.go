package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"webserver/internal/domain/model"

	"webserver/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

type Service struct {
	store *repository.Store
}

func NewService(db *repository.Store) *Service {
	return &Service{
		store: db,
	}
}

func (s *Service) GetAllResult() ([]byte, error) {
	result, err := s.store.GetAllMetrics()
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (s *Service) CompareResults(dateByte []byte) ([]byte, error) {
	result, err := countDifference(s.store, dateByte)
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (s *Service) GetFullResultByDateAndQuery(date []byte) ([]byte, error) {
	result, err := s.store.GetFullByQuery(string(date))
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (s *Service) GetNsPerOp(dateByte []byte) ([]byte, error) {
	result, err := s.store.GetNsPerOpByDate(string(dateByte))
	if err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func countDifference(store *repository.Store, dateByte []byte) (*model.FullDelta, error) {
	date := strings.Trim(string(dateByte), "[]")
	dateArr := strings.Split(date, ";")

	allBenchmarksArr := make([]model.AllBenchmarks, 0)
	benchFirstMetrics, err := store.GetMetricsByDate(dateArr[0])
	if err != nil {
		return nil, err
	}
	benchSecondMetrics, err := store.GetMetricsByDate(dateArr[1])
	if err != nil {
		return nil, err
	}
	allBenchmarksArr = append(allBenchmarksArr, *benchFirstMetrics, *benchSecondMetrics)
	delta := compareMetrics(allBenchmarksArr)
	return delta, nil
}

//сравнение метрик между двумя датами
func compareMetrics(m []model.AllBenchmarks) *model.FullDelta {
	var d model.FullDelta

	for i, el := range m[0].AllMetrics {
		for j, el2 := range m[1].AllMetrics {
			if el.Metrics[i].Query == el2.Metrics[j].Query {
				d.Delta = append(d.Delta, model.DeltaMetrics{
					Query:                     el.Metrics[i].Query,
					NsPerOp:                   getDeltaPercentage(el.Metrics[i].NsPerOp, el2.Metrics[j].NsPerOp),
					OffHeapDelta:              getDeltaPercentage(el.Metrics[i].OffHeap, el2.Metrics[j].OffHeap),
					InHeapDelta:               getDeltaPercentage(el.Metrics[i].InHeap, el2.Metrics[j].InHeap),
					InStackDelta:              getDeltaPercentage(el.Metrics[i].InStack, el2.Metrics[j].InStack),
					TotalUsedMemoryDelta:      getDeltaPercentage(el.Metrics[i].TotalUsedMemory, el2.Metrics[j].TotalUsedMemory),
					AllocationRatesDelta:      getDeltaPercentage(el.Metrics[i].AllocationRates, el2.Metrics[j].AllocationRates),
					NumberOfLiveObjectsDelta:  getDeltaPercentage(el.Metrics[i].NumberOfLiveObjects, el2.Metrics[j].NumberOfLiveObjects),
					RateObjectsAllocatedDelta: getDeltaPercentage(el.Metrics[i].RateObjectsAllocated, el2.Metrics[j].RateObjectsAllocated),
					GoroutinesDelta:           getDeltaPercentage(el.Metrics[i].Goroutines, el2.Metrics[j].Goroutines),
				})
			}
		}
	}
	return &d

}

func getDeltaPercentage(first, second int) string {
	val := strconv.Itoa(first)
	d := float64(second-first) / float64(first)
	if d > 0 {
		return fmt.Sprintf("%s(-%.2f)", val, d*-100)
	}
	return fmt.Sprintf("%s(+%.2f)", val, d*100)
}
