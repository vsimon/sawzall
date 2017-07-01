package sawzall_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

func TestSawzall(t *testing.T) {
	// create a buffer for the output
	output := &bytes.Buffer{}

	// create a temporary file in the testdata dir for representing the log file
	testFilePath := "testdata/.access.log"
	testFile, err := os.Create(testFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFilePath)
	defer testFile.Close()

	// create sawzall with the test log file path
	s, err := sawzall.NewSawzall(testFilePath, output)
	if err != nil {
		t.Fatal(err)
	}

	// write well-known sample log data to the test log file
	sampleFile, err := os.Open("testdata/sample.log")
	if err != nil {
		t.Fatal(err)
	}
	defer sampleFile.Close()
	_, err = io.Copy(testFile, sampleFile)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	s.WriteSummary()

	expectedCases := []string{
		"20x:3551|s",
		"40x:155|s",
		"30x:40|s",
	}

	for _, ec := range expectedCases {
		if strings.Contains(ec, output.String()) {
			t.Errorf("expected output %v to contain %v", output.String(), ec)
		}
	}

	err = s.Stop()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

func TestInvalidFileSawzall(t *testing.T) {
	// create a buffer for the output
	output := &bytes.Buffer{}
	_, err := sawzall.NewSawzall("", output)
	if err == nil {
		t.Errorf("expected error to be not nil, got %v", err)
	}
}
