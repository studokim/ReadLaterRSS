package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/studokim/ReadLaterRSS/internal"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

func main() {
	args := os.Args
	if len(args) >= 7 &&
		(args[1] == "--listen" || args[1] == "-l") &&
		(args[3] == "--website" || args[3] == "-s") &&
		(args[5] == "--author" || args[5] == "-a") {
		staticServer := http.FileServer(http.FS(static))
		http.Handle("/static/", staticServer)
		h := internal.NewHandler(templates, args[4], args[6])
		http.HandleFunc("/add", h.Add)
		http.HandleFunc("/rss", h.Rss)
		http.ListenAndServe(":"+args[2], nil)
	} else {
		fmt.Println("Usage: ./ReadLaterRSS --listen <port> --website <yoursite.com> --author <your.name>")
		fmt.Println("or     ./ReadLaterRSS -l       <port> -s       <yoursite.com> -a      <your.name>")
	}
}
