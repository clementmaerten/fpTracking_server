package main

import (
	"log"
	"strings"
	"github.com/clementmaerten/fpTracking"
)

//This function listen to the progress channel and update the user's session with progress information
//This function is supposed to be executed by a goroutine
func listenFpTrackingProgressChannel(totalLength int, sortedVisitFrequencies []int, lengths map[int]int,
	ch <- chan fpTracking.ProgressMessage) {

	currentVisitFrequency := sortedVisitFrequencies[0]
	indexAtNewVisitFrequency := 0
	globalProgression := 0
	
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

		} else if strings.Compare(rq.Task, fpTracking.CLOSE_GOROUTINE) == 0 {
			log.Println("progression : 100")
			return
		} else {
			//This case should never happen
			log.Println("Wrong task for listenFpTrackingProgressChannel")
			return
		}
	}
}