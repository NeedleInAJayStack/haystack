// Copyright (c) 2014 - Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package auth

import (
	"bytes"
	"fmt"
	"hash"
	"strconv"
	"strings"
)

// Scram implements a SCRAM-* client (SCRAM-SHA-1, SCRAM-SHA-256, etc).
//
// A Scram may be used within a SASL conversation with logic resembling:
//
//	var in []byte
//	var client = scram.NewScram(sha1.New, user, pass)
//	for client.Step(in) {
//	        out := client.Out()
//	        // send out to server
//	        in := serverOut
//	}
//	if client.Err() != nil {
//	        // auth failed
//	}
type Scram struct {
	newHash func() hash.Hash

	user string
	pass string
	step int
	out  bytes.Buffer
	err  error

	clientNonce []byte
	serverNonce []byte
	saltedPass  []byte
	authMsg     bytes.Buffer
}

// NewScram returns a new SCRAM-* client with the provided hash algorithm.
//
// For SCRAM-SHA-256, for example, use:
//
//	client := scram.NewScram(sha256.New, user, pass)
func NewScram(newHash func() hash.Hash, user, pass string) *Scram {
	c := &Scram{
		newHash: newHash,
		user:    user,
		pass:    pass,
	}
	c.out.Grow(256)
	c.authMsg.Grow(256)
	return c
}

// Out returns the data to be sent to the server in the current step.
func (c *Scram) Out() []byte {
	if c.out.Len() == 0 {
		return nil
	}
	return c.out.Bytes()
}

// Err returns the error that occurred, or nil if there were no errors.
func (c *Scram) Err() error {
	return c.err
}

// SetNonce sets the client nonce to the provided value.
// If not set, the nonce is generated automatically out of crypto/rand on the first step.
func (c *Scram) SetNonce(nonce []byte) {
	c.clientNonce = nonce
}

var escaper = strings.NewReplacer("=", "=3D", ",", "=2C")

// Step processes the incoming data from the server and makes the
// next round of data for the server available via Scram.Out.
// Step returns false if there are no errors and more data is
// still expected.
func (c *Scram) Step(in []byte) bool {
	c.out.Reset()
	if c.step > 2 || c.err != nil {
		return false
	}
	c.step++
	switch c.step {
	case 1:
		c.err = c.step1(in)
	case 2:
		c.err = c.step2(in)
	case 3:
		c.err = c.step3(in)
	}
	return c.step > 2 || c.err != nil
}

func (c *Scram) step1(in []byte) error {
	if len(c.clientNonce) == 0 {
		nonce, err := scramGenerateNonce()
		if err != nil {
			return err
		}
		c.clientNonce = nonce
	}
	c.authMsg.WriteString("n=")
	escaper.WriteString(&c.authMsg, c.user)
	c.authMsg.WriteString(",r=")
	c.authMsg.Write(c.clientNonce)

	c.out.WriteString("n,,")
	c.out.Write(c.authMsg.Bytes())
	return nil
}

func (c *Scram) step2(in []byte) error {
	c.authMsg.WriteByte(',')
	c.authMsg.Write(in)

	fields := bytes.Split(in, []byte(","))
	if len(fields) != 3 {
		return fmt.Errorf("expected 3 fields in first SCRAM-SHA-256 server message, got %d: %q", len(fields), in)
	}
	if !bytes.HasPrefix(fields[0], []byte("r=")) || len(fields[0]) < 2 {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-256 nonce: %q", fields[0])
	}
	if !bytes.HasPrefix(fields[1], []byte("s=")) || len(fields[1]) < 6 {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-256 salt: %q", fields[1])
	}
	if !bytes.HasPrefix(fields[2], []byte("i=")) || len(fields[2]) < 6 {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-256 iteration count: %q", fields[2])
	}

	c.serverNonce = fields[0][2:]
	if !bytes.HasPrefix(c.serverNonce, c.clientNonce) {
		return fmt.Errorf("server SCRAM-SHA-256 nonce is not prefixed by client nonce: got %q, want %q+\"...\"", c.serverNonce, c.clientNonce)
	}

	salt := make([]byte, b64Std.DecodedLen(len(fields[1][2:])))
	n, err := b64Std.Decode(salt, fields[1][2:])
	if err != nil {
		return fmt.Errorf("cannot decode SCRAM-SHA-256 salt sent by server: %q", fields[1])
	}
	salt = salt[:n]
	iterCount, err := strconv.Atoi(string(fields[2][2:]))
	if err != nil {
		return fmt.Errorf("server sent an invalid SCRAM-SHA-256 iteration count: %q", fields[2])
	}
	c.saltedPass = scramSaltPassword(c.newHash, c.pass, salt, iterCount)

	c.authMsg.WriteString(",c=biws,r=")
	c.authMsg.Write(c.serverNonce)

	c.out.WriteString("c=biws,r=")
	c.out.Write(c.serverNonce)
	c.out.WriteString(",p=")
	c.out.Write(scramClientProof(c.newHash, c.saltedPass, c.authMsg.Bytes()))
	return nil
}

func (c *Scram) step3(in []byte) error {
	var isv, ise bool
	var fields = bytes.Split(in, []byte(","))
	if len(fields) == 1 {
		isv = bytes.HasPrefix(fields[0], []byte("v="))
		ise = bytes.HasPrefix(fields[0], []byte("e="))
	}
	if ise {
		return fmt.Errorf("SCRAM-SHA-256 authentication error: %s", fields[0][2:])
	} else if !isv {
		return fmt.Errorf("unsupported SCRAM-SHA-256 final message from server: %q", in)
	}
	if !bytes.Equal(scramServerSignature(c.newHash, c.saltedPass, c.authMsg.Bytes()), fields[0][2:]) {
		return fmt.Errorf("cannot authenticate SCRAM-SHA-256 server signature: %q", fields[0][2:])
	}
	return nil
}
