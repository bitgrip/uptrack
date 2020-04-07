package job

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"regexp"
)

import "strings"

// Body is taken from plain body or transformed from base64 body if exists
func (j *UpJob) Body() io.Reader {
	if len(j.PlainBody) > 0 {
		return strings.NewReader(j.PlainBody)
	}
	if len(j.Base64Body) > 0 {
		decoded, _ := base64.StdEncoding.DecodeString(j.Base64Body)
		return strings.NewReader(string(decoded))
	}
	return nil
}
func (j *UpJob) HostString() (string, error) {
	if len(j.Host) > 0 {
		return j.Host, nil
	} else {
		url, err := url.Parse(j.URL)
		if err != nil {
			log.Panic(fmt.Sprintf("Fatal error: %s \n for url : %s in Job %s", err, j.URL, j.Name))
			return "", err
		}
		host := url.Host
		trimmedHost, _, err := net.SplitHostPort(host)
		if err == nil {
			host = trimmedHost
		}
		regex, _ := regexp.Compile("^www.")
		host = regex.ReplaceAllString(host, "")
		return host, nil
	}

}
