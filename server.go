package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/bastienl/docker-wappalyzer-api/extraction"
)

var port = flag.Int("port", 3001, "Port number to serve the API on")

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
		return
	})

	http.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeResponseError(w, http.StatusBadRequest, errors.New("Did not POST"))
			return
		}

		err := r.ParseForm()
		if err != nil {
			writeResponseError(w, http.StatusBadRequest, err)
			return
		}

		URLToExtractFrom := r.PostFormValue("url")

		_, err = url.ParseRequestURI(URLToExtractFrom)
		if err != nil {
			writeResponseError(w, http.StatusUnprocessableEntity, err)
			return
		}

		result, err := extraction.Extract(URLToExtractFrom)
		if err != nil {
			writeResponseError(w, http.StatusInternalServerError, err)
			return
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			writeResponseError(w, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start http server: %s\n", err)
		os.Exit(1)
	}
}

func writeResponseError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("%v: %v", http.StatusText(statusCode), err)
	w.WriteHeader(statusCode)
	fmt.Fprint(w, fmt.Sprintf("%d %v", statusCode, http.StatusText(statusCode)))
}
