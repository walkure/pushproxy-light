package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

var preSharedKey = flag.String("preSharedKey", "", "Metrics sharing key")
var metricsLifetime = flag.Int("metricsLifetime", 5, "Metrics TTL in minutes")
var httpListener = flag.String("httpListener", ":8080", "HTTP Proxy/PushReceiver Listener Address")
var proxyMetricsPath = flag.String("proxyMetricsPath", "/metrics", "Path of Metrics URI")

var storage *cache.Cache

func main() {

	flag.Parse()

	if *preSharedKey == "" {
		flag.PrintDefaults()
		log.Fatalf("`-preSharedKey` is mandatory argument.")
	}

	cacheLifetime := time.Duration(*metricsLifetime) * time.Minute
	storage = cache.New(cacheLifetime, cacheLifetime*3)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/push/{host}", pushHandler)

	http.ListenAndServe(*httpListener, &httpHandler{router: router, proxy: http.HandlerFunc(proxyHandler)})
}

type httpHandler struct {
	router http.Handler
	proxy  http.Handler
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Host != "" {
		// Proxy request
		h.proxy.ServeHTTP(w, r)
	} else {
		// Client requests
		h.router.ServeHTTP(w, r)
	}
}
