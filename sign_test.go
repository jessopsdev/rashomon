package rashomon

import (
	"crypto/ed25519"
	"crypto/rand"
	"net/http"
	"testing"
)

func Test_Sign(t *testing.T) {
	var tests = []struct {
		testDesc string
		request  CanonicalRequest
	}{
		{
			"Basic GET",
			CanonicalRequest{
				http.MethodGet,
				"https://localhost:9090/v1/thing?param=1",
				"a3c3c023ghj23s",
				"2025-04-20T23:20:40+09:00",
				"whatever",
				"",
				"",
			},
		},
		{
			"Basic POST",
			CanonicalRequest{
				http.MethodPost,
				"https://localhost:9090/v1/thing2?param=1&sort=true",
				"8d98a306-3745-4d74-8bf9-e03a7019d49d",
				"2025-04-20T23:20:40+09:00",
				"dcdgqwfo",
				"text/plain",
				"cool stuff in the body",
			},
		},
		{
			"JSON POST",
			CanonicalRequest{
				http.MethodPost,
				"https://localhost:9090/v1/thing2?param=1&sort=true",
				"8d98a306-3745-4d74-8bf9-e03a7019d49d",
				"2025-04-20T23:20:40+09:00",
				"dcdgqwfo",
				"application/json",
				`{"some":"data","more":1234}`,
			},
		},
		{
			"Newline in Body",
			CanonicalRequest{
				http.MethodPost,
				"https://localhost:9090/v1/thing2?param=1&sort=true",
				"8d98a306-3745-4d74-8bf9-e03a7019d49d",
				"2025-04-20T23:20:40+09:00",
				"dcdgqwfo",
				"text/plain",
				"line1\nline2",
			},
		},
	}

	public_key, private_key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Errorf("Error generating key pair")
	}

	for _, tt := range tests {
		t.Run(tt.testDesc, func(t *testing.T) {
			signature := Sign(&private_key, &tt.request)

			valid := VerifySignature(&public_key, &tt.request, &signature)

			if !valid {
				t.Errorf("\n[Expected signature to verify but failed]")
			}
		})
	}
}
