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
	result, err := s.store.GetMetrics()
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
	result, err := s.store.GetMetricsByDate(string(date))
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

	if len(dateArr) < 2 {
		return nil, fmt.Errorf("not enough argument in query")
	}

	benchFirstMetrics, err := store.GetMetricsByDate(dateArr[0])
	if err != nil {
		return nil, err
	}

	benchSecondMetrics, err := store.GetMetricsByDate(dateArr[1])
	if err != nil {
		return nil, err
	}

	delta := compareMetrics(benchFirstMetrics, benchSecondMetrics)

	return delta, nil
}

//сравнение метрик между двумя датами
func compareMetrics(firstDateMetrics, secondDateMetrics *model.FullResult) *model.FullDelta {
	var d model.FullDelta

	for _, el := range firstDateMetrics.Metrics {
		for _, el2 := range secondDateMetrics.Metrics {
			if el.Query == el2.Query {
				d.Delta = append(d.Delta, model.DeltaMetrics{
					Query:                     el.Query,
					NsPerOp:                   getDeltaPercentage(el.NsPerOp, el2.NsPerOp),
					OffHeapDelta:              getDeltaPercentage(el.OffHeap, el2.OffHeap),
					InHeapDelta:               getDeltaPercentage(el.InHeap, el2.InHeap),
					InStackDelta:              getDeltaPercentage(el.InStack, el2.InStack),
					TotalUsedMemoryDelta:      getDeltaPercentage(el.TotalUsedMemory, el2.TotalUsedMemory),
					AllocationRatesDelta:      getDeltaPercentage(el.AllocationRates, el2.AllocationRates),
					NumberOfLiveObjectsDelta:  getDeltaPercentage(el.NumberOfLiveObjects, el2.NumberOfLiveObjects),
					RateObjectsAllocatedDelta: getDeltaPercentage(el.RateObjectsAllocated, el2.RateObjectsAllocated),
					GoroutinesDelta:           getDeltaPercentage(el.Goroutines, el2.Goroutines),
				})
			}
		}
	}

	return &d
}

func getDeltaPercentage(first, second int) string {
	val := strconv.Itoa(first)

	diff := float64(second-first) / float64(first)
	if diff > 0 {
		return fmt.Sprintf("%s(-%.2f)", val, diff*-100)
	}

	return fmt.Sprintf("%s(+%.2f)", val, diff*100)
}
