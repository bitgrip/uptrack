package job

import (
	"encoding/base64"
	"fmt"
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

// URL combines the current URI with the given baseURL if not empty
func (j *UpJob) ConcatUrl(baseURL string) {
	j.URL = fmt.Sprintf("%s%s", strings.TrimSpace(baseURL), strings.TrimSpace(j.URLSuffix))
}
