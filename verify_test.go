package rashomon

import (
	"crypto/ed25519"
	"crypto/rand"
	"net/http"
	"testing"
)

var EMPTY_BYTE = []byte("")

func Test_VerifySignature(t *testing.T) {
	var tests = []struct {
		testDesc      string
		request       CanonicalRequest
		signature     *[]byte
		expectedValid bool
	}{
		{
			"Happy Path",
			CanonicalRequest{
				http.MethodPost,
				"https://localhost:9090/v1/thing2?param=1&sort=true",
				"8d98a306-3745-4d74-8bf9-e03a7019d49d",
				"2025-04-20T23:20:40+09:00",
				"dcdgqwfo",
				"text/plain",
				"line1\nline2",
			},
			&EMPTY_BYTE,
			false,
		},
		{
			"Bad Signature",
			CanonicalRequest{
				http.MethodPost,
				"https://localhost:9090/v1/thing2?param=1&sort=true",
				"8d98a306-3745-4d74-8bf9-e03a7019d49d",
				"2025-04-20T23:20:40+09:00",
				"dcdgqwfo",
				"text/plain",
				"line1\nline2",
			},
			&EMPTY_BYTE,
			false,
		},
	}

	public_key, private_key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Errorf("Error generating key pair")
	}

	for _, tt := range tests {
		t.Run(tt.testDesc, func(t *testing.T) {
			var signature []byte

			if tt.signature == nil {
				signature = Sign(&private_key, &tt.request)
			} else {
				signature = *tt.signature
			}

			valid := VerifySignature(&public_key, &tt.request, &signature)

			if valid != tt.expectedValid {
				t.Errorf("\n[Got validity]\n%v \n[Want validity]\n%v", valid, tt.expectedValid)
			}
		})
	}
}
