package api

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"webserver/internal/domain/service"
	"webserver/logger"
)

type server struct {
	service *service.Service
	logger  *logger.Logger
}

func NewServer(srvc *service.Service, log *logger.Logger) *server {
	return &server{
		service: srvc,
		logger:  log,
	}
}

func (s *server) Start(port string) error {
	requestHandler := s.createRequestHandler()
	log.Println("server started at port:", port)

	if err := fasthttp.ListenAndServe(":"+port, requestHandler); err != nil {
		return fmt.Errorf("error start server: %w", err)
	}

	return nil
}

func (s *server) createRequestHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/full":
			writeCors(ctx)
			s.processFull(ctx)
			s.logger.Middleware(ctx)
		case "/timestat":
			writeCors(ctx)
			s.processStat(ctx)
			s.logger.Middleware(ctx)
		case "/compare":
			writeCors(ctx)
			s.processCompare(ctx)
			s.logger.Middleware(ctx)
		case "/all":
			writeCors(ctx)
			s.processAll(ctx)
			s.logger.Middleware(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
			s.logger.Middleware(ctx)
		}
	}
}

func (s *server) processStat(ctx *fasthttp.RequestCtx) {
	req := ctx.QueryArgs().Peek("date")
	if req == nil {
		s.logger.Print("error request args is nil")
		writeError(ctx, "internal error", fasthttp.StatusNoContent)
		return
	}
	stat, err := s.service.GetNsPerOp(req)
	if err != nil {
		s.logger.Print("error get ns per op: ", err)
		writeError(ctx, "internal error", fasthttp.StatusNoContent)
		return
	}
	ctx.Response.SetBody(stat)
}

func (s *server) processFull(ctx *fasthttp.RequestCtx) {
	req := ctx.QueryArgs().Peek("date")
	if req == nil {
		s.logger.Print("error request args is nil")
		writeError(ctx, "internal error", fasthttp.StatusNoContent)
		return
	}
	mem, err := s.service.GetFullResultByDateAndQuery(req)
	if err != nil {
		s.logger.Print("error get full result: ", err)
		writeError(ctx, "internal error", fasthttp.StatusNoContent)
		return
	}
	ctx.Response.SetBody(mem)
}

func (s *server) processCompare(ctx *fasthttp.RequestCtx) {
	req := ctx.QueryArgs().Peek("date")
	if req == nil {
		s.logger.Print("error request args is nil")
		writeError(ctx, "internal error", fasthttp.StatusNoContent)
		return
	}
	cmp, err := s.service.CompareResults(req)
	if err != nil {
		return
	}
	ctx.Response.SetBody(cmp)
}

func (s *server) processAll(ctx *fasthttp.RequestCtx) {
	cmp, err := s.service.GetAllResult()
	if err != nil {
		s.logger.Print("error get all result: ", err)
		writeError(ctx, "internal error", fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetBody(cmp)
}

func writeError(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	ctx.SetBodyString(msg)
	ctx.SetStatusCode(statusCode)
}

func writeCors(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.Response.Header.SetBytesV("Access-Control-Allow-Origin", []byte("Accept, Content-Type, Content-Length, Accept-Encoding"))
}
