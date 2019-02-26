package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DataWriter ...
type DataWriter struct {
	Ch  chan int
	lag time.Duration
}

// Start ...
func (w *DataWriter) Start() {
	w.setupInterruptHandler()
	w.lag = 50 * time.Millisecond
	for {
		data, ok := <-w.Ch
		if !ok {
			break
		}
		w.writeToDb(data)
	}
}

func (w *DataWriter) writeToDb(data int) {
	time.Sleep(w.lag)
	fmt.Printf("Writing %d to DB\n", data)
}

func (w *DataWriter) setupInterruptHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for i := 1; ; i++ {
			<-c
			fmt.Println("\r- Ctrl+C pressed in Terminal")
			switch {
			case i == 1:
				w.lag = 500 * time.Millisecond
			case i == 2:
				w.lag = 50 * time.Millisecond
			case i == 3:
				os.Exit(0)
			}
		}
	}()
}
