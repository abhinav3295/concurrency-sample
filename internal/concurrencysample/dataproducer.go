package concurrencysample

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Producer ...
type Producer interface {
	GetChannel() <-chan int
}

// DataProducer ...
type DataProducer struct {
	ch   chan int
	stop int32
}

// GetChannel returns the communication channel as read-only
func (d *DataProducer) GetChannel() <-chan int {
	return d.ch
}

// NewDataProducer ...
func NewDataProducer(buffer int) *DataProducer {
	ch := make(chan int, buffer)
	return &DataProducer{
		ch: ch,
	}
}

// Run ...
func (d *DataProducer) Run() {
	d.stop = 0
	d.setupInterruptHandler()
	i := 0

	go func() {
		defer close(d.ch)
		for d.stop == 0 {
			d.doSomething(i)
			d.ch <- i
			i = i + 1
		}
	}()
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
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for i := 1; ; i++ {
			<-c
			if i >= 3 {
				atomic.CompareAndSwapInt32(&d.stop, 0, 1)
			}
		}
	}()
}
