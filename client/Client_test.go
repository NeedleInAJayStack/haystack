package client

import (
	"fmt"
	"testing"

	"gitlab.com/NeedleInAJayStack/haystack"
)

func TestClient(t *testing.T) {

	// Testing with local SkySpark instance
	uri := "http://localhost:8080/api/demo"
	user := "test"
	pass := "test"

	haystackClient := NewClient(uri, user, pass)
	openErr := haystackClient.Open()
	if openErr != nil {
		t.Error(openErr)
	}

	about, err := haystackClient.Call("about", haystack.Grid{})
	if err != nil {
		fmt.Println(about.ToZinc())
		t.Error(err)
	}
}
