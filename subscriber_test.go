package rtiddsgo

import (
	"testing"
)

func TestCreateSubscriber(t *testing.T) {
	p, _ := New(33)
	defer p.Free()

	_, err := p.CreateSubscriber()
	if err != nil {
		t.Error(err)
	}
}