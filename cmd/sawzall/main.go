package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

var logFile string
var outputFile string
var interval time.Duration

func init() {
	flag.StringVar(&logFile, "logFile", "/var/log/nginx/access.log", "Nginx log file to read.")
	flag.StringVar(&outputFile, "outputFile", "/var/log/stats.log", "Output file to write statsd-compatible summaries")
	flag.DurationVar(&interval, "interval", 5*time.Second, "Interval to output summaries to output file")
}

func main() {
	flag.Parse()

	// open the output file for writing
	outputFileWriter, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}
	defer outputFileWriter.Close()

	sawzall, err := sawzall.NewSawzall(logFile, outputFileWriter)
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}
	defer sawzall.Stop()

	// create a quit channel and catch the sigterm signal to cause the program
	// to exit by closing the quit channel
	quitChan := make(chan struct{})
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		close(quitChan)
	}()

	fmt.Println("Press Ctrl+C to exit...")

	// main run loop, periodically write summary every 'interval' seconds and stop
	// if quit is given
	for {
		select {
		case <-time.After(interval):
			sawzall.WriteSummary()
		case <-quitChan:
			return
		}
	}
}
