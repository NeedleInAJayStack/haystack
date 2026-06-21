package auth

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"hash"
	"strconv"
	"strings"
)

const scramDefaultIterations = 4096

var unescaper = strings.NewReplacer("=2C", ",", "=3D", "=")

// ScramServer implements the server side of a SCRAM-* conversation.
type ScramServer struct {
	newHash func() hash.Hash

	user string
	pass string
	step int
	out  bytes.Buffer
	err  error

	clientNonce []byte
	serverNonce []byte
	salt        []byte
	iterCount   int
	saltedPass  []byte
	authMsg     bytes.Buffer
}

// NewScramServer returns a new SCRAM-* server with the provided hash algorithm.
func NewScramServer(newHash func() hash.Hash, user, pass string) *ScramServer {
	s := &ScramServer{
		newHash: newHash,
		user:    user,
		pass:    pass,
	}
	s.out.Grow(256)
	s.authMsg.Grow(256)
	return s
}

// Out returns the data to be sent to the client in the current step.
func (s *ScramServer) Out() []byte {
	if s.out.Len() == 0 {
		return nil
	}
	return s.out.Bytes()
}

// Err returns the error that occurred, or nil if there were no errors.
func (s *ScramServer) Err() error {
	return s.err
}

// SetNonce sets the server nonce to the provided value.
func (s *ScramServer) SetNonce(nonce []byte) {
	s.serverNonce = nonce
}

// SetSalt sets the salt to the provided value.
func (s *ScramServer) SetSalt(salt []byte) {
	s.salt = salt
}

// SetIterations sets the SCRAM iteration count.
func (s *ScramServer) SetIterations(iter int) {
	s.iterCount = iter
}

// Step processes the incoming data from the client and makes the next round of
// data for the client available via ScramServer.Out. Step returns false if
// there are no errors and more data is still expected.
func (s *ScramServer) Step(in []byte) bool {
	s.out.Reset()
	if s.step > 1 || s.err != nil {
		return false
	}
	s.step++
	switch s.step {
	case 1:
		s.err = s.step1(in)
	case 2:
		s.err = s.step2(in)
	}
	return s.step > 1 || s.err != nil
}

func (s *ScramServer) step1(in []byte) error {
	if !bytes.HasPrefix(in, []byte("n,,")) {
		return fmt.Errorf("expected SCRAM client first message to start with %q: %q", "n,,", in)
	}
	clientFirstBare := in[len("n,,"):]
	fields := bytes.Split(clientFirstBare, []byte(","))
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 fields in first SCRAM client message, got %d: %q", len(fields), in)
	}
	if !bytes.HasPrefix(fields[0], []byte("n=")) {
		return fmt.Errorf("client sent an invalid SCRAM username: %q", fields[0])
	}
	if !bytes.HasPrefix(fields[1], []byte("r=")) || len(fields[1]) < 3 {
		return fmt.Errorf("client sent an invalid SCRAM nonce: %q", fields[1])
	}

	user := unescaper.Replace(string(fields[0][2:]))
	if user != s.user {
		return fmt.Errorf("unexpected SCRAM username: %q", user)
	}

	s.clientNonce = append(s.clientNonce[:0], fields[1][2:]...)
	if len(s.serverNonce) == 0 {
		nonce, err := scramGenerateNonce()
		if err != nil {
			return err
		}
		s.serverNonce = nonce
	}
	if len(s.salt) == 0 {
		salt, err := scramGenerateSalt()
		if err != nil {
			return err
		}
		s.salt = salt
	}
	if s.iterCount <= 0 {
		s.iterCount = scramDefaultIterations
	}

	s.authMsg.Write(clientFirstBare)
	serverFirst := s.serverFirstMessage()
	s.authMsg.WriteByte(',')
	s.authMsg.Write(serverFirst)
	s.out.Write(serverFirst)
	return nil
}

func (s *ScramServer) step2(in []byte) error {
	fields := bytes.Split(in, []byte(","))
	if len(fields) != 3 {
		return fmt.Errorf("expected 3 fields in final SCRAM client message, got %d: %q", len(fields), in)
	}
	if !bytes.Equal(fields[0], []byte("c=biws")) {
		return fmt.Errorf("client sent an invalid SCRAM channel binding: %q", fields[0])
	}
	if !bytes.HasPrefix(fields[1], []byte("r=")) || len(fields[1]) < 3 {
		return fmt.Errorf("client sent an invalid SCRAM combined nonce: %q", fields[1])
	}
	if !bytes.HasPrefix(fields[2], []byte("p=")) || len(fields[2]) < 3 {
		return fmt.Errorf("client sent an invalid SCRAM proof: %q", fields[2])
	}

	combinedNonce := fields[1][2:]
	expectedNonce := make([]byte, 0, len(s.clientNonce)+len(s.serverNonce))
	expectedNonce = append(expectedNonce, s.clientNonce...)
	expectedNonce = append(expectedNonce, s.serverNonce...)
	if !bytes.Equal(combinedNonce, expectedNonce) {
		return fmt.Errorf("client sent an unexpected SCRAM combined nonce: got %q, want %q", combinedNonce, expectedNonce)
	}

	s.saltPassword(s.salt, s.iterCount)
	s.authMsg.WriteString(",c=biws,r=")
	s.authMsg.Write(combinedNonce)
	if !bytes.Equal(fields[2][2:], s.clientProof()) {
		return fmt.Errorf("client sent an invalid SCRAM proof: %q", fields[2][2:])
	}

	s.out.WriteString("v=")
	s.out.Write(s.serverSignature())
	return nil
}

func (s *ScramServer) serverFirstMessage() []byte {
	serverFirst := new(bytes.Buffer)
	serverFirst.Grow(256)
	serverFirst.WriteString("r=")
	serverFirst.Write(s.clientNonce)
	serverFirst.Write(s.serverNonce)
	serverFirst.WriteString(",s=")
	encodedSalt := make([]byte, b64Std.EncodedLen(len(s.salt)))
	b64Std.Encode(encodedSalt, s.salt)
	serverFirst.Write(encodedSalt)
	serverFirst.WriteString(",i=")
	serverFirst.WriteString(strconv.Itoa(s.iterCount))
	return serverFirst.Bytes()
}

func (s *ScramServer) saltPassword(salt []byte, iterCount int) {
	mac := hmac.New(s.newHash, []byte(s.pass))
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
	s.saltedPass = hi
}

func (s *ScramServer) clientProof() []byte {
	mac := hmac.New(s.newHash, s.saltedPass)
	mac.Write([]byte("Client Key"))
	clientKey := mac.Sum(nil)
	hash := s.newHash()
	hash.Write(clientKey)
	storedKey := hash.Sum(nil)
	mac = hmac.New(s.newHash, storedKey)
	mac.Write(s.authMsg.Bytes())
	clientProof := mac.Sum(nil)
	for i, b := range clientKey {
		clientProof[i] ^= b
	}
	clientProof64 := make([]byte, b64Std.EncodedLen(len(clientProof)))
	b64Std.Encode(clientProof64, clientProof)
	return clientProof64
}

func (s *ScramServer) serverSignature() []byte {
	mac := hmac.New(s.newHash, s.saltedPass)
	mac.Write([]byte("Server Key"))
	serverKey := mac.Sum(nil)

	mac = hmac.New(s.newHash, serverKey)
	mac.Write(s.authMsg.Bytes())
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
