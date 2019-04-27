package main

import (
	"rtiddsgo/verification"
	"rtiddsgo/verification/src"
	"time"
)

func main() {
	w := eb.NewWriterCom_Two_X(0, verification.TopicName, "", "", "", "")
	defer w.Free()
	time.Sleep(3 * time.Second) // Let the DDS discovery process settle...

	w.Dw.Write(verification.ComTwoX)
}
