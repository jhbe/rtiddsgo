package rtiddsgo
/*
import (
	"testing"
	"rtiddsgo/test"
)

func TestCreateDataWriter(t *testing.T) {
	p, _ := New(33)
	defer p.Free()

	test.Message_RegisterType(p)
	topic, _ := p.CreateTopic("MyMessage", test.Message_GetTypeName())
	defer topic.Free()
	pub, _ := p.CreatePublisher()

	_, err := CreateDataWriter(pub, topic, "", "")
	if err != nil {
		t.Error(err)
	}
}
*/