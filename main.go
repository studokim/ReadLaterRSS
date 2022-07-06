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

//go:embed html
var html embed.FS

func main() {
	args := os.Args
	if len(args) >= 9 &&
		(args[1] == "--listen" || args[1] == "-l") &&
		(args[3] == "--website" || args[3] == "-s") &&
		(args[5] == "--author" || args[5] == "-a") &&
		(args[7] == "--email" || args[7] == "-m") {
		staticServer := http.FileServer(http.FS(static))
		http.Handle("/static/", staticServer)
		fmt.Println("Loading history...")
		h := internal.NewHandler(html, args[4], args[6], args[8])
		h.RegisterEndpoints()
		fmt.Printf("Serving at http://127.0.0.1:%s\n", args[2])
		http.ListenAndServe(":"+args[2], nil)
	} else {
		fmt.Println("Usage: ./ReadLaterRSS --listen <port> --website <yoursite.com> --author <your.name> --email <you@yoursite.com>")
		fmt.Println("or     ./ReadLaterRSS -l       <port> -s       <yoursite.com> -a      <your.name>   -m <you@yoursite.com>")
	}
}
