package main

import (
	"log"
	"net/http"

	"github.com/qichengzx/seqsvr/service"
)

func main() {
	var conf = service.NewConfig()
	var svr = service.New(conf)

	http.HandleFunc("/new", svr.ServeHttp)

	log.Println("ID server is listening on", conf.PORT)
	log.Fatal(http.ListenAndServe(conf.PORT, nil))
}
