package rtiddsgo

import "C"

//export OnDataAvailable
func OnDataAvailable(index int) {
	onDataAvailableCallbacks.Invoke(index)
}