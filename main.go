package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
	"log"
	"restapi2/internal/handlers"
	"restapi2/internal/service"
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
	cache := service.NewCache()
	err := cache.Load()
	if err != nil {
		log.Fatal("error load cache:", err)
	}
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stat":
			handlers.StatHandler(ctx, cache)
		case "/full":
			handlers.FullHandler(ctx, cache)
		case "/compare":
			handlers.CompareHandler(ctx, cache)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatal(err)
	}

}
