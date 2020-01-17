package job

import "fmt"

import "strings"

// Body is taken from plain body or transformed form base64 body if exists
func (check UpCheck) Body() []byte {
	return nil
}

// URL combines the current URI with the given baseURL if not empty
func (check UpCheck) URL(baseURL string) string {
	return fmt.Sprintf(
		"%s%s",
		strings.TrimSpace(baseURL),
		strings.TrimSpace(check.URI),
	)
}
