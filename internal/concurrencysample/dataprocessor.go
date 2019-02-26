package concurrencysample

import (
	"fmt"
	"time"
)

// DataProcessor ...
type DataProcessor struct {
	Ch chan int
}

// Run ...
func (d *DataProcessor) Run() {
	i := 0
	for {
		i = d.doSomething(i)
		select {
		case d.Ch <- i:
		default:
			d.writeToKafka(i)
		}
	}
}

func (d *DataProcessor) doSomething(i int) int {
	time.Sleep(100 * time.Millisecond)
	j := i + 1
	return j
}

func (d *DataProcessor) writeToKafka(data int) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Write %d to KAFKA\n", data)
}
