package rashomon

import (
	"crypto/ed25519"
)

func Sign(privateKey *ed25519.PrivateKey, req *CanonicalRequest) []byte {
	return ed25519.Sign(*privateKey, []byte(req.String()))
}
