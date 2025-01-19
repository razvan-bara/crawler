package main

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	ipCountTable   = make(map[string]int)
	ipTimeoutTable = map[string]time.Time{}
)

func interceptorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIp := strings.Split(r.RemoteAddr, ":")[0]
		ipTimeout, ok := ipTimeoutTable[clientIp]
		if ok && time.Now().Before(ipTimeout) {
			log.Println(clientIp, "is timed out")
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}

		ipCount, ok := ipCountTable[clientIp]
		if ipCount > 1 {
			ipTimeoutTable[clientIp] = time.Now().Add(5 * time.Second)
			ipCountTable[clientIp] = 0
			w.WriteHeader(http.StatusGatewayTimeout)
			log.Println(clientIp, "timed out")
			return
		}
		ipCountTable[clientIp]++

		log.Printf("Intercepting request for: %s", r.URL.Path)

		next.ServeHTTP(w, r)
	})
}

func serveHTML(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join("./html_pages", filename)
		log.Printf("Serving file: %s", filePath)
		http.ServeFile(w, r, filePath)
	}
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime)

	r := mux.NewRouter()

	routes := map[string]string{
		"/citations_article1":      "citations_article1.html",
		"/citations_article2":      "citations_article2.html",
		"/citations_page_1":        "citations_page_1.html",
		"/citations_page_2":        "citations_page_2.html",
		"/citations_page_fragment": "citations_page_fragment.html",
		// real paths below
		"/pid/80/2813":  "dblp_dc.html",
		"/pid/20/123":   "dblp_se.html",
		"/pers/d/{*}":   "dblp_se.html",
		"/pers":         "dblp_pers_index_1.html",
		"/pers?pos=301": "dblp_pers_index_2.html",
		"/pers?pos=601": "dblp_pers_index_3.html",
	}

	for route, file := range routes {
		r.HandleFunc(route, serveHTML(file))
	}

	handlerWithMiddleware := interceptorMiddleware(r)

	addr := ":8080"
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, handlerWithMiddleware); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
