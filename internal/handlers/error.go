package handlers

import "github.com/valyala/fasthttp"

const internalError = "internal error"

func writeError(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	ctx.SetBodyString(msg)
	ctx.SetStatusCode(statusCode)
}
