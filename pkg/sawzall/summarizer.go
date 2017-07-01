package sawzall

import (
	"sync"

	"gopkg.in/tomb.v2"
)

type Summary map[string]int64

type Summarizer interface {
	GetAndReset() Summary
	Stop() error
}

type AggregateSummarizer struct {
	logChan chan *LogEntry
	results Summary
	mu      sync.Mutex
	t       tomb.Tomb
}

// NewAggregateSummarizer takes a log entry channel returns an
// AggregateSummarizer which summarizes its values.
func NewAggregateSummarizer(logChan chan *LogEntry) *AggregateSummarizer {
	a := &AggregateSummarizer{
		logChan: logChan,
		results: Summary{},
	}

	a.t.Go(a.summarize)

	return a
}

// summarize does the processing. It is internal to the AggregateSummarizer.
func (a *AggregateSummarizer) summarize() error {
	for {
		select {
		case entry, ok := <-a.logChan:
			// if the channel was closed, nothing to do, just return
			if !ok {
				return nil
			}
			a.mu.Lock()
			if entry.Code >= 600 {
			} else if entry.Code >= 500 {
				a.results["50x"]++
				a.results[entry.Route]++
			} else if entry.Code >= 400 {
				a.results["40x"]++
			} else if entry.Code >= 300 {
				a.results["30x"]++
			} else if entry.Code >= 200 {
				a.results["20x"]++
			}
			a.mu.Unlock()
		case <-a.t.Dying():
			return nil
		}
	}
}

// GetAndReset returns a summary of the results collected thus far and resets
// the results values.
func (a *AggregateSummarizer) GetAndReset() Summary {
	a.mu.Lock()
	defer a.mu.Unlock()

	snapshot := Summary{}
	for k, v := range a.results {
		snapshot[k] = v
	}
	a.results = Summary{}

	return snapshot
}

// Stop stops the summarize loop
func (a *AggregateSummarizer) Stop() error {
	a.t.Kill(nil)
	return a.t.Wait()
}
