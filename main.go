package main

import (
	"flag"
	"flash-chart/app/controllers"
	"log"
	"os"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":"+os.Getenv("PORT"), "num of port")
	flag.Parse()
}

func main() {
	log.Println(controllers.StartWebServer(addr))
}
