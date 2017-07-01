package sawzall

import (
	"fmt"
	"regexp"
	"strconv"
)

type Parser interface {
	ParseString(line string) (*LogEntry, error)
}

type NginxParser struct {
	regex *regexp.Regexp
}

// NewNginxParser returns an NginxParser which parses strings based on a log
// format.
func NewNginxParser() *NginxParser {
	// For the following log format:
	// log_format combined $remote_addr - $http_x_forwarded_for - $http_x_realip -
	// [$time_local] $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme
	// "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent";

	// This matches the request and status part ...
	regex := regexp.MustCompile(`"\w{3,6} (.*) \w{0,4}/\d\.\d" (\d{3})`)
	return &NginxParser{regex: regex}
}

func (n *NginxParser) ParseString(line string) (*LogEntry, error) {
	fields := n.regex.FindStringSubmatch(line)
	if fields == nil {
		return nil, fmt.Errorf("Could not parse log line '%v'", line)
	}
	code, err := strconv.Atoi(fields[2])
	if err != nil {
		return nil, err
	}
	route := fields[1]
	return &LogEntry{Code: code,
		Route: route,
	}, nil
}
