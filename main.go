package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
	"log"
)

type Bench struct {
	BenchDate string
	Catalog   string
	Filters   string
}

type MemMetrics struct {
	BenchDate            string
	OffHeap              int
	InHeap               int
	InStack              int
	TotalUsedMemory      int
	AllocationRates      int
	NumberOfLiveObjects  int
	RateObjectsAllocated int
	Goroutines           int
}

func main() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stat":
			statHandler(ctx)
		case "/full":
			fullHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatal(err)
	}

}

func getNsPerOp(date string) (*Bench, error) {
	var b Bench
	database, _ := sql.Open("sqlite3", "./benchresult.db")
	query, err := database.Query("select benchDate,catalog,filters from metrics where benchDate = $1", date)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		err = query.Scan(&b.BenchDate,
			&b.Catalog,
			&b.Filters)
		if err != nil {
			return nil, err
		}
	}
	return &b, nil

}

func getBenchMem(date string) (*MemMetrics, error) {
	var b MemMetrics
	database, _ := sql.Open("sqlite3", "./benchresult.db")
	query, err := database.Query("select benchdate,offHeap,inStack,inHeap,totalUsedMemory,allocationRates,numberOfLiveObjects,rateObjectsAllocated,Goroutines from benchMem where benchDate = $1", date)
	if err != nil {
		return nil, err
	}
	for query.Next() {
		err = query.Scan(
			&b.BenchDate,
			&b.OffHeap,
			&b.InHeap,
			&b.InStack,
			&b.TotalUsedMemory,
			&b.AllocationRates,
			&b.NumberOfLiveObjects,
			&b.RateObjectsAllocated,
			&b.Goroutines)
		if err != nil {
			return nil, err
		}
	}
	return &b, nil
}

func statHandler(ctx *fasthttp.RequestCtx) {
	reqBody := ctx.Request.Body()
	stat, err := getNsPerOp(string(reqBody))
	if err != nil {
		log.Println(err)
		return
	}
	data, err := json.Marshal(stat)
	if err != nil {
		return
	}
	ctx.Response.SetBody(data)
}

func fullHandler(ctx *fasthttp.RequestCtx) {
	reqBody := ctx.Request.Body()
	mem, err := getBenchMem(string(reqBody))
	if err != nil {
		log.Println(err)
		return
	}
	data, err := json.Marshal(mem)
	if err != nil {
		return
	}
	ctx.Response.SetBody(data)
}
