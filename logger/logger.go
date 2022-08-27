package logger

import (
	"github.com/valyala/fasthttp"
	"log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Middleware(ctx *fasthttp.RequestCtx) {
	log.Printf("%s - %d - %s", ctx.Method(), ctx.Response.StatusCode(), ctx.RequestURI())
}

func (l *Logger) Print(any ...any) {
	log.Print(any)
}

func (l *Logger) Fatal(any ...any) {
	log.Fatal(any)
}
