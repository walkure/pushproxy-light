package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/walkure/pushproxy-light/converter"
)

func pushHandler(w http.ResponseWriter, r *http.Request) {

	signature, body, err := getBodyAndSignature(r)
	if err != nil {
		log.Printf("Request invalid:%+v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	host := vars["host"]

	if host == "" {
		log.Printf("host required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	retrieveData(w, host, signature, body)
}

func retrieveData(w http.ResponseWriter, host, signature, body string) {
	c, err := generateSignature(signature[0], *preSharedKey, body)
	if err != nil {
		log.Printf("Cannot generate signature:%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if c != signature[1:] {
		log.Printf("Signature mismatch[%v] != [%v]\n", c, signature)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var promDataBuilder strings.Builder
	if err := converter.ConvertMetrics(&promDataBuilder, body); err != nil {
		log.Printf("Signature mismatch[[%v]\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	promData := promDataBuilder.String()

	storage.SetDefault(host, promData)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, promData)

}

func generateSignature(alg byte, shareKey, body string) (string, error) {

	var buf bytes.Buffer
	buf.Grow(len(shareKey) + len(body))
	buf.WriteString(shareKey)
	buf.WriteString(body)

	switch alg {
	case '1':
		hashBinary := md5.Sum(buf.Bytes())
		return fmt.Sprintf("%x", hashBinary), nil
	case '5':
		hashBinary := sha256.Sum256(buf.Bytes())
		return fmt.Sprintf("%x", hashBinary), nil
	case '6':
		hashBinary := sha512.Sum512(buf.Bytes())
		return fmt.Sprintf("%x", hashBinary), nil
	default:
		return "", fmt.Errorf("Alg type=%c unknown", alg)
	}
}

func getBodyAndSignature(r *http.Request) (signature, body string, e error) {

	switch r.Method {
	case http.MethodGet:
		signature = r.URL.Query().Get("signature")
		body = r.URL.Query().Get("body")

		if signature == "" {
			e = errors.New("Signature not found")
			return
		}

		if body == "" {
			e = errors.New("Body not found")
			return
		}
	case http.MethodPost:
		switch r.Header.Get("Content-Type") {

		case "application/json":
			signature = r.Header.Get("X-Signature")
			if signature == "" {
				e = errors.New("Signature not found")
				return
			}

			var sb strings.Builder
			if _, err := io.Copy(&sb, r.Body); err != nil {
				// r.Body does not requires Close (see https://golang.org/pkg/net/http/#Request )
				e = fmt.Errorf("Retrieve body error: %w", err)
				return
			}
			body = sb.String()
			if body == "" {
				e = errors.New("Body not found")
				return
			}

		case "application/x-www-form-urlencoded":
			if err := r.ParseForm(); err != nil {
				e = fmt.Errorf("ParseForm Error: %w", err)
				return
			}
			signature = r.PostForm.Get("signature")
			body = r.PostForm.Get("body")
		}
	default:
		e = fmt.Errorf("Method[%s] not allowed.", http.MethodPost)
		return
	}
	return
}
