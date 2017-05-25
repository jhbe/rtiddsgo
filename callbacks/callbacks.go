package callbacks

import "fmt"

type Callback func()
type Callbacks struct {
	callbacks map[int]Callback
	next int
}

func New() Callbacks {
	return Callbacks{callbacks: make(map[int]Callback)}
}

func (cbs *Callbacks)Add(callback Callback) int {
	index := cbs.next
	cbs.next++
	cbs.callbacks[index] = callback
	return index
}

func (cbs Callbacks)Invoke(i int) error {
	fn, exist := cbs.callbacks[i]
	if !exist {
		return fmt.Errorf("Index %d does not hold a callback function.", i)
	}
	fn()
	return nil
}
