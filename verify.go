package rashomon

import (
	"crypto/ed25519"
)

func VerifySignature(publicKey *ed25519.PublicKey, req *CanonicalRequest, signature *[]byte) bool {
	return ed25519.Verify(*publicKey, []byte(req.String()), *signature)
}
