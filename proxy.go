package main

import (
	"fmt"
	"net/http"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != *proxyMetricsPath {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	promData, found := storage.Get(r.URL.Host)

	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, promData)
}
