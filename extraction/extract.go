package extraction

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tebeka/selenium"
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

// FailedExtractions counts how many extractions have failed in a row.
var FailedExtractions = 0

// Extract extracts all the technologies present on the provided URL
func Extract(wd selenium.WebDriver, URL string) (Result, error) {
	defer func() {
		if e := recover(); e != nil {
			FailedExtractions++
		}
	}()

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

	_, err = wd.NewSession()
	if err != nil {
		log.Printf("error creating a new session: %v\n", err)
		FailedExtractions++
		return result, err
	}

	err = wd.Get(URL)
	if err != nil {
		log.Printf("error fetching %v: %v\n", URL, err)
		FailedExtractions++
		return result, err
	}

	data, err := wd.ExecuteScript(wappalyzerFile+" "+appsFile+" "+driverFile+" "+detectionFile+" "+"return getDetectedApps();", nil)
	if err != nil {
		log.Printf("error: %v", err.Error())
		return result, err
	}

	applications := []Application{}

	for _, v := range data.([]interface{}) {
		applications = append(applications, Application{Name: fmt.Sprintf("%v", v)})
	}

	// we end the current session
	wd.Quit()

	result.Applications = applications

	log.Printf("found %d applications for %v\n", len(applications), URL)

	FailedExtractions = 0

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
