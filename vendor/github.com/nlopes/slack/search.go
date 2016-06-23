package slack

import (
	"errors"
	"net/url"
	"strconv"
)

const (
	DEFAULT_SEARCH_SORT      = "score"
	DEFAULT_SEARCH_SORT_DIR  = "desc"
	DEFAULT_SEARCH_HIGHLIGHT = false
	DEFAULT_SEARCH_COUNT     = 100
	DEFAULT_SEARCH_PAGE      = 1
)

type SearchParameters struct {
	Sort          string
	SortDirection string
	Highlight     bool
	Count         int
	Page          int
}

type CtxChannel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CtxMessage struct {
	User      string `json:"user"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"`
	Type      string `json:"type"`
}

type SearchMessage struct {
	Type      string     `json:"type"`
	Channel   CtxChannel `json:"channel"`
	User      string     `json:"user"`
	Username  string     `json:"username"`
	Timestamp string     `json:"ts"`
	Text      string     `json:"text"`
	Permalink string     `json:"permalink"`
	Previous  CtxMessage `json:"previous"`
	Previous2 CtxMessage `json:"previous_2"`
	Next      CtxMessage `json:"next"`
	Next2     CtxMessage `json:"next_2"`
}

type SearchMessages struct {
	Matches    []SearchMessage `json:"matches"`
	Paging     `json:"paging"`
	Pagination `json:"pagination"`
	Total      int `json:"total"`
}

type SearchFiles struct {
	Matches    []File `json:"matches"`
	Paging     `json:"paging"`
	Pagination `json:"pagination"`
	Total      int `json:"total"`
}

type searchResponseFull struct {
	Query          string `json:"query"`
	SearchMessages `json:"messages"`
	SearchFiles    `json:"files"`
	SlackResponse
}

func NewSearchParameters() SearchParameters {
	return SearchParameters{
		Sort:          DEFAULT_SEARCH_SORT,
		SortDirection: DEFAULT_SEARCH_SORT_DIR,
		Highlight:     DEFAULT_SEARCH_HIGHLIGHT,
		Count:         DEFAULT_SEARCH_COUNT,
		Page:          DEFAULT_SEARCH_PAGE,
	}
}

func (api *Client) _search(path, query string, params SearchParameters, files, messages bool) (response *searchResponseFull, error error) {
	values := url.Values{
		"token": {api.config.token},
		"query": {query},
	}
	if params.Sort != DEFAULT_SEARCH_SORT {
		values.Add("sort", params.Sort)
	}
	if params.SortDirection != DEFAULT_SEARCH_SORT_DIR {
		values.Add("sort_dir", params.SortDirection)
	}
	if params.Highlight != DEFAULT_SEARCH_HIGHLIGHT {
		values.Add("highlight", strconv.Itoa(1))
	}
	if params.Count != DEFAULT_SEARCH_COUNT {
		values.Add("count", strconv.Itoa(params.Count))
	}
	if params.Page != DEFAULT_SEARCH_PAGE {
		values.Add("page", strconv.Itoa(params.Page))
	}
	response = &searchResponseFull{}
	err := post(path, values, response, api.debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil

}

func (api *Client) Search(query string, params SearchParameters) (*SearchMessages, *SearchFiles, error) {
	response, err := api._search("search.all", query, params, true, true)
	if err != nil {
		return nil, nil, err
	}
	return &response.SearchMessages, &response.SearchFiles, nil
}

func (api *Client) SearchFiles(query string, params SearchParameters) (*SearchFiles, error) {
	response, err := api._search("search.files", query, params, true, false)
	if err != nil {
		return nil, err
	}
	return &response.SearchFiles, nil
}

func (api *Client) SearchMessages(query string, params SearchParameters) (*SearchMessages, error) {
	response, err := api._search("search.messages", query, params, false, true)
	if err != nil {
		return nil, err
	}
	return &response.SearchMessages, nil
}
