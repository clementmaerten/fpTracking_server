//Website for the differents response types : http://www.alexedwards.net/blog/golang-response-snippets
//Website for the web sessions : http://www.gorillatoolkit.org/pkg/sessions
//Website for the ECMA6 : https://www.wanadev.fr/21-introduction-a-ecmascript-6-le-javascript-de-demain/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sort"
	"strconv"
	"os/signal"
	"encoding/json"
	"path"
	"html/template"
	"github.com/clementmaerten/fpTracking"
	_ "github.com/go-sql-driver/mysql"
)

const TEMPLATES_FOLDER = "templates"

func handledFunctions() {
	http.HandleFunc("/",indexHandler)
	http.HandleFunc("/test-post/", testPostHandler)
	http.HandleFunc("/tracking-parallel/",trackingParallelHandler)
}

func main() {
	log.Println("Starting HTTP server on http://localhost:8080")

	// subscribe to SIGINT signals
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	go func() {
		<-quit
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}
	}()

	//Handled functions
	handledFunctions()

	//Handled static files (like CSS and JS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//Start the server and listen to requests
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}


	log.Println("Server stopped")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server","A Fingerprint tracking Go WebServer")
	w.Header().Set("Content-Type","text/html; charset=UTF-8")

	html_file := path.Join(TEMPLATES_FOLDER,"index.html")
	tmpl, err := template.ParseFiles(html_file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w,nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func testPostHandler(w http.ResponseWriter, r *http.Request) {
	
	log.Println("testPostHandler launched")

	r.ParseForm()
	fmt.Println(r.Form)
	//fmt.Println(r.FormValue("name"))
	for key, value := range r.Form {
		fmt.Println(key,value)
		//fmt.Println(key,":",r.FormValue(key))
	}
}

func trackingParallelHandler(w http.ResponseWriter, r *http.Request) {


	//Parse the parameters in a map
	r.ParseForm()

	number, err1 := strconv.Atoi(r.FormValue("fpTrackingParallelNumber"))
	minNbPerUser, err2 := strconv.Atoi(r.FormValue("fpTrackingParallelMinNbPerUser"))
	goroutineNumber, err3 := strconv.Atoi(r.FormValue("fpTrackingParallelGoroutineNumber"))
	train := float64(0)
	if err1 != nil || err2 != nil || err3 != nil || number <= 0 || minNbPerUser <= 0 || goroutineNumber <= 0 {
		log.Println("Error in the format in trackingParallelHandler")
		if err1 != nil {
			http.Error(w, err1.Error(), http.StatusBadRequest)
		} else if err2 != nil {
			http.Error(w, err2.Error(), http.StatusBadRequest)
		} else if err3 != nil {
			http.Error(w, err3.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, "Nul or negative parameters", http.StatusBadRequest)
		}
		return
	}

	//conversion of string slice visitFrequencies to int slice
	var visitFrequencies []int
	for _, stringValue := range r.Form["fpTrackingParallelVisitFrequency"] {
		intValue, err := strconv.Atoi(stringValue)
		if err != nil {
			log.Println("Error in the format of visitFrequencies in trackingParallelHandler")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		visitFrequencies = append(visitFrequencies,intValue)
	}
	if len(visitFrequencies) < 1 {
		log.Println("Not enough arguments for visitFrequencies in trackingParallelHandler")
		http.Error(w, "Not enough arguments for visitFrequencies", http.StatusBadRequest)
		return
	}

	//Sort the visitFrequencies
	sort.Ints(visitFrequencies)


	log.Println("trackingParallelHandler launched with number =",number,
		", minNbPerUser =",minNbPerUser,", visitFrequencies =",visitFrequencies,", goroutineNumber =",goroutineNumber)

	progressChannel := make(chan fpTracking.ProgressMessage, 100)
	defer close(progressChannel)
	go listenFpTrackingProgressChannel(progressChannel)

	fingerprintManager := fpTracking.FingerprintManager{
		Number: number,
		Train:  train,
		MinNumberFpPerUser: minNbPerUser,
		DBInfo: fpTracking.DBInformation {
			DBType: "mysql",
			User: "root",
			Password: "mysql",
			TCP: "",
			DBName: "fingerprint",
		},
	}

	_, test := fingerprintManager.GetFingerprints()

	var jsonResults []fpTracking.ResultsForVisitFrequency

	for _, visitFrequency := range visitFrequencies {
		scenarioResult := fpTracking.ReplayScenarioParallelWithProgressInformation(test,
			visitFrequency, fpTracking.RuleBasedLinkingParallel, goroutineNumber, progressChannel)

		jsonResults = append(jsonResults,fpTracking.AnalyseScenarioResultInJSON(visitFrequency, scenarioResult, test))
	}

	progressChannel <- fpTracking.ProgressMessage{Task : fpTracking.CLOSE_GOROUTINE}

	js, err := json.Marshal(jsonResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Server","A Fingerprint tracking Go WebServer")
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	w.Write(js)
}

//This function listen to the progress channel and update the user's session with progress information
//This function is supposed to be executed by a goroutine
func listenFpTrackingProgressChannel(ch <- chan fpTracking.ProgressMessage) {
	for {
		rq := <- ch
		if strings.Compare(rq.Task, fpTracking.SEND_PROGRESS_INFORMATION) == 0 {
			log.Println("visitFrequency :",rq.VisitFrequency,", progression :",rq.Progression)
		} else if strings.Compare(rq.Task, fpTracking.CLOSE_GOROUTINE) == 0 {
			return
		} else {
			//This case should never happen
			log.Println("Wrong task for listenFpTrackingProgressChannel")
			return
		}
	}
}