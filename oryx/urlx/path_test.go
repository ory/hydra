// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package urlx

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURLFilePath(t *testing.T) {
	type testData struct {
		urlStr          string
		expectedUnix    string
		expectedWindows string
		shouldSucceed   bool
	}
	var testURLs = []testData{
		{"File:///home/test/file1.txt", "/home/test/file1.txt", "\\home\\test\\file1.txt", true},
		{"fIle:/home/test/file2.txt", "/home/test/file2.txt", "\\home\\test\\file2.txt", true},
		{"fiLe:///../test/update/file3.txt", "/../test/update/file3.txt", "\\..\\test\\update\\file3.txt", true},
		{"filE://../test/update/file4.txt", "../test/update/file4.txt", "..\\test\\update\\file4.txt", true},
		{"file://C:/users/test/file5.txt", "/C:/users/test/file5.txt", "C:\\users\\test\\file5.txt", true},
		{"file:///C:/users/test/file5b.txt", "/C:/users/test/file5b.txt", "C:\\users\\test\\file5b.txt", true},
		{"file://anotherhost/share/users/test/file6.txt", "/share/users/test/file6.txt", "\\\\anotherhost\\share\\users\\test\\file6.txt", false}, // this is not supported
		{"file://file7.txt", "file7.txt", "file7.txt", true},
		{"file://path/file8.txt", "path/file8.txt", "path\\file8.txt", true},
		{"file://C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343/access-rules.json", "/C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343/access-rules.json", "C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343\\access-rules.json", true},
		{"file:///C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343/access-rules.json", "/C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343/access-rules.json", "C:\\Users\\RUNNER~1\\AppData\\Local\\Temp\\9ccf9f68-121c-451a-8a73-2aa360925b5a386398343\\access-rules.json", true},
		{"file8.txt", "file8.txt", "file8.txt", true},
		{"../file9.txt", "../file9.txt", "..\\file9.txt", true},
		{"./file9b.txt", "./file9b.txt", ".\\file9b.txt", true},
		{"file://./file9c.txt", "./file9c.txt", ".\\file9c.txt", true},
		{"file://./folder/.././file9d.txt", "./folder/.././file9d.txt", ".\\folder\\..\\.\\file9d.txt", true},
		{"..\\file10.txt", "..\\file10.txt", "..\\file10.txt", true},
		{"C:\\file11.txt", "/C:\\file11.txt", "C:\\file11.txt", true},
		{"\\\\hostname\\share\\file12.txt", "/share/file12.txt", "\\\\hostname\\share\\file12.txt", true},
		{"file:///home/test/file 13.txt", "/home/test/file 13.txt", "\\home\\test\\file 13.txt", true},
		{"file:///home/test/file%2014.txt", "/home/test/file 14.txt", "\\home\\test\\file 14.txt", true},
		{"http://server:80/test/file%2015.txt", "/test/file 15.txt", "/test/file 15.txt", true},
		{"file:///dir/file\\ with backslash", "/dir/file\\ with backslash", "\\dir\\file\\ with backslash", true},
		{"file://dir/file\\ with backslash", "dir/file\\ with backslash", "dir\\file\\ with backslash", true},
		{"file:///dir/file with windows path forbidden chars \\<>:\"|%3F*", "/dir/file with windows path forbidden chars \\<>:\"|?*", "\\dir\\file with windows path forbidden chars \\<>:\"|?*", true},
		{"file://dir/file with windows path forbidden chars \\<>:\"|%3F*", "dir/file with windows path forbidden chars \\<>:\"|?*", "dir\\file with windows path forbidden chars \\<>:\"|?*", true},
		{"file:///path/file?query=1", "/path/file", "\\path\\file", true},
		{"http://host:80/path/file?query=1", "/path/file", "/path/file", true},
		{"file://////C:/file.txt", "////C:/file.txt", "C:\\file.txt", true},
		{"file://////C:\\file.txt", "////C:\\file.txt", "C:\\file.txt", true},
	}
	for _, td := range testURLs {
		u, err := Parse(td.urlStr)
		assert.NoError(t, err)
		if err != nil {
			continue
		}
		p := GetURLFilePath(u)
		if runtime.GOOS == "windows" {
			if td.shouldSucceed {
				assert.Equal(t, td.expectedWindows, p)
			} else {
				assert.NotEqual(t, td.expectedWindows, p)
			}
		} else {
			if td.shouldSucceed {
				assert.Equal(t, td.expectedUnix, p)
			} else {
				assert.NotEqual(t, td.expectedUnix, p)
			}
		}
	}
	assert.Empty(t, GetURLFilePath(nil))
}
