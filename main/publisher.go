package main

import (
	"log"
	"rtiddsgo"
	"rtiddsgo/example"
	"time"
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

	pub, err := p.CreatePublisher("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Free()

	dw, err := example.NewMyModule_MyMessageDataWriter(pub, topic, "", "")
	if err != nil {
		log.Fatal(err)
	}
	defer dw.Free()

	m := example.MyModule_MyMessage{}
	m.TheShort = 20
	m.TheUnsignedShort = 21
	m.TheLong = 22
	m.TheUnsignedLong = 23
	m.TheFloat = 24.0
	m.TheDouble = 25.0
	m.TheBool = true
	m.Text = "JHBE WAS HERE!!!"
	m.FiveEnums = []example.MyModule_MyEnum{example.MyEnum_One, example.MyEnum_Two, example.MyEnum_One}
	m.TheUnion.MyModule_MyUnion_D = example.MyEnum_Two
	m.TheUnion.TheError.TheBool = true
	m.TheUnion.TheError.Description = "This part a union with discriminant MyEnum_Two."

	for i := 0; i < 10; i++ {
		dw.Write(m)
		log.Println("Wrote data:", m)
		m.TheLong++

		time.Sleep(time.Second)
	}

	log.Println("Stopped.")
}
