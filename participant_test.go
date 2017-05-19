package rtiddsgo

import "testing"

func TestParticipant(t *testing.T) {
	p, err := New(33, "", "")
	defer p.Free()
	if err != nil {
		t.Error(err)
	}
}
