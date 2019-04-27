package main

import (
	"log"
	"os"
	"os/signal"
	"reflect"
	"rtiddsgo/verification"
	"rtiddsgo/verification/src"
)

func main() {
	rxChan := make(chan bool)
	r := eb.NewReaderCom_Two_X(0, verification.TopicName, "", "", "", "", func(alive bool, data eb.Com_Two_X) {
		if alive {
			if reflect.DeepEqual(verification.ComTwoX, data) {
				log.Printf("Verified OK")
				rxChan <- true
			} else {
				log.Printf("Expected\n%q\n, got\n%q\n", verification.ComTwoX, data)
				rxChan <- false
			}
		} else {
			rxChan <- false
		}
	})
	defer r.Free()

	ctrlCChan := make(chan os.Signal, 1)
	signal.Notify(ctrlCChan, os.Interrupt)

	for keepGoing := true; keepGoing; {
		select {
		case result := <-rxChan:
			if !result {
				os.Exit(-1)
			}
			log.Printf("Stopping...")
			keepGoing = false
		case <-ctrlCChan:
			keepGoing = false
		}
	}
	log.Printf("Stopped.")
}
