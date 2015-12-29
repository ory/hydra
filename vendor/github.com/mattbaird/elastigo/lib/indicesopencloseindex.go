// Copyright 2013 Matthew Baird
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
)

func (c *Conn) OpenIndices() (BaseResponse, error) {
	return c.openCloseOperation("_all", "_open")
}

func (c *Conn) CloseIndices() (BaseResponse, error) {
	return c.openCloseOperation("_all", "_close")
}

func (c *Conn) OpenIndex(index string) (BaseResponse, error) {
	return c.openCloseOperation(index, "_open")
}

func (c *Conn) CloseIndex(index string) (BaseResponse, error) {
	return c.openCloseOperation(index, "_close")
}

func (c *Conn) openCloseOperation(index, mode string) (BaseResponse, error) {
	var url string
	var retval BaseResponse

	if len(index) > 0 {
		url = fmt.Sprintf("/%s/%s", index, mode)
	} else {
		url = fmt.Sprintf("/%s", mode)
	}

	body, errDo := c.DoCommand("POST", url, nil, nil)
	if errDo != nil {
		return retval, errDo
	}
	jsonErr := json.Unmarshal(body, &retval)
	if jsonErr != nil {
		return retval, jsonErr
	}
	return retval, errDo
}
