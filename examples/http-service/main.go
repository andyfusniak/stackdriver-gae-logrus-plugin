package main

import (
	"fmt"
	"net/http"
	"os"

	stackdriver "github.com/andyfusniak/stackdriver-gae-logrus-plugin"
	"github.com/andyfusniak/stackdriver-gae-logrus-plugin/middleware"

	log "github.com/sirupsen/logrus"
)

var version string

var projectID = "glowing-market-234808"

func sayHello(w http.ResponseWriter, r *http.Request) {
	contextLogger := log.WithContext(r.Context())

	contextLogger.WithFields(log.Fields{
		"battery": "50",
	}).Debug("Flux capacitor low")

	contextLogger.WithFields(log.Fields{
		"status": "busted",
	}).Info("Warp speed activated")

	contextLogger.WithFields(log.Fields{
		"status": "hmmm",
	}).Warn("You have been warning")

	contextLogger.WithFields(log.Fields{
		"status": "busted",
	}).Error("These are not the drones you are looking for")

	w.Write([]byte("hello"))
}

func main() {
	// Log as JSON Stackdriver with entry threading
	// instead of the default ASCII formatter.
	formatter := stackdriver.GAEStandardFormatter(
		stackdriver.WithProjectID(projectID),
	)
	log.SetFormatter(formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Log the debug severity or above.
	log.SetLevel(log.DebugLevel)

	mux := http.NewServeMux()
	mux.HandleFunc("/", sayHello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Stackdriver GAE Example HTTP Service version %s listening on port %s (projectID=%s)\n", version, port, projectID)
	fmt.Println("If you are running, simulate a GAE request locally using")
	fmt.Println("curl -v -H 'X-Cloud-Trace-Context: 1ad1e4f50427b51eadc9b36064d40cc2/8196282844182683029;o=1' http://localhost:8080/")
	fmt.Println()
	fmt.Println("Use make deploy to run this example on Google App Engine")
	fmt.Println()
	fmt.Println("JSON logging is sent to stdout")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), middleware.XCloudTraceContext(mux)))
}
