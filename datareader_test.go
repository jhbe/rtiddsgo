package rtiddsgo_test
/*
import (
	"testing"
	"rtiddsgo"
)

func TestCreateDataReader(t *testing.T) {
	p, _ := rtiddsgo.New(33)
	defer p.Free()

	topic, _ := p.CreateTopic("MyMessage", "")
	defer topic.Free()
	sub, _ := p.CreateSubscriber()

	dr, err := rtiddsgo.CreateDataReader(sub, topic, "", "", func() {})
	if err != nil {
		t.Error(err)
	}
	defer dr.Free()
}
*/