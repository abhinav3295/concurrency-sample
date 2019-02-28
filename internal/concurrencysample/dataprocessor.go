package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DataProcessor ...
type DataProcessor struct {
	ch   chan<- int
	stop bool
}

// NewDataProcessor ...
func NewDataProcessor(ch chan<- int) *DataProcessor {
	return &DataProcessor{
		ch: ch,
	}
}

// Run ...
func (d *DataProcessor) Run() {
	d.stop = false
	d.setupInterruptHandler()
	i := 0
	for !d.stop {
		d.doSomething(i)
		select {
		case d.ch <- i:
		default:
			d.writeToRedis(i)
		}
		i = i + 1
	}
	close(d.ch)
}

func (d *DataProcessor) doSomething(i int) int {
	time.Sleep(100 * time.Millisecond)
	j := i
	return j
}

func (d *DataProcessor) writeToRedis(data int) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Write %d to REDIS\n", data)
}

func (d *DataProcessor) setupInterruptHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for i := 1; ; i++ {
			<-c
			if i >= 3 {
				d.stop = true
			}
		}
	}()
}
