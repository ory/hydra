package slack

import (
	"errors"
	"net/url"
)

type TeamResponse struct {
	Team TeamInfo `json:"team"`
	SlackResponse
}

type TeamInfo struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Domain      string                 `json:"domain"`
	EmailDomain string                 `json:"email_domain"`
	Icon        map[string]interface{} `json:"icon"`
}

func teamRequest(path string, values url.Values, debug bool) (*TeamResponse, error) {
	response := &TeamResponse{}
	err := post(path, values, response, debug)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, errors.New(response.Error)
	}

	return response, nil
}

// GetTeamInfo gets the Team Information of the user
func (api *Client) GetTeamInfo() (*TeamInfo, error) {
	values := url.Values{
		"token": {api.config.token},
	}

	response, err := teamRequest("team.info", values, api.debug)
	if err != nil {
		return nil, err
	}
	return &response.Team, nil
}
