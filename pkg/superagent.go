package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"

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

func (s *SuperAgent) DELETE() error {
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
		return errors.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}

	return nil
}

func (s *SuperAgent) GET(o interface{}) error {
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
		return errors.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	} else if err := json.NewDecoder(resp.Body).Decode(o); err != nil {
		return errors.New(err)
	}

	return nil
}

func (s *SuperAgent) POST(o interface{}) error {
	return s.send("POST", o)
}

func (s *SuperAgent) PUT(o interface{}) error {
	return s.send("PUT", o)
}

func (s *SuperAgent) send(method string, o interface{}) error {
	if s.Client == nil {
		s.Client = http.DefaultClient
	}

	data, err := json.Marshal(o)
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
		return errors.Errorf("Expected status code %d, got %d", expectedStatus, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(o); err != nil {
		return errors.New(err)
	}

	return nil
}
