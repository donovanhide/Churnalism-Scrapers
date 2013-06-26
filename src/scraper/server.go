package scraper

import (
	"flag"
	"github.com/donovanhide/eventsource"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func Serve(srv *eventsource.Server) {
	flag.Parse()
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/articles", srv.Handler("articles"))
	http.HandleFunc("/statistics", srv.Handler("statistics"))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
