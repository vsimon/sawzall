package sawzall_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

func TestStatsDWriter(t *testing.T) {
	// create a buffer for the output
	output := &bytes.Buffer{}
	w := sawzall.NewStatsDWriter(output)

	summary := sawzall.Summary{"20x": 2, "30x": 2, "40x": 2, "50x": 3,
		"/login": 2, "/logout": 1}
	expected := []string{"20x:2|s", "30x:2|s", "40x:2|s", "50x:3|s", "/login:2|s",
		"/logout:1|s"}

	w.Write(summary)

	// check if the resulting output contains expected statsd lines
	for _, line := range expected {
		if !strings.Contains(output.String(), line) {
			t.Errorf("expected output to contain %v, got %v", line, output.String())
		}
	}
}
