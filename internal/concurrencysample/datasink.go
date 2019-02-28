package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// DbSink ...
type DbSink struct {
	workers int
	lag     time.Duration
	wg      sync.WaitGroup
}

// NewDbSink ...
func NewDbSink(workers int) *DbSink {
	return &DbSink{
		workers: workers,
	}
}

// Listen ...
func (w *DbSink) Listen(producer Producer) {
	w.setupInterruptHandler()
	w.lag = 10 * time.Millisecond
	w.wg.Add(w.workers)
	for i := 0; i < w.workers; i++ {
		go w.startWorker(producer.GetChannel())
	}
}

// WaitForFinish ...
func (w *DbSink) WaitForFinish() {
	w.wg.Wait()
}
func (w *DbSink) startWorker(ch <-chan int) {
	for {
		data, ok := <-ch
		if !ok {
			break
		}
		w.writeToDb(data)
	}
	w.wg.Done()
}

func (w *DbSink) writeToDb(data int) {
	time.Sleep(w.lag)
	fmt.Printf("Writing %d to DB\n", data)
}

func (w *DbSink) setupInterruptHandler() {
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
