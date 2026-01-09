package crypto

import (
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/ed25519"
)

func VerifySignature(
	publicKeyBase64 string,
	messageBase64 string,
	signatureBase64 string,
) error {

	pubKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return errors.New("invalid public key encoding")
	}

	messageBytes, err := base64.StdEncoding.DecodeString(messageBase64)
	if err != nil {
		return errors.New("invalid message encoding")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return errors.New("invalid signature encoding")
	}

	if len(pubKeyBytes) != ed25519.PublicKeySize {
		return errors.New("invalid public key size")
	}

	if !ed25519.Verify(pubKeyBytes, messageBytes, sigBytes) {
		return errors.New("signature verification failed")
	}

	return nil
}
