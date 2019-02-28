package main

import (
	"github.com/abhinav3295/go-meetups/internal/concurrencysample"
)

func main() {
	producer := concurrencysample.NewDataProducer(4)
	sink := concurrencysample.NewDbSink(2)

	sink.Listen(producer)
	producer.Run()

	sink.WaitForFinish()
}
