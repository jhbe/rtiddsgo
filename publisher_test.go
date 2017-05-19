package rtiddsgo

import (
	"testing"
)

func TestCreatePublisher(t *testing.T) {
	p, _ := New(33)
	defer p.Free()

	_, err := p.CreatePublisher("", "")
	if err != nil {
		t.Error(err)
	}
}