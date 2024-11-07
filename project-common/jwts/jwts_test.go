package jwts

import "testing"

func TestJwt(t *testing.T) {

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzEwMjUxMTUsInRva2VuIjoiMTAxMSJ9.407JF50fKXipGRqYJ0Wk2YTGqET3RKUp48qgfYtaztk"
	secret := "msys-project"

	ParseJwt(tokenString, secret)
}
