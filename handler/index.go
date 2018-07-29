package handler

import (
	"fmt"
	"github.com/yanchenghust/goblog/init/log"
	"html"
	"net/http"
)

type IndexHandler struct {
}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Hello, key is %q", html.EscapeString(r.Form.Get("key")))
	log.Infof("Hello, key is %q", html.EscapeString(r.Form.Get("key")))
}
