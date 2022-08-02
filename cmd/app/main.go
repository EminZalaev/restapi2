package app

import (
	"github.com/valyala/fasthttp"
	"log"
	"restapi2/internal/service"
	"restapi2/internal/startup"
)

func Run() {
	cache := service.NewCache()
	err := cache.Load()
	if err != nil {
		log.Fatal("error load cache:", err)
	}
	requestHandler := startup.CreateRequestHandler()
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatal(err)
	}
}
