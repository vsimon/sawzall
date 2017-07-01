package sawzall_test

import (
	"reflect"
	"testing"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

func TestSummarizer(t *testing.T) {
	// create a log channel and summarizer
	logChan := make(chan *sawzall.LogEntry)
	s := sawzall.NewAggregateSummarizer(logChan)

	// feed the log channel a log entry
	logChan <- &sawzall.LogEntry{-1, "/"}
	logChan <- &sawzall.LogEntry{0, "/"}
	logChan <- &sawzall.LogEntry{99, "/"}
	logChan <- &sawzall.LogEntry{200, "/"}
	logChan <- &sawzall.LogEntry{200, "/"}
	logChan <- &sawzall.LogEntry{300, "/"}
	logChan <- &sawzall.LogEntry{301, "/"}
	logChan <- &sawzall.LogEntry{403, "/"}
	logChan <- &sawzall.LogEntry{404, "/"}
	logChan <- &sawzall.LogEntry{500, "/login"}
	logChan <- &sawzall.LogEntry{501, "/login"}
	logChan <- &sawzall.LogEntry{599, "/logout"}
	logChan <- &sawzall.LogEntry{600, "/"}
	logChan <- &sawzall.LogEntry{999, "/"}

	// get and check expected results
	summary := s.GetAndReset()
	expected := sawzall.Summary{"20x": 2, "30x": 2, "40x": 2, "50x": 3,
		"/login": 2, "/logout": 1}
	if !reflect.DeepEqual(summary, expected) {
		t.Errorf("expected summary to be %v, got %v", expected, summary)
	}

	err := s.Stop()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}
