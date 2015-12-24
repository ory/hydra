// Copyright 2015 Niels Freier
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package elastigo

import (
	"encoding/json"
	"fmt"
	"time"
)

type GetSnapshotsResponse struct {
	Snapshots []struct {
		Snapshot  string    `json:"snapshot"`
		Indices   []string  `json:"indices"`
		State     string    `json:"state"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	} `json:"snapshots"`
}

// CreateSnapshotRepository creates a new snapshot repository on the cluster
// http://www.elastic.co/guide/en/elasticsearch/reference/1.3/modules-snapshots.html
func (c *Conn) CreateSnapshotRepository(name string, args map[string]interface{}, settings interface{}) (BaseResponse, error) {
	var url string
	var retval BaseResponse
	url = fmt.Sprintf("/_snapshot/%s", name)
	body, err := c.DoCommand("POST", url, args, settings)
	if err != nil {
		return retval, err
	}
	if err == nil {
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}

	return retval, nil

}

// TakeSnapshot takes a snapshot of the current state of the cluster with a specific name and for a existing repositoriy
// http://www.elastic.co/guide/en/elasticsearch/reference/1.3/modules-snapshots.html
func (c *Conn) TakeSnapshot(repository, name string, args map[string]interface{}, query interface{}) (BaseResponse, error) {
	var url string
	var retval BaseResponse
	url = fmt.Sprintf("/_snapshot/%s/%s", repository, name)
	body, err := c.DoCommand("PUT", url, args, query)
	if err != nil {
		return retval, err
	}
	if err == nil {
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}

	return retval, nil
}

// RestoreSnapshot restores a snapshot of the current state of the cluster with a specific name and for a existing repositoriy
// http://www.elastic.co/guide/en/elasticsearch/reference/1.3/modules-snapshots.html
func (c *Conn) RestoreSnapshot(repository, name string, args map[string]interface{}, query interface{}) (BaseResponse, error) {
	var url string
	var retval BaseResponse
	url = fmt.Sprintf("/_snapshot/%s/%s/_restore", repository, name)
	body, err := c.DoCommand("POST", url, args, query)
	if err != nil {
		return retval, err
	}
	if err == nil {
		fmt.Println(string(body))
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}

	return retval, nil
}

// GetSnapshots returns all snapshot of the specified name for a specific repository
// http://www.elastic.co/guide/en/elasticsearch/reference/1.3/modules-snapshots.html
func (c *Conn) GetSnapshotByName(repository, name string, args map[string]interface{}) (GetSnapshotsResponse, error) {
	return c.getSnapshots(repository, name, args)
}

// GetSnapshots returns all snapshot for a specific repository
// http://www.elastic.co/guide/en/elasticsearch/reference/1.3/modules-snapshots.html
func (c *Conn) GetSnapshots(repository string, args map[string]interface{}) (GetSnapshotsResponse, error) {
	return c.getSnapshots(repository, "_all", args)
}

func (c *Conn) getSnapshots(repository, name string, args map[string]interface{}) (GetSnapshotsResponse, error) {
	var url string
	var retval GetSnapshotsResponse
	url = fmt.Sprintf("/_snapshot/%s/%s", repository, name)
	body, err := c.DoCommand("GET", url, args, nil)
	if err != nil {
		return retval, err
	}
	if err == nil {
		jsonErr := json.Unmarshal(body, &retval)
		if jsonErr != nil {
			return retval, jsonErr
		}
	}

	return retval, nil
}
