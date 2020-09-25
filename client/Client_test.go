package client

import "testing"

func TestClient(t *testing.T) {

	uri := "http://localhost:8080/api/demo"
	user := "test"
	pass := "test"

	haystackClient := NewClient(uri, user, pass)
	err := haystackClient.Open()

	if err != nil {
		t.Error(err)
	}
}
