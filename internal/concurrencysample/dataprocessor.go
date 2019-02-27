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
	Ch          chan int
	StopChannel chan bool
}

// Run ...
func (d *DataProcessor) Run() {
	d.StopChannel = make(chan bool, 1)
	d.setupInterruptHandler()
	stop := false
	i := 0
	for !stop {
		d.doSomething(i)
		select {
		case stop = <-d.StopChannel:
			if stop {
				close(d.Ch)
			}

		case d.Ch <- i:
		default:
			d.writeToKafka(i)
		}
		i = i + 1
	}
}

func (d *DataProcessor) doSomething(i int) int {
	time.Sleep(100 * time.Millisecond)
	j := i
	return j
}

func (d *DataProcessor) writeToKafka(data int) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Write %d to KAFKA\n", data)
}

func (d *DataProcessor) setupInterruptHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for i := 1; ; i++ {
			<-c
			switch {
			case i >= 4:
				d.StopChannel <- true
			}
		}
	}()
}
