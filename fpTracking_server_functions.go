package main

import (
	"log"
	"fmt"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/clementmaerten/fpTracking"
)



type progressInformationStruct struct {
	Progression int
	Results []fpTracking.ResultsForVisitFrequency
}

//This function listen to the progress channel and update the user's session with progress information
//This function is supposed to be executed by a goroutine
func listenFpTrackingProgressChannel(totalLength int, sortedVisitFrequencies []int, lengths map[int]int,
	userId string, ch <- chan fpTracking.ProgressMessage) {

	currentVisitFrequency := sortedVisitFrequencies[0]
	indexAtNewVisitFrequency := 0
	globalProgression := 0

	if progressInformationSession[userId] == nil {
		progressInformationSession[userId] = &progressInformationStruct{}
	}
	
	for {
		rq := <- ch
		if strings.Compare(rq.Task, fpTracking.SEND_PROGRESS_INFORMATION) == 0 {

			if rq.VisitFrequency != currentVisitFrequency {
				indexAtNewVisitFrequency += lengths[currentVisitFrequency]
				currentVisitFrequency = rq.VisitFrequency
				log.Println("new visitFrequency :",currentVisitFrequency)
			}

			globalProgression = (indexAtNewVisitFrequency + rq.Index) * 100 / totalLength

			log.Println("progression :",globalProgression)
			progressInformationSession[userId].Progression = globalProgression

		} else if strings.Compare(rq.Task, fpTracking.SEND_RESULTS_FOR_VISIT_FREQUENCY) == 0 {

			progressInformationSession[userId].Results = append(progressInformationSession[userId].Results, rq.ResForVisitFreq)
		} else if strings.Compare(rq.Task, fpTracking.CLOSE_GOROUTINE) == 0 {

			globalProgression = 100
			progressInformationSession[userId].Progression = globalProgression
			return
		} else {
			//This case should never happen
			log.Println("Wrong task for listenFpTrackingProgressChannel")
			return
		}
	}
}

func generateNewId() string {
	gen, _ := uuid.NewV4()
	return fmt.Sprintf("%s",gen)
}

func launchTrackingAlgorithm(number int, minNbPerUser int, goroutineNumber int,
		train float64, visitFrequencies []int, userId string) {

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
		userId, progressChannel)


	for _, visitFrequency := range visitFrequencies {
		scenarioResult := fpTracking.ReplayScenarioParallelWithProgressInformation(test,
			visitFrequency, fpTracking.RuleBasedLinkingParallel, goroutineNumber,
			visitFrequencyToReplaySequence[visitFrequency], progressChannel)

		log.Println("We send the results for visitFrequency",visitFrequency)
		progressChannel <- fpTracking.ProgressMessage {
			Task : fpTracking.SEND_RESULTS_FOR_VISIT_FREQUENCY,
			ResForVisitFreq : fpTracking.AnalyseScenarioResultInJSON(visitFrequency, scenarioResult, test),
		}
	}

	progressChannel <- fpTracking.ProgressMessage{Task : fpTracking.CLOSE_GOROUTINE}

	/*js, err := json.Marshal(jsonResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Server","A Fingerprint tracking Go WebServer")
	w.Header().Set("Content-Type","application/json; charset=utf-8")
	w.Write(js)*/

	//log.Println("TrackingAlgorithm finished !")
	//log.Println("userId :",userId)
	//log.Println("Progress information :",progressInformationSession[userId])
}