package job

import (
	"encoding/base64"
	"io"
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
