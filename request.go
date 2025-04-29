package rashomon

import (
	"strings"
)

type CanonicalRequest struct {
	Method      string
	URL         string
	KeyID       string
	Timestamp   string
	Nonce       string
	ContentType string
	Body        string
}

func (c *CanonicalRequest) String() string {
	return strings.Join([]string{
		c.Method,
		c.URL,
		"x-rmn-key-id:" + c.KeyID,
		"x-rmn-timestamp:" + c.Timestamp,
		"x-rmn-nonce:" + c.Nonce,
		"content-type:" + c.ContentType,
		"body:" + c.Body,
	}, "\n")
}
