package airflow

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type AuthProvider interface {
	ApplyAuth(req *http.Request) error
}

type NoAuth struct{}

func (a *NoAuth) ApplyAuth(req *http.Request) error {
	return nil
}

type BasicAuth struct {
	Username string
	Password string
}

func (a *BasicAuth) ApplyAuth(req *http.Request) error {
	credentials := base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password))
	req.Header.Set("Authorization", "Basic "+credentials)
	return nil
}

type StaticToken struct {
	Token string
}

func (a *StaticToken) ApplyAuth(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return nil
}

type MWAAuth struct {
	Profile string
	Region  string
	Token   string
}

func (a *MWAAuth) ApplyAuth(req *http.Request) error {
	if a.Token == "" {
		return fmt.Errorf("MWAA token not initialized")
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)
	return nil
}
