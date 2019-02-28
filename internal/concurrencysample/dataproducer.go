package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// DataProducer ...
type DataProducer struct {
	ch   chan<- int
	stop bool
}

// NewDataProducer ...
func NewDataProducer(ch chan<- int) *DataProducer {
	return &DataProducer{
		ch: ch,
	}
}

// Run ...
func (d *DataProducer) Run() {
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

func (d *DataProducer) doSomething(i int) int {
	time.Sleep(100 * time.Millisecond)
	j := i
	return j
}

func (d *DataProducer) writeToRedis(data int) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Write %d to REDIS\n", data)
}

func (d *DataProducer) setupInterruptHandler() {
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
