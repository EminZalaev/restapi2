package startup

import (
	"github.com/valyala/fasthttp"
	"log"
	"restapi2/internal/handlers"
	"restapi2/internal/service"
)

func CreateRequestHandler() func(ctx *fasthttp.RequestCtx) {
	cache := service.NewCache()
	err := cache.Load()
	if err != nil {
		log.Print("error load json file to cache")
	}
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/stat":
			handlers.StatHandler(ctx, cache)
		case "/compare":
			handlers.CompareHandler(ctx, cache)
		case "/full":
			handlers.FullHandler(ctx, cache)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
}
