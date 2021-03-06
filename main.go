package main

import (
	"github.com/yanchenghust/goblog/handler"
	"github.com/yanchenghust/goblog/init/log"
	"golang.org/x/net/context"
	"net/http"
)

func main() {
	ctx := context.Background()
	log.InitLog()

	http.Handle("/index", handler.IndexHandler{})
	http.Handle("/favicon.ico", handler.IndexHandler{})
	http.Handle("/", handler.IndexHandler{})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Errorf("server error, err: %v", err)
	}
	log.StopLog()
}
