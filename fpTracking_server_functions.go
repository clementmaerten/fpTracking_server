package main

import (
	"log"
	"fmt"
	"math"
	"time"
	"strings"
	"github.com/satori/go.uuid"
	"github.com/clementmaerten/fpTracking"
)

type progressInformationStruct struct {
	creationDate time.Time
	inProgress bool
	Progression int
	CurrentVisitFrequency int
	AverageTrackingTimeGraph []fpTracking.GraphicPoint
	MaximumAverageTrackingTimeGraph []fpTracking.GraphicPoint
	NbIdsFrequencyGraph []fpTracking.GraphicPoint
	OwnershipFrequencyGraph []fpTracking.GraphicPoint
	ExecutingTime float64
}

//This function listen to the progress channel and update the user's session with progress information
//This function is supposed to be executed by a goroutine
func listenFpTrackingProgressChannel(totalLength int, sortedVisitFrequencies []int, lengths map[int]int,
	userId string, ch chan fpTracking.ProgressMessage) {

	currentVisitFrequency := sortedVisitFrequencies[0]
	indexAtNewVisitFrequency := 0
	globalProgression := 0
	
	for {
		rq := <- ch
		if strings.Compare(rq.Task, fpTracking.SEND_PROGRESS_INFORMATION) == 0 {

			isVisitFrequencyChanged := false

			if rq.VisitFrequency != currentVisitFrequency {
				indexAtNewVisitFrequency += lengths[currentVisitFrequency]
				currentVisitFrequency = rq.VisitFrequency
				log.Println("new visitFrequency :",currentVisitFrequency)
				isVisitFrequencyChanged = true
			}

			globalProgression = (indexAtNewVisitFrequency + rq.Index) * 100 / totalLength

			log.Println("progression :",globalProgression)

			//We lock the mutex in order to have a clean write access to progressInformationSession
			lock.Lock()
			progressInformationSession[userId].Progression = globalProgression
			if isVisitFrequencyChanged {
				progressInformationSession[userId].CurrentVisitFrequency = currentVisitFrequency
			}
			lock.Unlock()

		} else if strings.Compare(rq.Task, fpTracking.SEND_NEW_COMPUTED_POINTS) == 0 {

			//We lock the mutex in order to have a clean write access to progressInformationSession
			lock.Lock()
			progressInformationSession[userId].AverageTrackingTimeGraph = append(progressInformationSession[userId].AverageTrackingTimeGraph,
				rq.GraphPoints["averageTrackingTime"])
			progressInformationSession[userId].MaximumAverageTrackingTimeGraph = append(progressInformationSession[userId].MaximumAverageTrackingTimeGraph,
				rq.GraphPoints["maximumAverageTrackingTime"])
			progressInformationSession[userId].NbIdsFrequencyGraph = append(progressInformationSession[userId].NbIdsFrequencyGraph,
				rq.GraphPoints["nbIdsFrequency"])
			progressInformationSession[userId].OwnershipFrequencyGraph = append(progressInformationSession[userId].OwnershipFrequencyGraph,
				rq.GraphPoints["ownershipFrequency"])
			lock.Unlock()
		} else if strings.Compare(rq.Task, fpTracking.CLOSE_GOROUTINE) == 0 {

			globalProgression = 100
			//We lock the mutex in order to have a clean write access to progressInformationSession
			lock.Lock()
			progressInformationSession[userId].Progression = globalProgression
			progressInformationSession[userId].inProgress = false
			progressInformationSession[userId].ExecutingTime = roundPlus(time.Since(progressInformationSession[userId].creationDate).Seconds(),2)
			lock.Unlock()

			close(ch)
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

	//We do the request in the SQL database
	fingerprintManager := fpTracking.FingerprintManager{
		Number: number,
		Train:  train,
		MinNumberFpPerUser: minNbPerUser,
		DBInfo: dbInfos,
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
	go listenFpTrackingProgressChannel(totalLength, visitFrequencies, lengths,
		userId, progressChannel)


	for _, visitFrequency := range visitFrequencies {
		scenarioResult := fpTracking.ReplayScenarioParallelWithProgressInformation(test,
			visitFrequency, fpTracking.RuleBasedLinkingParallel, goroutineNumber,
			visitFrequencyToReplaySequence[visitFrequency], progressChannel)

		log.Println("We compute the results for visitFrequency",visitFrequency)
		progressChannel <- fpTracking.ProgressMessage {
			Task : fpTracking.SEND_NEW_COMPUTED_POINTS,
			GraphPoints : computeGraphicsPoints(fpTracking.AnalyseScenarioResultInStruct(visitFrequency, scenarioResult, test)),
		}
	}

	progressChannel <- fpTracking.ProgressMessage{Task : fpTracking.CLOSE_GOROUTINE}
}

//Returns whether the tracking alorithm is currently running or not for the user in argument
func isCurrentlyRunningForUser(id string) bool {
	//We only need a read access to the global variable progressInformationSession
	lock.RLock()
	value, is_present := progressInformationSession[id]
	result := is_present && value.inProgress
	lock.RUnlock()

	return result
}

//Delete userId previous session + all old sessions
func checkAndDeleteOldSessions(id string) {

	lock.Lock()
	for userId, progressInfo := range progressInformationSession {
		//We delete a session if it was created more than one day ago
		if time.Since(progressInfo.creationDate).Hours() >= float64(24) || userId == id {
			log.Println("userId",userId,"was deleted in the progressInformationSession map")
			delete(progressInformationSession,userId)
		}
	}
	lock.Unlock()
}

func computeGraphicsPoints(results fpTracking.ResultsForVisitFrequency) map[string]fpTracking.GraphicPoint {

	var nbRawDays []float64
	var maxChainRatio []float64

	nbAssignedIdsMean := 0.0

	for _, res1 := range results.Res1 {
		//we recompute ratio so that the first fingerprint of each assigned id is not counted
		res1.Ratio = float64((res1.NbOriginalFp - res1.NbAssignedIds)/res1.NbAssignedIds)

		nbRawDays = append(nbRawDays, res1.Ratio * float64(results.VisitFrequency))
		maxChainRatio = append(maxChainRatio, float64((res1.MaxChain - 1)*results.VisitFrequency))

		nbAssignedIdsMean += float64(res1.NbAssignedIds)
	}

	maxChainRatioMean := getMeanFromFloatSlice(maxChainRatio)
	nbRawDaysMean := getMeanFromFloatSlice(nbRawDays)

	nbAssignedIdsMean /= float64(len(results.Res1))

	ownershipMean := 0.0
	for _, res2 := range results.Res2 {
		ownershipMean += res2.Ownership
	}
	ownershipMean /= float64(len(results.Res2))


	m := make(map[string]fpTracking.GraphicPoint)
	m["averageTrackingTime"] = fpTracking.GraphicPoint{VisitFrequency : results.VisitFrequency, Value : roundPlus(nbRawDaysMean,4)}
	m["maximumAverageTrackingTime"] = fpTracking.GraphicPoint{VisitFrequency : results.VisitFrequency, Value : roundPlus(maxChainRatioMean,4)}
	m["nbIdsFrequency"] = fpTracking.GraphicPoint{VisitFrequency : results.VisitFrequency, Value : roundPlus(nbAssignedIdsMean,4)}
	m["ownershipFrequency"] = fpTracking.GraphicPoint{VisitFrequency : results.VisitFrequency, Value : roundPlus(ownershipMean,4)}

	return m
}

func getMeanFromFloatSlice(floatSlice []float64) float64 {
	mean := 0.0

	for _, fl := range floatSlice {
		mean += fl
	}

	return (mean / float64(len(floatSlice)))
}

func round(f float64) float64 {
    return math.Floor(f + .5)
}

func roundPlus(f float64, places int) (float64) {
	shift := math.Pow(10, float64(places))
	return round(f * shift) / shift;
}