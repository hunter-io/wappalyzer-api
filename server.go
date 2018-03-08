package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hunter-io/docker-wappalyzer-api/extraction"
	"github.com/tebeka/selenium"
)

var (
	port        int    // port number the API will be served on
	seleniumURL string // URL of the remote Selenium instance to connect to
)

func init() {
	flag.IntVar(&port, "port", 3001, "port number to serve the API on")
	flag.StringVar(&seleniumURL, "seleniumURL", "http://localhost:4444/wd/hub", "Selenium URL to connect to")

	flag.Parse()
}

func main() {
	caps := selenium.Capabilities{"browserName": "chrome"}

	wd, err := selenium.NewRemote(caps, seleniumURL)
	if err != nil {
		log.Fatalf("cannot connect to Selenium: %v", err)
		return
	}
	defer wd.Quit()

	// limits the concurrency to one extraction at a time
	ch := make(chan bool, 1)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
		return
	})

	// health-check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if extraction.FailedExtractions < 10 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "200 OK")
			return
		}

		writeResponseError(w, http.StatusInternalServerError, errors.New("App is unhealthy"))
		return
	})

	// extract endpoint
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
			// passed URL is invalid
			writeResponseError(w, http.StatusUnprocessableEntity, err)
			return
		}

		ch <- true

		result, err := extraction.Extract(wd, URLToExtractFrom)

		<-ch

		if err != nil {
			log.Printf("extraction failed (count: %d\n)", extraction.FailedExtractions)
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

	serverErr := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if serverErr != nil {
		log.Fatalf("cannot start http server: %v", serverErr)
		return
	}
}

func writeResponseError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("%v: %v", http.StatusText(statusCode), err)
	w.WriteHeader(statusCode)
	fmt.Fprint(w, fmt.Sprintf("%d %v", statusCode, http.StatusText(statusCode)))
}
