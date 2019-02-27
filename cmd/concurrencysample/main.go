package main

import (
	"github.com/abhinav3295/go-meetups/internal/concurrencysample"
)

func main() {
	ch := make(chan int, 4)
	writer := concurrencysample.DataWriter{
		Ch: ch,
	}
	app := concurrencysample.DataProcessor{
		Ch: ch,
	}
	go writer.Start(2)
	app.Run()
	writer.WaitForFinish()
}
