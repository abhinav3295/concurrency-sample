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
	ch  <-chan int
	lag time.Duration
	wg  sync.WaitGroup
}

// NewDataWriter ...
func NewDataWriter(ch <-chan int) *DataWriter {
	return &DataWriter{
		ch: ch,
	}
}

// Start ...
func (w *DataWriter) Start(workerCount int) {
	w.setupInterruptHandler()
	w.lag = 10 * time.Millisecond
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
		data, ok := <-w.ch
		if !ok {
			break
		}
		w.writeToDb(data)
	}
	w.wg.Done()
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
			switch i {
			case 1:
				w.lag = 1000 * time.Millisecond
			case 2:
				w.lag = 10000 * time.Millisecond
			case 3:
				w.lag = 10 * time.Millisecond
			}
		}
	}()
}
