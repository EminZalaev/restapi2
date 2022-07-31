package handlers

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"log"
	"restapi2/internal/service"
)

func FullHandler(ctx *fasthttp.RequestCtx, bench *service.CacheBench) {
	writeCors(ctx)

	fullResult := bench.Get(0)

	if fullResult == nil {
		log.Print("error full result is nil: ", ctx.Response.String())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fullResult)
	if err != nil {
		log.Print("error marshall full result json: ", err, ctx.Response.String())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(data)
}
