package handlers

import (
	"encoding/json"
	"log"

	"restapi2/internal/service"

	"github.com/valyala/fasthttp"
)

func CompareHandler(ctx *fasthttp.RequestCtx, cache *service.CacheBench) {
	writeCors(ctx)
	data := ctx.Request.Body()
	if data == nil {
		log.Print("error get data from request: ", ctx.Request.Body())
		writeError(ctx, "error: no data", fasthttp.StatusNoContent)
		return
	}

	delta, err := service.CountDifference(string(data), cache)
	if err != nil {
		log.Print("error get difference from service: ", err, ctx.Request.Body())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	jsonDelta, err := json.Marshal(delta)
	if err != nil {
		log.Print("error marshal json delta: ", err, ctx.Request.Body())
		writeError(ctx, internalError, fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.SetBody(jsonDelta)

	log.Print(string(ctx.Request.RequestURI()))
}

func writeCors(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "localhost:8080")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", []byte("Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"))
}
