package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func serve(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s\n", getCurrentIP(*r), r.Method, r.URL)
	content := r.URL.Query().Get("text")
	bs, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		log.Printf("err: %s", content)
		w.Header().Set("Content-Type", "txt/plain")
		fmt.Fprintf(w, "Error!")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(bs)
}

func empty(w http.ResponseWriter, r *http.Request) {
}

func getCurrentIP(r http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if len(ip) > 0 {
		ip = strings.Split(ip, ",")[0]
	}
	if len(ip) == 0 {
		ip = r.Header.Get("X-Real-IP")
	}
	if len(ip) == 0 {
		ip = r.RemoteAddr
	}
	return ip
}

func main() {
	httpPort := 9981
	http.HandleFunc("/favicon.ico", empty)
	http.HandleFunc("/qrcode", serve)
	fmt.Printf("listening on %v\n", httpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}
