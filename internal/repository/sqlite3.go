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

func (s *Store) GetAllMetrics() (*model.AllBenchmarks, error) {
	result := make([]model.Metrics, 0)
	mtrc := &model.Metrics{}

	var date, oldDate string

	allRes := make([]model.AllBenchmarkMetrics, 0)
	dbRows, err := s.db.Query("SELECT d.benchdate,queries,offheap,instack,totalusedmemory,allocationrates,numberofliveobjects,rateobjectsallocated,goroutines,inheap,nsperop from memorymetrics m, date d where m.benchdate = d.id")
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		err = dbRows.Scan(&date,
			&mtrc.Query,
			&mtrc.OffHeap,
			&mtrc.InStack,
			&mtrc.TotalUsedMemory,
			&mtrc.AllocationRates,
			&mtrc.NumberOfLiveObjects,
			&mtrc.RateObjectsAllocated,
			&mtrc.Goroutines,
			&mtrc.InHeap,
			&mtrc.NsPerOp)
		if err != nil {
			return nil, err
		}
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

		if date != oldDate {
			if len(result) > 1 {
				allRes = append(allRes, model.AllBenchmarkMetrics{
					Date:    oldDate,
					Metrics: result[:len(result)-1],
				})
				result = result[len(result)-1:]
			}
			oldDate = date
		}

	}

	allRes = append(allRes, model.AllBenchmarkMetrics{
		Date:    oldDate,
		Metrics: result,
	})

	return &model.AllBenchmarks{
		AllMetrics: allRes,
	}, nil
}

func (s *Store) GetMetricsByDate(date string) (*model.AllBenchmarks, error) {
	result := make([]model.Metrics, 0)
	mtrc := &model.Metrics{}

	var oldDate string
	allRes := make([]model.AllBenchmarkMetrics, 0)
	dbRowsFirst, err := s.db.Query("SELECT queries,offheap,instack,totalusedmemory,allocationrates,numberofliveobjects,rateobjectsallocated,goroutines,inheap,nsperop from memorymetrics m, date d where m.benchdate = d.id and d.benchdate=$1", date)
	if err != nil {
		return nil, err
	}
	for dbRowsFirst.Next() {
		err = dbRowsFirst.Scan(
			&mtrc.Query,
			&mtrc.OffHeap,
			&mtrc.InStack,
			&mtrc.TotalUsedMemory,
			&mtrc.AllocationRates,
			&mtrc.NumberOfLiveObjects,
			&mtrc.RateObjectsAllocated,
			&mtrc.Goroutines,
			&mtrc.InHeap,
			&mtrc.NsPerOp)
		if err != nil {
			return nil, err
		}
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

		allRes = append(allRes, model.AllBenchmarkMetrics{
			Date:    oldDate,
			Metrics: result,
		})
	}
	return &model.AllBenchmarks{
		AllMetrics: allRes,
	}, nil
}

func (s *Store) GetFullByQuery(date string) (*model.FullResult, error) {
	result := make([]*model.Metrics, 0)
	mtrc := &model.Metrics{}
	dbRows, err := s.db.Query("SELECT queries,offheap,instack,totalusedmemory,allocationrates,numberofliveobjects,rateobjectsallocated,goroutines,inheap,nsperop from memorymetrics m, date d where m.benchdate = d.id and d.benchdate=$1", date)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		err = dbRows.Scan(&mtrc.Query,
			&mtrc.OffHeap,
			&mtrc.InStack,
			&mtrc.TotalUsedMemory,
			&mtrc.AllocationRates,
			&mtrc.NumberOfLiveObjects,
			&mtrc.RateObjectsAllocated,
			&mtrc.Goroutines,
			&mtrc.InHeap,
			&mtrc.NsPerOp)
		if err != nil {
			return nil, err
		}
		result = append(result, &model.Metrics{
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
	dbRows, err := s.db.Query("SELECT queries,nsperop from memorymetrics m, date d where m.benchdate = d.id and d.benchdate=$1", date)
	if err != nil {
		return nil, err
	}
	for dbRows.Next() {
		err = dbRows.Scan(&mtrc.Query,
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
