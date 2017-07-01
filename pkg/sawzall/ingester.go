package sawzall

import (
	"log"
	"os"

	"github.com/hpcloud/tail"
	"gopkg.in/tomb.v2"
)

type Ingester interface {
	Stop() error
}

type LogIngester struct {
	logFile string
	parser  Parser
	logChan chan *LogEntry
	t       tomb.Tomb
	tailer  *tail.Tail
}

// NewLogIngester takes a log file path, a parser and a log entry channel and
// returns a LogIngester which sends log entries into the log entry channel.
func NewLogIngester(logFile string, parser Parser, logChan chan *LogEntry) (*LogIngester, error) {
	seekInfo := tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}

	tailer, err := tail.TailFile(logFile, tail.Config{
		Follow:    true,
		Poll:      true,
		ReOpen:    true,
		MustExist: true,
		Location:  &seekInfo,
		Logger:    tail.DiscardingLogger,
	})

	if err != nil {
		return nil, err
	}

	ing := &LogIngester{
		logFile: logFile,
		parser:  parser,
		logChan: logChan,
		tailer:  tailer,
	}

	ing.t.Go(ing.ingestLoop)

	return ing, nil
}

// ingestLoop does the processing. It is internal to the LogIngester.
func (l *LogIngester) ingestLoop() error {
	for {
		select {
		case line, ok := <-l.tailer.Lines:
			if !ok {
				return nil
			}
			entry, err := l.parser.ParseString(line.Text)
			if err != nil {
				log.Printf("%v. Continuing...\n", err)
				continue
			}
			l.logChan <- entry
		case <-l.t.Dying():
			return nil
		}
	}
}

// Stop stops the ingest loop
func (l *LogIngester) Stop() error {
	l.t.Kill(nil)
	return l.t.Wait()
}
