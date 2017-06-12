package extraction

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/wirepair/autogcd"
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

var navigationTimeout = time.Second * 10
var stableAfter = time.Millisecond * 450

// Healthy is set to false when an unexpected error occured, that might indicate
// Chrome should be restarted
var Healthy = true

// Extract extracts all the technologies present on the passed URL
func Extract(auto *autogcd.AutoGcd, URL string) (Result, error) {
	result := Result{URL: URL, Applications: make([]Application, 0)}

	tab, err := auto.NewTab()
	if err != nil {
		Healthy = false
		log.Printf("error creating a tab: %v", err)
		return result, err
	}
	defer auto.CloseTab(tab)

	tab.SetNavigationTimeout(navigationTimeout)
	tab.SetStabilityTime(stableAfter)

	_, err = tab.Navigate(URL)
	if err != nil {
		log.Printf("error navigating to URL %v: %v\n", URL, err)
		return result, nil
	}

	tab.WaitStable()

	// appending to the page all the required wappalyzer files
	wappalyzerFile, err := getFileAsString("/extraction/js/wappalyzer.js")
	if err != nil {
		log.Printf("error opening wappalyzer file: %v\n", err)
		return result, err
	}

	_, err = tab.EvaluateScript(wappalyzerFile)
	if err != nil {
		log.Printf("error evaluating wappalyzer script: %v\n", err)
		return result, err
	}

	appsFile, err := getFileAsString("/extraction/js/apps.js")
	if err != nil {
		log.Printf("error opening apps file: %v\n", err)
		return result, err
	}

	_, err = tab.EvaluateScript(appsFile)
	if err != nil {
		log.Printf("error evaluating apps script: %v\n", err)
		return result, err
	}

	driverFile, err := getFileAsString("/extraction/js/driver.js")
	if err != nil {
		log.Printf("error opening driver file: %v\n", err)
		return result, err
	}

	_, err = tab.EvaluateScript(driverFile)
	if err != nil {
		log.Printf("error evaluating driver script: %v\n", err)
		return result, err
	}

	// tiny JS file which declares the function getDetectedApps()
	// that makes the "bridge" between Chrome and the Go code
	detectionFile, err := getFileAsString("/extraction/js/detection.js")
	if err != nil {
		log.Printf("error opening detection file: %v\n", err)
		return result, err
	}

	_, err = tab.EvaluateScript(detectionFile)
	if err != nil {
		log.Printf("error evaluating detection script: %v\n", err)
		return result, err
	}

	data, err := tab.EvaluateScript("getDetectedApps()")
	if err != nil {
		log.Printf("error evaluating detection script: %v\n", err)
		return result, err
	}

	applications := []Application{}

	for _, v := range data.Value.([]interface{}) {
		application := Application{Name: fmt.Sprintf("%v", v)}
		applications = append(applications, application)
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
