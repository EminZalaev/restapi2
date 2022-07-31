package handlers

import (
	"encoding/json"
	"log"

	"restapi2/internal/service"

	"github.com/valyala/fasthttp"
)

func StatHandler(ctx *fasthttp.RequestCtx, cache *service.CacheBench) {
	writeCors(ctx)

	benchStat, err := service.GetBenchStat(cache)
	if err != nil {
		log.Print("error get stat from service: ", err, ctx.Response.String())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(benchStat)
	if err != nil {
		log.Print("error unmarshal json stat: ", err, ctx.Response.String())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	if data == nil {
		log.Print("error: stat is nil", ctx.Response.String())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(data)
}
