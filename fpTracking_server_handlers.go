package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"encoding/json"
	"path"
	"html/template"
	"github.com/clementmaerten/fpTracking"
	"github.com/gorilla/sessions"
	_ "github.com/go-sql-driver/mysql"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//We check if the user has already the cookie
	//If he doesn't have, we give him one
	//The cookie contains a personal userId
	session, _ := store.Get(r, "fpTracking-cookie")
	if session.IsNew {
		log.Println("We create a new cookie")
		session.Options = &sessions.Options{
			Path: "/",
			MaxAge: 86400, //The cookie last 1 day at maximum
			HttpOnly: true,
		}
		session.Values["userId"] = generate_new_id()
		session.Save(r, w)
	}
	log.Println("userId :",session.Values["userId"])


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

func checkProgressionHandler(w http.ResponseWriter, r *http.Request) {

	//We check if the user has a cookie with a userId
	session, _ := store.Get(r, "fpTracking-cookie")
	if session.IsNew {
		http.Error(w, "You don't have the cookie", http.StatusForbidden)
		return
	}

	//We check if the tracking algorithm has begun
	if _, is_present := progressInformationSession[session.Values["userId"].(string)]; !is_present {
		http.Error(w, "The tracking algorithm wasn't launched", http.StatusForbidden)
		return
	}

	js, err := json.Marshal(progressInformationSession[session.Values["userId"].(string)])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Server","A Fingerprint tracking Go WebServer")
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	w.Write(js)
}

func trackingParallelHandler(w http.ResponseWriter, r *http.Request) {


	//We check if the user has a cookie with a userId
	session, _ := store.Get(r, "fpTracking-cookie")
	if session.IsNew {
		http.Error(w, "You don't have the cookie", http.StatusForbidden)
		return
	}


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


	//We calculate all the replaySequence to know the total number of fingerprints to analyze
	//We store these replay sequences and we send them to the ReplayScenario program.
	visitFrequencyToReplaySequence := make(map[int][]fpTracking.SequenceElt)
	lengths := make(map[int]int)
	totalLength := 0
	for _, visitFreq := range visitFrequencies {
		visitFrequencyToReplaySequence[visitFreq] = fpTracking.GenerateReplaySequence(test,visitFreq)
		lengths[visitFreq] = len(visitFrequencyToReplaySequence[visitFreq])
		totalLength += lengths[visitFreq]
	}

	//We create the channel and we lauch the goroutine which is going to listen to the messages
	progressChannel := make(chan fpTracking.ProgressMessage, 100)
	defer close(progressChannel)
	go listenFpTrackingProgressChannel(totalLength, visitFrequencies, lengths,
		session.Values["userId"].(string), progressChannel)


	for _, visitFrequency := range visitFrequencies {
		scenarioResult := fpTracking.ReplayScenarioParallelWithProgressInformation(test,
			visitFrequency, fpTracking.RuleBasedLinkingParallel, goroutineNumber,
			visitFrequencyToReplaySequence[visitFrequency], progressChannel)

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