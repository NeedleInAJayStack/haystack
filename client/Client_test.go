package client

import "testing"

func TestClient(t *testing.T) {

	// Testing with local SkySpark instance
	uri := "http://localhost:8080/api/demo"
	user := "test"
	pass := "test"

	haystackClient := NewClient(uri, user, pass)
	err := haystackClient.Open()

	if err != nil {
		t.Error(err)
	}

	// {userAuth:{c:10000 hash:"SHA-256" salt:"p-_WN0lGn71_uyeAU4MzDRFXgw0k_yTumRIyielX9_0" scheme:"scram"}}

	// Messages from Haystack example:
	// C: n,,n=user,r=rOprNGfwEbeRWgbNEkqO
	// S: r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF,s=W22ZaJ0SNY7soEsUEjb6gQ==,i=4096
	// C: c=biws,r=rOprNGfwEbeRWgbNEkqO%hvYDpWUa2RaTCAfuxFIlj)hNlF,p=dHzbZapWIk4jUhN+Ute9ytag9zjfMHgsqmmiz7AndVQ=
}
