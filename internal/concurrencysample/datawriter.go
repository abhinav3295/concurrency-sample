package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// DataWriter ...
type DataWriter struct {
	Ch  chan int
	Lag time.Duration
	wg  sync.WaitGroup
}

// Start ...
func (w *DataWriter) Start(workerCount int) {
	w.setupInterruptHandler()
	w.Lag = 10 * time.Millisecond
	w.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go w.startWorker()
	}
}

// WaitForFinish ...
func (w *DataWriter) WaitForFinish() {
	w.wg.Wait()
}
func (w *DataWriter) startWorker() {
	for {
		data, ok := <-w.Ch
		if !ok {
			break
		}
		w.writeToDb(data)
	}
	w.wg.Done()
}

func (w *DataWriter) writeToDb(data int) {
	time.Sleep(w.Lag)
	fmt.Printf("Writing %d to DB\n", data)
}

func (w *DataWriter) setupInterruptHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for i := 1; ; i++ {
			<-c
			switch {
			case i == 1:
				w.Lag = 1000 * time.Millisecond
			case i == 2:
				w.Lag = 10000 * time.Millisecond
			case i == 3:
				w.Lag = 10 * time.Millisecond
			}
		}
	}()
}
