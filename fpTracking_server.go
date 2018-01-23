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
	http.HandleFunc("/testPost/", testPostHandler)
	http.HandleFunc("/tracking_parallel/",trackingParallelHandler)
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
	}
}

func trackingParallelHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("trackingParallelHandler launched")

	number, err1 := strconv.Atoi("5000")
	train, err2 := strconv.ParseFloat("0.01",64)
	if (err1 != nil || err2 != nil) {
		fmt.Println("The format is not respected !")
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}

	fingerprintManager := fpTracking.FingerprintManager{
		Number: number,
		Train:  train,
		MinNumberFpPerUser: 6,
		DBInfo: fpTracking.DBInformation {
			DBType: "mysql",
			User: "root",
			Password: "mysql",
			TCP: "",
			DBName: "fingerprint",
		},
	}

	fmt.Printf("Start fetching fingerprints\n")
	_, test := fingerprintManager.GetFingerprints()
	fmt.Printf("Fetched %d fingerprints\n", len(test))
	visitFrequencies := []int{1, 2, 3, 4, 5, 6, 7, 8, 10, 15, 20}
	var jsonResults []fpTracking.ResultsForVisitFrequency

	for _, visitFrequency := range visitFrequencies {
		//DANS UN PREMIER TEMPS EN SEQUENTIEL POUR LE TEST
		scenarioResult := fpTracking.ReplayScenario(test, visitFrequency, fpTracking.RuleBasedLinking)

		jsonResults = append(jsonResults,fpTracking.AnalyseScenarioResultInJSON(visitFrequency, scenarioResult, test))
	}

	js, err := json.Marshal(jsonResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Server","A Fingerprint tracking Go WebServer")
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	w.Write(js)
}