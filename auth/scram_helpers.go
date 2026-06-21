package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"
)

// Need both because Haystack seems to switch between them for different parts.
var b64Std = base64.StdEncoding
var b64Uri = base64.RawURLEncoding // No padding

func scramSaltPassword(newHash func() hash.Hash, password string, salt []byte, iterCount int) []byte {
	mac := hmac.New(newHash, []byte(password))
	mac.Write(salt)
	mac.Write([]byte{0, 0, 0, 1})
	ui := mac.Sum(nil)
	hi := make([]byte, len(ui))
	copy(hi, ui)
	for i := 1; i < iterCount; i++ {
		mac.Reset()
		mac.Write(ui)
		mac.Sum(ui[:0])
		for j, b := range ui {
			hi[j] ^= b
		}
	}
	return hi
}

func scramClientProof(newHash func() hash.Hash, saltedPass []byte, authMsg []byte) []byte {
	mac := hmac.New(newHash, saltedPass)
	mac.Write([]byte("Client Key"))
	clientKey := mac.Sum(nil)
	hash := newHash()
	hash.Write(clientKey)
	storedKey := hash.Sum(nil)
	mac = hmac.New(newHash, storedKey)
	mac.Write(authMsg)
	clientProof := mac.Sum(nil)
	for i, b := range clientKey {
		clientProof[i] ^= b
	}
	clientProof64 := make([]byte, b64Std.EncodedLen(len(clientProof)))
	b64Std.Encode(clientProof64, clientProof)
	return clientProof64
}

func scramServerSignature(newHash func() hash.Hash, saltedPass []byte, authMsg []byte) []byte {
	mac := hmac.New(newHash, saltedPass)
	mac.Write([]byte("Server Key"))
	serverKey := mac.Sum(nil)

	mac = hmac.New(newHash, serverKey)
	mac.Write(authMsg)
	serverSignature := mac.Sum(nil)

	encoded := make([]byte, b64Std.EncodedLen(len(serverSignature)))
	b64Std.Encode(encoded, serverSignature)
	return encoded
}

func scramGenerateNonce() ([]byte, error) {
	const nonceLen = 16
	buf := make([]byte, nonceLen+b64Uri.EncodedLen(nonceLen))
	if _, err := rand.Read(buf[:nonceLen]); err != nil {
		return nil, fmt.Errorf("cannot read random SCRAM nonce from operating system: %v", err)
	}
	nonce := buf[nonceLen:]
	b64Uri.Encode(nonce, buf[:nonceLen])
	return nonce, nil
}

func scramGenerateSalt() ([]byte, error) {
	const saltLen = 16
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("cannot read random SCRAM salt from operating system: %v", err)
	}
	return salt, nil
}
