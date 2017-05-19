package main

import (
	"log"
	"rtiddsgo"
	"rtiddsgo/example"
	"os"
	"os/signal"
)

func main() {
	p, err := rtiddsgo.New(33, "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer p.Free()

	err = example.MyModule_MyMessage_RegisterType(p)
	if err != nil {
		log.Fatal(err)
	}

	topic, err := p.CreateTopic("MyMessage", example.MyModule_MyMessage_GetTypeName(), "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer topic.Free()

	sub, err := p.CreateSubscriber("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer sub.Free()

	rxChan := make(chan example.MyModule_MyMessage)
	dr, err := example.NewMyModule_MyMessageDataReader(sub, topic, "", "", func(message example.MyModule_MyMessage) {
		// Send the message to the main thread.
		rxChan <- message
	})
	if err != nil {
		log.Fatal(err)
	}
	defer dr.Free()

	ctrlCChan := make(chan os.Signal, 1)
	signal.Notify(ctrlCChan, os.Interrupt)

	log.Println("Waiting for data. CTRL_C to stop.")
	for keepGoing := true; keepGoing; {
		select {
		case mess := <-rxChan:
			log.Println(mess)
		case <-ctrlCChan:
			log.Println("Stopping.")
			keepGoing = false
		}
	}

	log.Println("Stopped.")
}

