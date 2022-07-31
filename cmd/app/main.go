package main

import (
	"github.com/valyala/fasthttp"
	"log"
	"restapi2/internal/startup"
)

func main() {
	requestHandler := startup.CreateRequestHandler()
	if err := fasthttp.ListenAndServe(":8080", requestHandler); err != nil {
		log.Fatal(err)
	}
}
