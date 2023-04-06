package client

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestOps(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		print(".env file not found, falling back to local environment")
	}

	client := NewClient(
		os.Getenv("TEST_URL"),
		os.Getenv("TEST_USERNAME"),
		os.Getenv("TEST_PASSWORD"),
	)
	if client.Open() != nil {
		panic("Failed to open connection. Is a local Haxall server running?")
	}

	ops, err := client.Ops()
	if err != nil {
		panic(err)
	}

	print(ops.ToZinc())
}
