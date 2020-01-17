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
func (job *UpJob) Body() io.Reader {
	if len(job.PlainBody) > 0 {
		return strings.NewReader(job.PlainBody)
	}
	if len(job.Base64Body) > 0 {
		decoded, _ := base64.StdEncoding.DecodeString(job.Base64Body)
		return strings.NewReader(string(decoded))
	}
	return nil
}
func (job *UpJob) HostString() (string, error) {
	if len(job.Host) > 0 {
		return job.Host, nil
	} else {
		url, err := url.Parse(job.URL)
		if err != nil {
			log.Panic(fmt.Sprintf("Fatal error: %s \n for url : %s in Job %s", err, job.URL, job.Name))
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
