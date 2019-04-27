package main

import (
	"rtiddsgo/example/src"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds)

	type info struct {
		alive bool
		message example.Com_Ex_A
	}
	rxChan := make(chan info)
	r := example.NewReaderCom_Ex_A(33, "TheA", "", "", "", "", func(alive bool, data example.Com_Ex_A) {
		rxChan <- info{alive, data} // Send the message to the main thread.
	})
	defer r.Free()

	ctrlCChan := make(chan os.Signal, 1)
	signal.Notify(ctrlCChan, os.Interrupt)

	log.Println("Waiting for data. CTRL_C to stop.")
	for keepGoing := true; keepGoing; {
		select {
		case mess := <-rxChan:
			if mess.alive {
				log.Println(mess.message)
			} else {
				log.Println("Instance no longer alive")
			}
		case <-ctrlCChan:
			log.Println("Stopping.")
			keepGoing = false
		}
	}

	log.Println("Stopped.")
}

