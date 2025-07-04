/*-
 * Copyright 2019 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package josex

import "io"

// Base64Reader wraps an input stream consisting of either standard or url-safe
// base64 data, and maps it to a raw (unpadded) standard encoding. This can be used
// to read any base64-encoded data as input, whether padded, unpadded, standard or
// url-safe.
type Base64Reader struct {
	In io.Reader
}

func (r Base64Reader) Read(p []byte) (n int, err error) {
	n, err = r.In.Read(p)
	if err != nil {
		return
	}

	for i := range n {
		switch p[i] {
		// Map - to +
		case 0x2D:
			p[i] = 0x2B
		// Map _ to /
		case 0x5F:
			p[i] = 0x2F
		// Strip =
		case 0x3D:
			n = i
		default:
		}
	}

	if n == 0 {
		err = io.EOF
	}

	return
}
