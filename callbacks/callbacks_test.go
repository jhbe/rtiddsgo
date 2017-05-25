package callbacks

import "testing"

func TestEmptySet(t *testing.T) {
	cbs := New()
	err := cbs.Invoke(0)
	if err == nil {
		t.Error("Expected to fail calling an empty set, but did not.")
	}
}

func TestOneFunction(t *testing.T) {
	invoked := false
	cbs := New()
	index := cbs.Add(func () {
		invoked = true
	})
	err := cbs.Invoke(index)
	if err != nil {
		t.Error("Expected to not fail calling an empty set, but did not.")
	}
	if !invoked {
		t.Error("Expected the invoked boolean to be set to true by the callback, but it remained false.")
	}
}

func TestTwoFunctions(t *testing.T) {
	counter := 0
	cbs := New()
	index_1 := cbs.Add(func () {
		counter += 1
	})
	index_2 := cbs.Add(func () {
		counter += 2
	})
	err := cbs.Invoke(index_1)
	if err != nil {
		t.Error("Expected to not fail calling an empty set, but did not.")
	}
	err = cbs.Invoke(index_2)
	if counter != 3 {
		t.Error("Expected the counter to be set to 1+2 by the callbacks, but it wasn't.")
	}
}