package model

type Bench struct {
	Metrics []BenchStat
}

type BenchStat struct {
	Date string `json:"date"`
	Stat Stat   `json:"stat"`
}

type Stat struct {
	Catalog string `json:"catalog"`
	Filters string `json:"filters"`
}
