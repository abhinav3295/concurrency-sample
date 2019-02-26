package concurrencysample

import (
	"fmt"
	"time"
)

// DataWriter ...
type DataWriter struct {
	Ch chan int
}

// Start ...
func (w *DataWriter) Start() {
	for {
		data, ok := <-w.Ch
		if !ok {
			break
		}
		w.writeToDb(data)
	}
}

func (w *DataWriter) writeToDb(data int) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Writing %d to DB\n", data)
}
