package repository

import (
	"database/sql"
	"fmt"

	"webserver/internal/domain/model"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetMetrics() ([]model.FullResult, error) {
	result := make([]model.Metrics, 0)
	mtrc := &model.Metrics{}

	var date, oldDate string
	allRes := make([]model.FullResult, 0)
	query := "SELECT d.benchdate,queries,offheap,instack,totalusedmemory,allocationrates,numberofliveobjects,rateobjectsallocated,goroutines,inheap,nsperop from memorymetrics m, date d where m.benchdate = d.id"
	dbRows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		if err = dbRows.Scan(&date,
			&mtrc.Query,
			&mtrc.OffHeap,
			&mtrc.InStack,
			&mtrc.TotalUsedMemory,
			&mtrc.AllocationRates,
			&mtrc.NumberOfLiveObjects,
			&mtrc.RateObjectsAllocated,
			&mtrc.Goroutines,
			&mtrc.InHeap,
			&mtrc.NsPerOp); err != nil {
			return nil, err
		}

		result, err = addMetrics(*mtrc, result)
		if err != nil {
			return nil, err
		}

		if date != oldDate {
			if len(result) > 1 {
				allRes = append(allRes, model.FullResult{
					Date:    oldDate,
					Metrics: result[:len(result)-1],
				})
				result = result[len(result)-1:]
			}
			oldDate = date
		}

	}

	allRes = append(allRes, model.FullResult{
		Date:    oldDate,
		Metrics: result,
	})

	return allRes, nil
}

func (s *Store) GetMetricsByDate(date string) (*model.FullResult, error) {
	result := make([]model.Metrics, 0)
	mtrc := &model.Metrics{}
	query := "SELECT queries,offheap,instack,totalusedmemory,allocationrates,numberofliveobjects,rateobjectsallocated,goroutines,inheap,nsperop from memorymetrics m, date d where m.benchdate = d.id and d.benchdate=$1"
	dbRows, err := s.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		if err := dbRows.Scan(&mtrc.Query,
			&mtrc.OffHeap,
			&mtrc.InStack,
			&mtrc.TotalUsedMemory,
			&mtrc.AllocationRates,
			&mtrc.NumberOfLiveObjects,
			&mtrc.RateObjectsAllocated,
			&mtrc.Goroutines,
			&mtrc.InHeap,
			&mtrc.NsPerOp); err != nil {
			return nil, err
		}

		result, err = addMetrics(*mtrc, result)
		if err != nil {
			return nil, err
		}
	}
	if mtrc.Query == "" {
		return nil, fmt.Errorf("error wrong input data")
	}

	return &model.FullResult{
		Date:    date,
		Metrics: result,
	}, nil
}

func (s *Store) GetNsPerOpByDate(date string) (*model.BenchStat, error) {
	result := make([]*model.Stat, 0)
	mtrc := &model.Stat{}
	query := "SELECT queries,nsperop from memorymetrics m, date d where m.benchdate = d.id and d.benchdate=$1"
	dbRows, err := s.db.Query(query, date)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		err = dbRows.Scan(
			&mtrc.Query,
			&mtrc.NsPerOp)
		if err != nil {
			return nil, fmt.Errorf("error scan from db")
		}
		result = append(result, &model.Stat{
			Query:   mtrc.Query,
			NsPerOp: mtrc.NsPerOp,
		})
	}

	if mtrc.Query == "" {
		return nil, fmt.Errorf("error wrong input data")
	}

	return &model.BenchStat{
		Date: date,
		Stat: result,
	}, nil
}

func addMetrics(mtrc model.Metrics, result []model.Metrics) ([]model.Metrics, error) {
	result = append(result, model.Metrics{
		Query:                mtrc.Query,
		OffHeap:              mtrc.OffHeap,
		InStack:              mtrc.InStack,
		TotalUsedMemory:      mtrc.TotalUsedMemory,
		AllocationRates:      mtrc.AllocationRates,
		NumberOfLiveObjects:  mtrc.NumberOfLiveObjects,
		RateObjectsAllocated: mtrc.RateObjectsAllocated,
		Goroutines:           mtrc.Goroutines,
		InHeap:               mtrc.InHeap,
		NsPerOp:              mtrc.NsPerOp,
	})
	return result, nil
}
