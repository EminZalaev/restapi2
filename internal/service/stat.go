package service

import "restapi2/internal/model"

func GetBenchStat(cache *CacheBench) (*model.Bench, error) {
	time, err := readBenchTimeStat(cache)
	if err != nil {
		return nil, err
	}
	reverseStat(time)
	return time, nil
}

func readBenchTimeStat(cache *CacheBench) (*model.Bench, error) {
	tmp := cache.Get(0)
	stat := make([]model.BenchStat, 0, len(tmp.Metrics))
	for i := 0; i < len(tmp.Metrics); i++ {
		stat = append(stat, model.BenchStat{
			Date: tmp.Metrics[i].Date,
			Stat: model.Stat{
				Catalog: tmp.Metrics[i].NsPerOp.Catalog,
				Filters: tmp.Metrics[i].NsPerOp.Filters,
			},
		})
	}
	return &model.Bench{
		Metrics: stat,
	}, nil
}

func reverseStat(t *model.Bench) {
	for i, j := 0, len(t.Metrics)-1; i < j; i, j = i+1, j-1 {
		t.Metrics[i], t.Metrics[j] = t.Metrics[j], t.Metrics[i]
	}
}
