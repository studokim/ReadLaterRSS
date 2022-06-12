package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/studokim/ReadLaterRSS/internal"
)

func main() {
	args := os.Args
	if len(args) > 1 && (args[1] == "--listen" || args[1] == "-l") {
		staticServer := http.FileServer(http.Dir("static"))
		http.Handle("/static/", http.StripPrefix("/static", staticServer))
		h := internal.NewHandler()
		http.HandleFunc("/add", h.Add)
		http.HandleFunc("/rss", h.Rss)
		http.ListenAndServe(":"+args[2], nil)
	} else {
		fmt.Println("Usage: ./ReadLaterRSS --listen <port>")
		fmt.Println("or     ./ReadLaterRSS -l       <port>")
	}
}
