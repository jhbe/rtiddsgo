package main

import (
	"rtiddsgo/example/src"
	"log"
	"time"
)

func main() {
	w := example.NewWriterCom_Ex_A(33, "TheA", "", "", "", "")
	defer w.Free()

	m := example.Com_Ex_A{0, "Hello"}
	for i := 0; i < 10; i++ {
		m.TheLong++
		w.Dw.Write(m)
		log.Println("Wrote data:", m)
		time.Sleep(time.Second)
	}
}

