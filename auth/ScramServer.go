package auth

import (
	"bytes"
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

	s.saltedPass = scramSaltPassword(s.newHash, s.pass, s.salt, s.iterCount)
	s.authMsg.WriteString(",c=biws,r=")
	s.authMsg.Write(combinedNonce)
	if !bytes.Equal(fields[2][2:], scramClientProof(s.newHash, s.saltedPass, s.authMsg.Bytes())) {
		return fmt.Errorf("client sent an invalid SCRAM proof: %q", fields[2][2:])
	}

	s.out.WriteString("v=")
	s.out.Write(scramServerSignature(s.newHash, s.saltedPass, s.authMsg.Bytes()))
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
