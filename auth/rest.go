package auth

import (
	"encoding/json"
	"net/url"
	"github.com/pkg/errors"
	"bytes"
)

func Auth(a *AuthRequest, endpoint *url.URL) (result *AuthResponse, error) {
// func request(method string, path string, param interface{}, header map[string]string) (response []byte, errors definitions.Errors) {
	body, err := json.Marshal(a)
	if err != nil {
		nil, errors.New("Error while marshal JSON")
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.New("Error while creating request") 
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Error while requesting to the auth server")
	}
	if resp.StatusCode > 300 {
		return nil, errors.New("HttpError:" + resp.StatusCode)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	err := json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		nil, errors.New("Error while marshal JSON")
	}
	return
}
