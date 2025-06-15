package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/studokim/ReadLaterRSS/internal"
)

//go:embed static
var static embed.FS

//go:embed html
var html embed.FS

func ipsV4() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if ip := v.IP.To4(); ip != nil {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	return ips, nil
}

func main() {
	args := os.Args
	if len(args) == 5 &&
		(args[1] == "--listen" || args[1] == "-l") &&
		(args[3] == "--website" || args[3] == "-s") {
		port := args[2]
		rootUrl := args[4]

		http.Handle("/static/", http.FileServer(http.FS(static)))
		_, err := internal.NewHandler(rootUrl, html)
		if err != nil {
			panic(err)
		}

		ips, err := ipsV4()
		if err != nil {
			panic(err)
		}
		for i, ip := range ips {
			if i == 0 {
				fmt.Printf("Serving at http://%s:%s\n", ip, port)
			} else {
				fmt.Printf("           http://%s:%s\n", ip, port)
			}
		}

		err = http.ListenAndServe("0.0.0.0:"+port, nil)
		panic(err) // always non-nil
	} else {
		fmt.Println("Usage: ./ReadLaterRSS --listen|-l <port> --website|-s <example.com>")
	}
}
