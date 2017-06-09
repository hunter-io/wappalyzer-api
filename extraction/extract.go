package extraction

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/wirepair/autogcd"
	"github.com/wirepair/gcd/gcdapi"
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

var (
	chromePath string
	userDir    string
	chromePort string
	debug      bool
)

var startupFlags = []string{"--disable-new-tab-first-run", "--no-first-run", "--disable-translate", "--headless", " --disable-gpu", "--ignore-certificate-errors", "--allow-running-insecure-content", "--no-sandbox"}

var waitFor = time.Second * 2
var navigationTimeout = time.Second * 10
var stableAfter = time.Millisecond * 450
var stabilityTimeout = time.Second * 2

func init() {
	flag.StringVar(&chromePath, "chromePath", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "path to chrome")
	flag.StringVar(&userDir, "tmpDir", "/tmp/", "temp directory")
	flag.StringVar(&chromePort, "chromePort", "9222", "debugger port")
	flag.BoolVar(&debug, "debug", false, "print console.log() outputs")
}

// Extract extracts all the technologies present on the passed URL
func Extract(URL string) (Result, error) {
	result := Result{URL: URL, Applications: make([]Application, 0)}

	settings := autogcd.NewSettings(chromePath, randUserDir())
	settings.RemoveUserDir(true)
	settings.AddStartupFlags(startupFlags)

	auto := autogcd.NewAutoGcd(settings)

	err := auto.Start()
	if err != nil {
		log.Printf("error starting Chrome: %v", err)
		return result, err
	}
	defer auto.Shutdown()

	auto.SetTerminationHandler(nil)

	tab, err := auto.GetTab()
	if err != nil {
		log.Printf("error creating a tab: %v", err)
		return result, err
	}

	tab.SetNavigationTimeout(navigationTimeout)
	tab.SetStabilityTime(stableAfter)

	if debug {
		msgHandler := func(callerTab *autogcd.Tab, message *gcdapi.ConsoleConsoleMessage) {
			fmt.Printf("%s\n", message.Text)
		}
		tab.GetConsoleMessages(msgHandler)
	}

	_, err = tab.Navigate(URL)
	if err != nil {
		log.Printf("error navigating to URL %v: %v\n", URL, err)
		return result, nil
	}

	// wait for page to load
	time.Sleep(waitFor)

	// appending to the tab all the wappalyzer files
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

	return result, nil
}

func randUserDir() string {
	dir, err := ioutil.TempDir(userDir, "autogcd")
	if err != nil {
		log.Printf("error getting temp dir: %s\n", err)
	}
	return dir
}

func getFileAsString(filePath string) (string, error) {
	pwd, _ := os.Getwd()
	file, err := ioutil.ReadFile(pwd + filePath)
	if err != nil {
		return "", err
	}

	return string(file), nil
}
