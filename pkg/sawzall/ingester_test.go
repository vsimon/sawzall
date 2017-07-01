package sawzall_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

func TestIngesterCanCreateLogEntries(t *testing.T) {
	ingTest, teardown := NewIngesterTest("can-create-log-entries", t)
	defer teardown()

	// create a log channel
	logChan := make(chan *sawzall.LogEntry)

	// create a log file for input
	filePath := ingTest.CreateFile("access.log", ``)

	ing, err := sawzall.NewLogIngester(filePath, sawzall.NewNginxParser(), logChan)
	if err != nil {
		t.Fatal(err)
	}

	// write test case data to log file
	ingTest.AppendFile("access.log", `10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" 500 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"
10.10.180.40 - 60.173.10.235 - - - [02/Aug/2015:21:36:33 +0000]  http http http "GET /cht/ HTTP/1.1" 404 162 "-" "Mozilla/4.0 (compatible; Win32; WinHttp.WinHttpRequest.5)"
10.10.180.161 - 180.76.15.140, 192.33.28.238 - - - [02/Aug/2015:20:27:44 +0000]  http http http "GET /page/privacy-policy HTTP/1.1" 301 178 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"
10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" 200 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"
10.10.180.40 - 60.173.10.235 - - - [02/Aug/2015:21:36:33 +0000]  http http http "GET /cht/ HTTP/1.1" 404 162 "-" "Mozilla/4.0 (compatible; Win32; WinHttp.WinHttpRequest.5)"
50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:18:56:27 +0000]  http https,http https,http "POST /login HTTP/1.1" 200 6058 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"
`)

	expectedCases := []*sawzall.LogEntry{
		&sawzall.LogEntry{500, "/our-products"},
		&sawzall.LogEntry{404, "/cht/"},
		&sawzall.LogEntry{301, "/page/privacy-policy"},
		&sawzall.LogEntry{200, "/our-products"},
		&sawzall.LogEntry{404, "/cht/"},
		&sawzall.LogEntry{200, "/login"},
	}

	for _, ec := range expectedCases {
		select {
		case logEntry := <-logChan:
			if *logEntry != *ec {
				t.Errorf("expected logEntry to be %v, got %v", *ec, *logEntry)
			}
		case <-time.After(time.Second):
			t.Fatal("timed out expecting a log entry")
		}
	}

	err = ing.Stop()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

func TestIngesterCanIgnoreInvalidLogEntries(t *testing.T) {
	ingTest, teardown := NewIngesterTest("can-ignore-invalid-log-entries", t)
	defer teardown()

	// create a log channel
	logChan := make(chan *sawzall.LogEntry)

	// create a log file for input
	filePath := ingTest.CreateFile("access.log", ``)

	ing, err := sawzall.NewLogIngester(filePath, sawzall.NewNginxParser(), logChan)
	if err != nil {
		t.Fatal(err)
	}

	// write test case data to log file
	ingTest.AppendFile("access.log", `"GET Mozilla/5.0 (X11; Linux x86_64)"
	""

	abcdefghi
	1 2 3 4 5 6 7 8 9
	😍 🙇 💁 🙅 🙆 🙋 🙎 🙍
	`)

	select {
	case <-logChan:
		t.Errorf("expected logEntry to not be created")
	case <-time.After(time.Second):
		break
	}

	err = ing.Stop()
	if err != nil {
		t.Errorf("expected error to be nil, got %v", err)
	}
}

// test helpers

type IngesterTest struct {
	path string
	*testing.T
}

func NewIngesterTest(name string, t *testing.T) (IngesterTest, func()) {
	// create a top-level .test dir to house the test files
	i := IngesterTest{".test/" + name, t}
	err := os.MkdirAll(i.path, os.ModeTemporary|0700)
	if err != nil {
		i.Fatal(err)
	}

	// provide a teardown function to remove temporary test dir/files
	teardown := func() {
		err := os.RemoveAll(".test")
		if err != nil {
			i.Fatal(err)
		}
	}

	return i, teardown
}

func (i IngesterTest) CreateFile(name string, contents string) string {
	filePath := i.path + "/" + name
	err := ioutil.WriteFile(filePath, []byte(contents), 0600)
	if err != nil {
		i.Fatal(err)
	}
	return filePath
}

func (i IngesterTest) AppendFile(name string, contents string) {
	f, err := os.OpenFile(i.path+"/"+name, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		i.Fatal(err)
	}
	defer f.Close()
	_, err = f.WriteString(contents)
	if err != nil {
		i.Fatal(err)
	}
}
