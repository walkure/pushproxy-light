package main

import (
	"fmt"
	"net/http"
	"time"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != *proxyMetricsPath {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, found := storage.Get(r.URL.Host)

	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	promData := data.(*PrometheusData)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.Header().Set("Last-Modified", promData.LastUpdated.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, promData.Body)
}
