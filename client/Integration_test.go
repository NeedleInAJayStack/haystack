package client

import (
	"crypto/tls"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestOps(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		print(".env file not found, falling back to local environment")
	}

	// Disable TLS verification. ONLY FOR TESTING
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	client := NewClient(
		os.Getenv("TEST_URL"),
		os.Getenv("TEST_USERNAME"),
		os.Getenv("TEST_PASSWORD"),
		// NiagaraDigestAuthenticator{},
	)
	openErr := client.Open()
	if openErr != nil {
		panic(openErr.Error())
	}

	ops, err := client.Ops()
	if err != nil {
		panic(err)
	}

	print(ops.ToZinc())
}
