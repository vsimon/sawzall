package sawzall_test

import (
	"testing"

	"github.com/vsimon/sawzall/pkg/sawzall"
)

func TestWellFormedLogs(t *testing.T) {
	p := sawzall.NewNginxParser()

	// a couple test-cases for well-formed but different codes and routes within
	// the logs
	testCases := []struct {
		line     string
		expected *sawzall.LogEntry
	}{
		{`10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" 500 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`,
			&sawzall.LogEntry{500, "/our-products"},
		},
		{`10.10.180.40 - 60.173.10.235 - - - [02/Aug/2015:21:36:33 +0000]  http http http "GET /cht/ HTTP/1.1" 404 162 "-" "Mozilla/4.0 (compatible; Win32; WinHttp.WinHttpRequest.5)"`,
			&sawzall.LogEntry{404, "/cht/"},
		},
		{`10.10.180.161 - 180.76.15.140, 192.33.28.238 - - - [02/Aug/2015:20:27:44 +0000]  http http http "GET /page/privacy-policy HTTP/1.1" 301 178 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"`,
			&sawzall.LogEntry{301, "/page/privacy-policy"},
		},
		{`10.10.180.161 - 50.112.166.232, 192.33.28.238 - - - [02/Aug/2015:15:56:14 +0000]  https https https "GET /our-products HTTP/1.1" 200 35967 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`,
			&sawzall.LogEntry{200, "/our-products"},
		},
		{`10.10.180.40 - 60.173.10.235 - - - [02/Aug/2015:21:36:33 +0000]  http http http "GET /cht/ HTTP/1.1" 404 162 "-" "Mozilla/4.0 (compatible; Win32; WinHttp.WinHttpRequest.5)"`,
			&sawzall.LogEntry{404, "/cht/"},
		},
		{`50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:18:56:27 +0000]  http https,http https,http "POST /login HTTP/1.1" 200 6058 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`,
			&sawzall.LogEntry{200, "/login"},
		},
		{`50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:18:56:27 +0000]  http https,http https,http "POST /guîlë HTTP/1.1" 200 6058 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`,
			&sawzall.LogEntry{200, "/guîlë"},
		},
		{`50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:18:56:27 +0000]  http https,http https,http "PUT / HTTP/1.1" 200 6058 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`,
			&sawzall.LogEntry{200, "/"},
		},
		{`172.19.0.1 - - - - - [01/Jul/2017:05:46:49 +0000] http - http "GET / HTTP/1.1" 200 612 "-" "curl/7.51.0"`,
			&sawzall.LogEntry{200, "/"},
		},
	}

	for _, tc := range testCases {
		logEntry, err := p.ParseString(tc.line)
		if err != nil {
			t.Errorf("expected error to be nil, got %v", err)
			continue
		}
		if *logEntry != *tc.expected {
			t.Errorf("expected logEntry to be %v, got %v", tc.expected, logEntry)
		}
	}
}

func TestMalformedLogs(t *testing.T) {
	p := sawzall.NewNginxParser()

	// a couple test-cases with empty strings, malformed and missing content
	// within the logs
	testCases := []struct {
		line string
	}{
		{``},
		{`""`},
		{`10.10.180.161 - 180.76.15.140, 192.33.28.238 - - - [02/Aug/2015:20:27:44 +0000]  http http http "GET /page HTTP/1.1" abc 178 "-" "Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)"`},
		{`50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:18:56:27 +0000]  http https,http https,http "GET HTTP/1.1" 200 6058 "-" "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.81 Safari/537.36"`},
	}

	for _, tc := range testCases {
		_, err := p.ParseString(tc.line)
		if err == nil {
			t.Errorf("expected error to not be nil, got %v", err)
		}
	}
}
