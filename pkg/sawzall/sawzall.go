package sawzall

import (
	"io"
	"log"
)

type LogEntry struct {
	Code  int
	Route string
}

type Sawzall struct {
	ingester   Ingester
	summarizer Summarizer
	writer     Writer
	logChan    chan *LogEntry
}

// create a new sawzall. this will immediately start ingesting logs. call Stop()
// to stop ingestion.
func NewSawzall(logFile string, outputFileWriter io.Writer) (*Sawzall, error) {
	// create a new log channel, it will be used to pass log entries to the
	// summarizer
	logChan := make(chan *LogEntry)
	ing, err := NewLogIngester(logFile, NewNginxParser(), logChan)
	if err != nil {
		return nil, err
	}
	return &Sawzall{
		ingester:   ing,
		summarizer: NewAggregateSummarizer(logChan),
		writer:     NewStatsDWriter(outputFileWriter),
		logChan:    logChan,
	}, nil
}

func (s *Sawzall) WriteSummary() {
	r := s.summarizer.GetAndReset()
	err := s.writer.Write(r)
	if err != nil {
		log.Fatal(err)
	}
}

// Stop stops all processing
func (s *Sawzall) Stop() error {
	err := s.ingester.Stop()
	if err != nil {
		return err
	}
	close(s.logChan)
	err = s.summarizer.Stop()
	if err != nil {
		return err
	}
	return nil
}
