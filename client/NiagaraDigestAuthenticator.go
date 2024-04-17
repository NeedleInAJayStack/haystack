package client

import (
	"crypto/sha256"
	"io"
	"net/http"
	"strings"

	"github.com/NeedleInAJayStack/haystack/auth"
)

// NiagaraDigestAuthenticator is an Authenticator that conforms to the Niagara Digest authentication
type NiagaraDigestAuthenticator struct{}

func (authenticator NiagaraDigestAuthenticator) Authenticate(
	uri string,
	username string,
	password string,
	client ClientHTTP,
) (map[string]string, error) {
	// Hit prelogin to set cookies
	preloginUri := uri + "prelogin?j_username=" + username
	println(preloginUri)
	preloginReq, _ := http.NewRequest("POST", preloginUri, strings.NewReader("j_username="+username))
	setHeaders(preloginReq, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	preloginResp, preloginErr := client.do(preloginReq)
	if preloginErr != nil {
		return map[string]string{}, preloginErr
	}
	if preloginResp.StatusCode != http.StatusOK {
		return map[string]string{}, NewHTTPError(preloginResp.StatusCode, preloginResp.Status)
	}
	for _, cookie := range preloginResp.Cookies() {
		print(cookie.Name)
		print(":")
		print(cookie.Value)
	}

	securityCheckUri := uri + "j_security_check/"
	niagaraCookie := http.Cookie{Name: "niagara_userid", Value: username}
	hash := sha256.New
	var in []byte
	var scram = auth.NewScram(hash, username, password)
	var jsession string
	messageNumber := 1
	for !scram.Step(in) {
		out := scram.Out()

		var data []byte
		if messageNumber == 1 {
			msg := "action=sendClientFirstMessage&clientFirstMessage=" + string(out)
			req, _ := http.NewRequest("POST", securityCheckUri, strings.NewReader(msg))
			setHeaders(req, map[string]string{"Content-Type": "application/x-niagara-login-support"})
			req.AddCookie(&niagaraCookie)
			resp, err := client.do(req)
			if err != nil {
				return map[string]string{}, err
			}
			if resp.StatusCode != http.StatusOK { // We expect 200s
				return map[string]string{}, NewHTTPError(resp.StatusCode, resp.Status)
			}
			for _, setCookieHeader := range resp.Header["Set-Cookie"] {
				setCookie := strings.Split(setCookieHeader, ",")
				for _, key := range setCookie {
					for _, split := range strings.Split(key, ";") {
						if strings.Contains(split, "JSESSIONID=") {
							jsession = strings.Split(split, "=")[1]
						}
					}
				}
			}
			body, _ := io.ReadAll(resp.Body)
			data = body
		}
		if messageNumber == 2 {
			msg := "action=sendClientFinalMessage&clientFinalMessage=" + string(out)
			req, _ := http.NewRequest("POST", securityCheckUri, strings.NewReader(msg))
			setHeaders(req, map[string]string{"Content-Type": "application/x-niagara-login-support"})
			req.AddCookie(&niagaraCookie)
			jsessionCookie := http.Cookie{Name: "JSESSIONID", Value: jsession}
			req.AddCookie(&jsessionCookie)
			resp, err := client.do(req)
			if err != nil {
				return map[string]string{}, err
			}
			if resp.StatusCode != http.StatusOK { // We expect 200s
				return map[string]string{}, NewHTTPError(resp.StatusCode, resp.Status)
			}
			body, _ := io.ReadAll(resp.Body)
			data = body
		}

		in = data
		messageNumber = messageNumber + 1
	}
	if scram.Err() != nil {
		return map[string]string{}, scram.Err()
	}

	// Validate login
	req, _ := http.NewRequest("POST", securityCheckUri, nil)
	setHeaders(req, map[string]string{"Content-Type": "application/x-niagara-login-support"})
	resp, _ := client.do(req)
	if resp.StatusCode != http.StatusOK { // We expect 200s
		return map[string]string{}, NewHTTPError(resp.StatusCode, resp.Status)
	}

	// Since login is cookie-based, we don't add any new headers, allowing the underlying client to support
	return map[string]string{}, nil
}
