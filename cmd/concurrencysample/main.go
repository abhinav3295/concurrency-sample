package main

import (
	"github.com/abhinav3295/go-meetups/internal/concurrencysample"
)

func main() {
	ch := make(chan int, 4)
	sink := concurrencysample.NewDbSink(ch)
	producer := concurrencysample.NewDataProducer(ch)

	go sink.Start(2)

	producer.Run()
	sink.WaitForFinish()
}
