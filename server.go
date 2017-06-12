package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/hunter-io/docker-wappalyzer-api/extraction"
	"github.com/wirepair/autogcd"
)

var (
	port       int
	chromePath string
	userDir    string
	chromePort string
)

var startupFlags = []string{"--disable-new-tab-first-run", "--no-first-run", "--disable-translate", "--headless", " --disable-gpu", "--ignore-certificate-errors", "--allow-running-insecure-content", "--no-sandbox"}

func init() {
	flag.IntVar(&port, "port", 3001, "port number to serve the API on")
	flag.StringVar(&chromePath, "chromePath", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&userDir, "tmpDir", "/tmp/", "temp directory")
	flag.StringVar(&chromePort, "chromePort", "9222", "debugger port")
}

func main() {
	flag.Parse()

	// starts autoGcd
	autoGcd, startErr := createAutoGcd()
	if startErr != nil {
		log.Fatalf("cannot start Chrome: %v", startErr)
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "200 OK")
		return
	})

	// health-check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if extraction.Healthy {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "200 OK")
			return
		}

		writeResponseError(w, http.StatusInternalServerError, errors.New("Chrome must be restarted"))
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

		result, err := extraction.Extract(autoGcd, URLToExtractFrom)
		if err != nil {
			// failure during the extraction
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

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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

func randUserDir() string {
	dir, err := ioutil.TempDir(userDir, "autogcd")
	if err != nil {
		log.Printf("error getting temp dir: %s\n", err)
	}
	return dir
}

func createAutoGcd() (*autogcd.AutoGcd, error) {
	settings := autogcd.NewSettings(chromePath, randUserDir())
	settings.RemoveUserDir(true)
	settings.AddStartupFlags(startupFlags)

	auto := autogcd.NewAutoGcd(settings)
	auto.SetTerminationHandler(nil)

	err := auto.Start()
	if err != nil {
		return nil, err
	}

	return auto, nil
}
