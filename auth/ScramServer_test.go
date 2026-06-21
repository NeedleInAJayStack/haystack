package auth

import (
	"bytes"
	"crypto/sha256"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScramClientServerHandshake(t *testing.T) {
	client := NewScram(sha256.New, "test", "test")
	client.SetNonce([]byte("clientnonce"))

	server := NewScramServer(sha256.New, "test", "test")
	server.SetNonce([]byte("servernonce"))
	server.SetSalt([]byte("salt-for-tests"))
	server.SetIterations(4096)

	var clientIn []byte
	for !client.Step(clientIn) {
		serverDone := server.Step(client.Out())
		assert.Nil(t, server.Err())
		clientIn = server.Out()
		if serverDone {
			break
		}
	}

	assert.Nil(t, client.Err())
	assert.Nil(t, server.Err())
	assert.Equal(t, "v=2BuVkVpguahLEJ+D5hifYQOF+3GxZnOoqQejYRDA8IU=", string(server.Out()))
	assert.True(t, client.Step(server.Out()))
	assert.Nil(t, client.Err())
}

func TestScramServerInvalidUsername(t *testing.T) {
	server := NewScramServer(sha256.New, "test", "test")
	server.SetNonce([]byte("servernonce"))
	server.SetSalt([]byte("salt-for-tests"))
	server.SetIterations(4096)

	done := server.Step([]byte("n,,n=wrong,r=clientnonce"))
	assert.True(t, done)
	assert.EqualError(t, server.Err(), `unexpected SCRAM username: "wrong"`)
}

func TestScramServerInvalidProof(t *testing.T) {
	server := NewScramServer(sha256.New, "test", "test")
	server.SetNonce([]byte("servernonce"))
	server.SetSalt([]byte("salt-for-tests"))
	server.SetIterations(4096)

	assert.False(t, server.Step([]byte("n,,n=test,r=clientnonce")))
	assert.Nil(t, server.Err())

	clientFinal := []byte("c=biws,r=clientnonceservernonce,p=invalid")
	done := server.Step(clientFinal)
	assert.True(t, done)
	assert.Error(t, server.Err())
	assert.True(t, strings.Contains(server.Err().Error(), "invalid SCRAM proof"))
}

func TestScramServerInvalidCombinedNonce(t *testing.T) {
	server := NewScramServer(sha256.New, "test", "test")
	server.SetNonce([]byte("servernonce"))
	server.SetSalt([]byte("salt-for-tests"))
	server.SetIterations(4096)

	assert.False(t, server.Step([]byte("n,,n=test,r=clientnonce")))
	assert.Nil(t, server.Err())

	client := NewScram(sha256.New, "test", "test")
	client.SetNonce([]byte("clientnonce"))
	assert.False(t, client.Step(nil))
	assert.False(t, client.Step(server.Out()))

	clientFinal := bytes.Replace(client.Out(), []byte("clientnonceservernonce"), []byte("clientnoncewrongnonce"), 1)
	done := server.Step(clientFinal)
	assert.True(t, done)
	assert.Error(t, server.Err())
	assert.True(t, strings.Contains(server.Err().Error(), "unexpected SCRAM combined nonce"))
}
