package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"

	"io/ioutil"

	"github.com/go-errors/errors"
)

type SuperAgent struct {
	Client *http.Client
	URL    string
}

func NewSuperAgent(rawurl string) *SuperAgent {
	return &SuperAgent{
		URL:    rawurl,
		Client: http.DefaultClient,
	}
}

func (s *SuperAgent) Delete() error {
	req, err := http.NewRequest("DELETE", s.URL, nil)
	if err != nil {
		return errors.New(err)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusNoContent, resp.StatusCode, body)
	}

	return nil
}

func (s *SuperAgent) Get(o interface{}) error {
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		return errors.New(err)
	} else if o == nil {
		return errors.New("Can not pass nil")
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", http.StatusOK, resp.StatusCode, body)
	} else if err := json.NewDecoder(resp.Body).Decode(o); err != nil {
		return errors.New(err)
	}

	return nil
}

func (s *SuperAgent) Create(o interface{}) error {
	return s.send("POST", o, o)
}

func (s *SuperAgent) POST(in, out interface{}) error {
	return s.send("POST", in, out)
}

func (s *SuperAgent) Update(o interface{}) error {
	return s.send("PUT", o, o)
}

func (s *SuperAgent) send(method string, in interface{}, out interface{}) error {
	if s.Client == nil {
		s.Client = http.DefaultClient
	}

	data, err := json.Marshal(in)
	if err != nil {
		return errors.New(err)
	}

	req, err := http.NewRequest(method, s.URL, bytes.NewReader(data))
	if err != nil {
		return errors.New(err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := s.Client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	expectedStatus := http.StatusOK
	if method == "POST" {
		expectedStatus = http.StatusCreated
	}
	if resp.StatusCode != expectedStatus {
		body, _ := ioutil.ReadAll(resp.Body)
		return errors.Errorf("Expected status code %d, got %d.\n%s\n", expectedStatus, resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return errors.New(err)
	}

	return nil
}
