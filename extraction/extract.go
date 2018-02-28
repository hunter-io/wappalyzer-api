package extraction

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Result contains the result of a wappalyzer extraction
type Result struct {
	URL          string        `json:"url"`
	Applications []Application `json:"applications"`
}

// Application contains an extracted application
type Application struct {
	Name string `json:"name"`
}

// Healthy is set to false when an unexpected error occured, that might indicate
// the container should be restarted
var Healthy = true

// Extract extracts all the technologies present on the passed URL
func Extract(ctxt context.Context, chrome *chromedp.CDP, URL string) (Result, error) {
	result := Result{URL: URL, Applications: make([]Application, 0)}

	wappalyzerFile, err := getFileAsString("/extraction/js/wappalyzer.js")
	if err != nil {
		log.Printf("error opening wappalyzer file: %v\n", err)
		return result, err
	}

	appsFile, err := getFileAsString("/extraction/js/apps.js")
	if err != nil {
		log.Printf("error opening apps file: %v\n", err)
		return result, err
	}

	driverFile, err := getFileAsString("/extraction/js/driver.js")
	if err != nil {
		log.Printf("error opening driver file: %v\n", err)
		return result, err
	}

	detectionFile, err := getFileAsString("/extraction/js/detection.js")
	if err != nil {
		log.Printf("error opening detection file: %v\n", err)
		return result, err
	}

	var apps []string
	var evaluationResult *runtime.RemoteObject

	err = chrome.Run(ctxt, chromedp.Tasks{
		chromedp.Navigate(URL),
		chromedp.Sleep(2 * time.Second),
		chromedp.Evaluate(wappalyzerFile, &evaluationResult),
		chromedp.Evaluate(appsFile, &evaluationResult),
		chromedp.Evaluate(driverFile, &evaluationResult),
		chromedp.Evaluate(detectionFile, &evaluationResult),
		chromedp.Evaluate(`getDetectedApps();`, &apps),
	})
	if err != nil {
		log.Printf("error detecting apps: %v\n", err)
		return result, err
	}

	applications := []Application{}

	for _, app := range apps {
		applications = append(applications, Application{Name: app})
	}

	result.Applications = applications

	log.Printf("found %d applications for %v\n", len(applications), URL)

	// the extraction succeeded, our app is healthy
	Healthy = true

	return result, nil
}

func getFileAsString(filePath string) (string, error) {
	pwd, _ := os.Getwd()
	file, err := ioutil.ReadFile(pwd + filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}
