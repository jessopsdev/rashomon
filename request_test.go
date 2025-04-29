package rashomon

import (
	"net/http"
	"testing"
)

func Test_CanonicalRequest_String(t *testing.T) {
	var tests = []struct {
		testDesc string
		request  CanonicalRequest
		want     string
		hasError bool
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
			"GET\nhttps://localhost:9090/v1/thing?param=1\nx-rmn-key-id:a3c3c023ghj23s\nx-rmn-timestamp:2025-04-20T23:20:40+09:00\nx-rmn-nonce:whatever\ncontent-type:\nbody:",
			false,
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
			"POST\nhttps://localhost:9090/v1/thing2?param=1&sort=true\nx-rmn-key-id:8d98a306-3745-4d74-8bf9-e03a7019d49d\nx-rmn-timestamp:2025-04-20T23:20:40+09:00\nx-rmn-nonce:dcdgqwfo\ncontent-type:text/plain\nbody:cool stuff in the body",
			false,
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
			"POST\nhttps://localhost:9090/v1/thing2?param=1&sort=true\nx-rmn-key-id:8d98a306-3745-4d74-8bf9-e03a7019d49d\nx-rmn-timestamp:2025-04-20T23:20:40+09:00\nx-rmn-nonce:dcdgqwfo\ncontent-type:application/json\nbody:{\"some\":\"data\",\"more\":1234}",
			false,
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
			"POST\nhttps://localhost:9090/v1/thing2?param=1&sort=true\nx-rmn-key-id:8d98a306-3745-4d74-8bf9-e03a7019d49d\nx-rmn-timestamp:2025-04-20T23:20:40+09:00\nx-rmn-nonce:dcdgqwfo\ncontent-type:text/plain\nbody:line1\nline2",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testDesc, func(t *testing.T) {
			// Generate string
			requestString := tt.request.String()
			if requestString != tt.want {
				t.Errorf("\n[Got result]\n%s \n[Want result]\n%s", requestString, tt.want)
			}
		})
	}
}

// func Test_ServerSide_makeRequestString(t *testing.T) {
// 	defer req.Body.Close()
// 	bodyBytes, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read request body: %w", err)
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.testDesc, func(t *testing.T) {
// 			req := httptest.NewRequest(tt.method, tt.url, tt.body)
// 			req.Header.Add("x-rmn-key-id", tt.keyId)
// 			req.Header.Add("x-rmn-timestamp", tt.timestamp)
// 			req.Header.Add("x-rmn-nonce", tt.nonce)
// 			req.Header.Add("Content-Type", tt.contentType)

// 			rr := httptest.NewRecorder()

// 			func(w http.ResponseWriter, r *http.Request) {
// 				requestString, err := makeRequestString(r)
// 				if err != nil {
// 					if !tt.hasError {
// 						t.Error("got unexpected error")
// 					}
// 				}

// 				if requestString != tt.want {
// 					t.Errorf("\n[Got result]\n%s \n[Want result]\n%s", requestString, tt.want)
// 				}

// 				w.WriteHeader(http.StatusOK)
// 			}(rr, req)
// 		})
// 	}
// }
