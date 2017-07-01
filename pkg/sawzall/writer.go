package sawzall

import (
	"fmt"
	"io"
)

type Writer interface {
	Write(results Summary) error
}

type StatsDWriter struct {
	writer io.Writer
}

// NewStatsDWriter takes a Writer interface and returns a StatsDWriter which
// writes statd-compatible summaries to the writer.
func NewStatsDWriter(writer io.Writer) *StatsDWriter {
	return &StatsDWriter{
		writer: writer,
	}
}

// Write writes summary results in a statsd-compatible format
func (w *StatsDWriter) Write(results Summary) error {
	for k, v := range results {
		_, err := fmt.Fprintf(w.writer, "%s:%d|s\n", k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
