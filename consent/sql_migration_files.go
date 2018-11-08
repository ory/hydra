// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// migrations/sql/shared/.gitkeep
// migrations/sql/shared/1.sql
// migrations/sql/shared/2.sql
// migrations/sql/shared/3.sql
// migrations/sql/mysql/.gitkeep
// migrations/sql/mysql/4.sql
// migrations/sql/mysql/5.sql
// migrations/sql/mysql/6.sql
// migrations/sql/postgres/.gitkeep
// migrations/sql/postgres/4.sql
// migrations/sql/postgres/5.sql
// migrations/sql/postgres/6.sql
// migrations/sql/tests/.gitkeep
// migrations/sql/tests/1_test.sql
// migrations/sql/tests/2_test.sql
// migrations/sql/tests/3_test.sql
// migrations/sql/tests/4_test.sql
// migrations/sql/tests/5_test.sql
// migrations/sql/tests/6_test.sql
package consent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _migrationsSqlSharedGitkeep = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func migrationsSqlSharedGitkeepBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlSharedGitkeep,
		"migrations/sql/shared/.gitkeep",
	)
}

func migrationsSqlSharedGitkeep() (*asset, error) {
	bytes, err := migrationsSqlSharedGitkeepBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/shared/.gitkeep", size: 0, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared1Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x95\x41\x8f\xda\x30\x10\x85\xcf\xce\xaf\x98\x23\xa8\xac\x54\xad\xba\xa7\x3d\xd1\x42\xa5\xaa\x14\x56\x08\x54\xed\xc9\x32\xf6\x40\xdc\x4d\x6c\x6a\x3b\x4b\xfb\xef\xab\x84\x04\x88\x93\x98\x2c\x5d\xae\xbc\x19\xcf\xbc\xef\xd9\xb9\xbb\x83\x0f\xa9\xdc\x19\xe6\x10\xd6\xfb\xe8\xcb\x72\x3a\x5e\x4d\x61\x35\xfe\x3c\x9b\x42\xfc\x57\x18\x46\x35\xcb\x5c\x7c\x4f\xb9\x56\x16\x95\xa3\x06\x7f\x67\x68\x1d\x0c\x22\xc2\x63\x96\x24\xa8\x76\x08\x40\x08\x01\x78\x65\x86\xc7\xcc\x0c\x3e\x7d\x1c\xc2\x7c\xb1\x82\xf9\x7a\x36\x83\xa7\xe5\xb7\x1f\xe3\xe5\x33\x7c\x9f\x3e\x8f\x22\xf2\x8a\x46\x6e\x25\x1a\x38\xfd\xda\x8a\x46\x11\xe1\x89\xcc\x4f\x93\xa2\xe8\x7c\x96\xdd\x3f\x3c\xd4\x74\x36\xdb\xfc\x42\xee\xc8\x15\x59\x39\x35\xcd\x4c\x52\x28\x1d\xfe\x71\xb5\x36\x2f\x72\x5f\xf5\x00\xd8\x68\x9d\xb4\x54\xa3\xa0\x96\xeb\x3d\x12\xe2\x97\x73\x6b\xb6\xe7\xf2\x8e\x95\x72\x1f\x51\x39\xc9\x59\xde\x89\x39\xe2\x64\x8a\xd6\xb1\x74\xdf\x3c\x87\xb9\xdc\xd3\x0b\x41\x65\xe7\x64\xfa\x75\xbc\x9e\xad\x40\xe9\xc3\x60\x38\x8a\x88\x96\x82\xe7\x6c\xf2\x89\x1a\x8b\x45\xc3\xc7\x00\xd1\x8b\x79\xa4\x56\xff\x09\xf6\xaa\x45\x27\xf2\x1e\x29\x9f\x7b\x1f\x27\x6f\x81\x0e\xd0\x18\x29\x0c\xbd\x6f\x02\x6f\x80\x76\x3d\x09\xef\x87\xd5\xa2\xb5\x52\xab\x1c\xab\x14\xc7\x4d\x03\x04\x3c\xaa\xfe\xa0\x00\x81\xed\xe6\x8b\x9f\xc5\x76\x25\x1d\xff\x98\x9a\x6f\xe1\x15\xbc\xb7\x86\xc6\x4c\x89\x04\x45\x33\x9a\x7d\xf7\xd8\x19\xa6\x2e\xb2\x59\x94\xf9\x69\x30\x98\x62\xba\x41\x13\x7a\x06\x8e\x0a\xba\xd5\xa6\x14\x49\x55\xeb\x81\xc6\x14\x7f\x55\x1d\xda\x4e\xa9\xa7\x85\x04\x2d\xad\x02\x53\x42\xa4\x8c\x73\xb4\x96\x3a\xfd\x82\xaa\x25\xcf\xa5\x4a\x8a\x4a\xd1\xf2\xd2\x35\xb2\x57\x9f\xe0\x28\x3a\x30\x4b\x33\x8b\x02\x3a\xcc\xb8\xe5\x61\xe9\xa6\xd8\x13\xa2\x97\xab\xd0\x8d\x3c\x91\xec\x03\x32\xc4\xb1\x0b\x23\xe3\x6f\xc3\xdc\x8b\x72\x03\xcd\x15\x32\x1d\x60\xa2\xcb\x6f\xfa\x44\x1f\x54\x34\x59\x2e\x9e\x7a\xdc\xb3\xc7\x4e\x61\x3b\xd1\xde\xfa\x32\x98\xdd\xfa\x8e\x0b\xff\xd6\x81\xce\x75\xff\x02\x00\x00\xff\xff\x5a\xbe\x93\x31\xd7\x08\x00\x00")

func migrationsSqlShared1SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlShared1Sql,
		"migrations/sql/shared/1.sql",
	)
}

func migrationsSqlShared1Sql() (*asset, error) {
	bytes, err := migrationsSqlShared1SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/shared/1.sql", size: 2263, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared2Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\xd1\x41\x4f\xf3\x20\x18\x07\xf0\x33\x7c\x8a\xe7\xb6\x2d\xef\x76\x59\xb2\x53\x4f\xbc\x2d\x46\x23\x5b\x17\xd2\x9a\xec\x44\x18\x50\x8b\x99\xa0\x40\x35\x7e\x7b\xd3\xa5\x73\xc6\xd8\x64\x99\x3d\xf4\xf4\xe7\xf9\xc1\xf3\x5f\x2c\xe0\xdf\xb3\x7d\x0c\x32\x19\xa8\x5f\x30\x61\x15\xe5\x50\x91\xff\x8c\x42\xfb\xa1\x83\x14\x5e\x76\xa9\x5d\x0a\xe5\x5d\x34\x2e\x89\x60\x5e\x3b\x13\x13\x90\xa2\x80\xc6\x07\x65\xb4\x88\xdd\xfe\xc9\xa8\x24\xac\x36\x2e\xd9\xc6\x9a\x00\x0f\x84\xe7\xb7\x84\x4f\x97\xab\xd5\x0c\x36\x35\x63\x50\xd0\x1b\x52\xb3\x0a\x26\x93\x6c\x1c\xe9\xff\xfd\x0c\x25\x93\xf5\xee\x64\x89\x56\x3a\x7d\x30\xfa\x4f\x66\xce\x29\xa9\xe8\x6f\xa8\xdf\x37\x5d\x54\x32\x19\xfd\xd3\x8f\x26\x46\xeb\x1d\x4c\x31\x1a\x3c\x00\x84\x10\x0c\xdf\x9b\x0c\xaa\x95\x61\x00\xcb\xea\x88\xce\x31\x52\x07\xdb\x6f\xca\xea\x73\x78\x2c\x7a\x7a\xc6\xf9\x0e\x68\x3c\xbc\xe5\x77\x6b\xc2\x77\x70\x4f\x77\xd3\xe1\xe0\x1c\xbe\xb0\x19\x9e\x65\x18\x7f\xaf\xb3\xf0\xef\xee\xf2\x42\x0b\x5e\x6e\x21\x2f\x59\xbd\xde\x8c\x2f\xf9\xfa\xee\x2e\x1c\x7f\x8c\x5d\x53\x52\xf6\x19\x00\x00\xff\xff\x6d\x3e\xf4\x52\xc9\x02\x00\x00")

func migrationsSqlShared2SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlShared2Sql,
		"migrations/sql/shared/2.sql",
	)
}

func migrationsSqlShared2Sql() (*asset, error) {
	bytes, err := migrationsSqlShared2SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/shared/2.sql", size: 713, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared3Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x9c\x90\xb1\xca\xc2\x30\x14\x46\xe7\xf4\x29\xee\xd6\xff\x47\x0a\x22\x6e\x4e\xb1\xa9\x38\xc4\x56\x42\xe3\x1a\x42\x7b\x69\x03\x35\xd1\x26\x45\x7c\x7b\xe9\x22\x42\x41\x6c\x97\x3b\x9e\x7b\xbe\x93\x24\xb0\xba\x9a\xa6\xd7\x01\x41\xde\x22\xca\xcb\x4c\x40\x49\xf7\x3c\x83\xf6\x59\xf7\x5a\x39\x3d\x84\x76\xa3\x2a\x67\x3d\xda\xa0\x7a\xbc\x0f\xe8\x03\x50\xc6\xa0\x73\x8d\xb1\xca\xa3\xf7\xc6\x59\x65\x6a\xb8\x50\x91\x1e\xa9\xf8\xdb\xae\xff\x21\x97\x9c\x03\xcb\x0e\x54\xf2\x12\xe2\x78\xb7\x04\x5d\xb5\xba\xeb\xd0\x36\xb8\x90\x3c\x5e\xb4\xc1\x54\x3a\x8c\x86\xcb\xdd\x09\x21\x24\xfa\x4c\xc5\xdc\xc3\xfe\xbe\x88\x89\xe2\x0c\x69\xc1\xe5\x29\x9f\x3c\x9e\x11\x66\x8a\x79\x07\x9a\x1f\xe1\xbb\xd3\x2b\x00\x00\xff\xff\xb0\x4c\xb3\x00\x17\x02\x00\x00")

func migrationsSqlShared3SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlShared3Sql,
		"migrations/sql/shared/3.sql",
	)
}

func migrationsSqlShared3Sql() (*asset, error) {
	bytes, err := migrationsSqlShared3SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/shared/3.sql", size: 535, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysqlGitkeep = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func migrationsSqlMysqlGitkeepBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlMysqlGitkeep,
		"migrations/sql/mysql/.gitkeep",
	)
}

func migrationsSqlMysqlGitkeep() (*asset, error) {
	bytes, err := migrationsSqlMysqlGitkeepBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/mysql/.gitkeep", size: 0, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql4Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd2\x41\x4b\x03\x31\x10\x05\xe0\xfb\xfe\x8a\xb9\xf5\x20\xbd\x78\x5d\x3c\xac\x26\x82\x90\x36\xa5\x26\xa0\xa7\x30\x6c\x86\x6e\x40\x27\x9a\x66\x11\xff\xbd\xb4\x14\x29\x92\xd4\x56\xf7\xb2\xec\xe9\x63\xde\x7b\x99\xcf\xe1\xea\x35\x6c\x12\x66\x02\xfb\xd6\x74\xca\xc8\x35\x98\xee\x56\x49\x18\x3e\x7d\x42\x17\x71\xcc\xc3\xb5\xeb\x23\x6f\x89\xb3\x4b\xf4\x3e\xd2\x36\x43\x27\x04\x1c\xfe\xc9\x3b\xcc\x0e\x47\x1f\x88\x7b\x02\x23\x9f\x0c\x2c\xad\x52\x6d\x5d\xdb\x7d\x89\x73\xe8\x31\x87\xc8\x13\xa1\x3f\x4e\x74\x03\xb2\x7f\x21\xbf\x57\x37\x09\xf9\x84\xd9\xd8\x95\xe8\xcc\x2f\x91\x1f\xa5\x29\x5f\x77\x33\x9b\xb5\x45\xa1\x12\xf3\x72\xa8\x16\x6d\x27\x15\xa2\xed\x9d\xf3\xb7\x5c\x68\xf1\x70\xff\x7c\xb2\x79\xfd\xd7\x49\x27\xb1\x6b\xf1\x0f\x78\x7d\xdc\x6f\xba\x39\x7e\xe7\x22\x7e\xf0\xf9\xed\x88\xb5\x5e\xc1\x9d\x56\x76\xb1\x2c\xc7\xb8\xbc\x94\xff\x90\xb5\x2e\x8e\xcd\x42\x21\x6d\xf3\x15\x00\x00\xff\xff\x99\xcd\x38\x04\xea\x03\x00\x00")

func migrationsSqlMysql4SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlMysql4Sql,
		"migrations/sql/mysql/4.sql",
	)
}

func migrationsSqlMysql4Sql() (*asset, error) {
	bytes, err := migrationsSqlMysql4SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/mysql/4.sql", size: 1002, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql5Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xbc\x93\x4f\x4b\xc4\x30\x10\xc5\xef\xf9\x14\x39\xba\xe8\x5e\xbc\xe6\x24\xb6\x87\xbd\x74\x75\xb1\xe0\x2d\xa4\xc9\xac\x3b\xa2\x89\xe6\x8f\x7f\xbe\xbd\x54\x14\x43\xe9\xac\x89\x14\x4f\x85\x32\xf3\xde\xcb\x8f\x37\xeb\x35\x3f\x7d\xc4\x3b\xaf\x22\xf0\xfe\x89\x5d\xee\xda\x8b\x9b\x96\xf7\xdd\xe6\xba\x6f\xf9\xa6\x6b\xda\x5b\x7e\x78\x37\x5e\x49\xa7\x52\x3c\x9c\x4b\x37\xec\x53\xd0\x2a\x82\x91\xe3\x0f\xb0\x11\xb5\x8a\xe8\xac\x0c\x10\xc2\xe7\xd7\x49\x34\x6f\x7c\xdb\x55\x2e\xf2\x13\xfd\x80\x60\xa3\x44\x73\xc6\x43\x1a\xee\x41\xc7\x6c\x6b\x25\xd8\x77\xba\x99\x58\xda\xd9\x30\xee\x7a\x78\x4e\x10\xa2\xd4\x68\x66\x53\x4c\xe6\x32\xcf\x95\xa8\x90\x0f\x69\x28\x93\xff\x7a\xc7\x8f\x38\x4d\x76\xea\xf1\x02\x1e\xcb\x4c\xc6\xc9\x3d\x82\x3f\x8e\x68\x42\xfd\x37\x52\xf3\xe3\xa5\xc0\x08\x33\x8a\x1b\x65\x56\x81\x8f\x70\x24\x29\x52\x96\x39\xcc\xfc\x38\x1a\xf7\x6a\x59\xb3\xdb\x5e\xfd\xef\x55\x08\x46\x9a\xfe\xb1\xf3\xa2\x58\xb0\xb0\xe5\xe5\x82\xa5\x95\x3e\xf2\xe8\x45\x5a\x4c\x27\x5e\xa4\xb7\xd5\xf2\x95\x25\x15\xec\x23\x00\x00\xff\xff\x26\xc5\x0b\x30\xb6\x05\x00\x00")

func migrationsSqlMysql5SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlMysql5Sql,
		"migrations/sql/mysql/5.sql",
	)
}

func migrationsSqlMysql5Sql() (*asset, error) {
	bytes, err := migrationsSqlMysql5SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/mysql/5.sql", size: 1462, mode: os.FileMode(420), modTime: time.Unix(1541418964, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql6Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\xce\xcf\x2b\x4e\xcd\x2b\x89\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x51\x70\x74\x71\x51\x48\x4c\x2e\x52\x08\x71\x8d\x08\x51\xf0\x0b\xf5\xf1\xb1\xe6\x0a\x0d\x70\x71\x0c\x21\xa0\x2d\xd8\x35\x04\xa4\xcd\x56\x5d\xdd\x9a\x78\xbb\x7c\xfd\x5d\x3c\xdd\x22\x91\xac\xf3\x87\x59\xc9\x85\xec\x7c\x97\xfc\xf2\x3c\xe2\x0d\x75\x09\xf2\x0f\x50\x70\xf6\xf7\x09\xf5\xf5\x03\x99\x6c\xcd\x05\x08\x00\x00\xff\xff\x40\xd9\xef\x84\x0a\x01\x00\x00")

func migrationsSqlMysql6SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlMysql6Sql,
		"migrations/sql/mysql/6.sql",
	)
}

func migrationsSqlMysql6Sql() (*asset, error) {
	bytes, err := migrationsSqlMysql6SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/mysql/6.sql", size: 266, mode: os.FileMode(420), modTime: time.Unix(1541664379, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgresGitkeep = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func migrationsSqlPostgresGitkeepBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlPostgresGitkeep,
		"migrations/sql/postgres/.gitkeep",
	)
}

func migrationsSqlPostgresGitkeep() (*asset, error) {
	bytes, err := migrationsSqlPostgresGitkeepBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/postgres/.gitkeep", size: 0, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres4Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd0\xb1\x4a\xc6\x30\x10\xc0\xf1\x39\x7d\x8a\xdb\x3a\x48\x17\xd7\x4e\xd1\xc4\x29\xb6\x52\x12\x70\x0b\x47\x73\xb4\x01\xbd\x68\x9a\x22\xbe\xbd\x28\x0e\x45\x5a\xa8\x7c\xdf\x12\x32\xfd\xef\xee\xd7\x34\x70\xf3\x1a\xa7\x8c\x85\xc0\xbd\x55\xd2\x58\x3d\x80\x95\x77\x46\xc3\xfc\x19\x32\xfa\x84\x6b\x99\x6f\xfd\x98\x78\x21\x2e\x3e\xd3\xfb\x4a\x4b\x01\xa9\x14\xfc\xfe\x29\x78\x2c\x1e\xd7\x10\x89\x47\x02\xab\x9f\x2d\x74\xce\x18\x50\xfa\x41\x3a\x63\xa1\xae\xdb\xe3\xf0\xf7\x4b\x5c\xe2\x88\x25\x26\xbe\x7e\xff\xcf\xe2\x7e\x46\x0e\x2f\x14\x7e\x06\x4c\x19\xf9\x5c\x5e\x08\x21\xaa\xad\x95\x4a\x1f\x7c\x5e\x4b\x0d\xfd\x13\xdc\xf7\xc6\x3d\x76\xfb\x57\xfd\x1f\xe8\x92\xe4\x91\xc9\xb6\xb9\x63\xd3\x56\x5f\x01\x00\x00\xff\xff\x50\x71\xaa\x54\x2e\x02\x00\x00")

func migrationsSqlPostgres4SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlPostgres4Sql,
		"migrations/sql/postgres/4.sql",
	)
}

func migrationsSqlPostgres4Sql() (*asset, error) {
	bytes, err := migrationsSqlPostgres4SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/postgres/4.sql", size: 558, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres5Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x93\x3d\x4f\xc3\x30\x10\x86\xf7\xfc\x0a\x8f\x54\xd0\x85\x35\x13\x22\x19\xba\xa4\x50\x11\x89\xcd\x72\xec\x2b\x3d\x04\x36\xf8\xce\x7c\xfc\x7b\x14\xd4\xaa\x51\x94\x0b\x8e\x98\x22\x45\xf7\xde\xf3\xe6\x89\xbd\x5e\xab\xcb\x57\x7c\x8a\x86\x41\xb5\x6f\xc5\xed\xae\xbe\x79\xa8\x55\xdb\x6c\xee\xdb\x5a\x6d\x9a\xaa\x7e\x54\x87\x6f\x17\x8d\x0e\x26\xf1\xe1\x5a\x87\x6e\x9f\xc8\x1a\x06\xa7\xfb\x17\xe0\x19\xad\x61\x0c\x5e\x13\x10\xfd\x3e\x83\x46\xf7\xa5\xb6\xcd\xc2\xa0\xba\xb0\x2f\x08\x9e\x35\xba\x2b\x45\xa9\x7b\x06\xcb\x83\xd4\xaa\x2c\x4e\xed\x26\x6a\xd9\xe0\xa9\xcf\x46\x78\x4f\x40\xac\x2d\xba\xc9\x16\xa3\xb9\x01\x73\x55\x2e\x58\x4f\xa9\xcb\x5b\x7f\xfc\x8e\xf3\x72\xd9\xec\x98\xf1\x01\x11\xf3\x20\xfd\xe4\x1e\x21\xce\x2b\x1a\x59\xff\xcb\xd4\xf4\x78\xae\x30\x01\x26\x79\x93\x60\x0b\xf4\x09\x44\xd1\xa2\x84\x1c\xca\x1c\x5e\x8e\x2a\x7c\xfa\xa2\xda\x6d\xef\xfe\x73\x2b\xca\x42\x5c\x21\x9c\xe0\x32\x3b\x70\x74\x9b\x1f\x38\xa9\x99\x29\x35\x7f\x66\x64\xd4\xfc\xef\x5f\x9c\x3b\x37\xfd\x09\x00\x00\xff\xff\xa9\xd8\x35\xfe\xaf\x04\x00\x00")

func migrationsSqlPostgres5SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlPostgres5Sql,
		"migrations/sql/postgres/5.sql",
	)
}

func migrationsSqlPostgres5Sql() (*asset, error) {
	bytes, err := migrationsSqlPostgres5SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/postgres/5.sql", size: 1199, mode: os.FileMode(420), modTime: time.Unix(1541418964, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres6Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\xf4\x09\x71\x0d\x52\x08\x71\x74\xf2\x71\x55\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\xce\xcf\x2b\x4e\xcd\x2b\x89\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x51\x70\x74\x71\x51\x48\x4c\x2e\x52\x08\x71\x8d\x08\x51\xf0\x0b\xf5\xf1\x51\x70\x71\x75\x73\x0c\xf5\x09\x51\x50\x57\xb7\xe6\xe2\x42\x36\xd6\x25\xbf\x3c\x8f\x78\x83\x5d\x82\xfc\x03\x14\x9c\xfd\x7d\x42\x7d\xfd\x40\x16\x58\x73\x01\x02\x00\x00\xff\xff\x9e\x91\x43\x31\xa2\x00\x00\x00")

func migrationsSqlPostgres6SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlPostgres6Sql,
		"migrations/sql/postgres/6.sql",
	)
}

func migrationsSqlPostgres6Sql() (*asset, error) {
	bytes, err := migrationsSqlPostgres6SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/postgres/6.sql", size: 162, mode: os.FileMode(420), modTime: time.Unix(1541662428, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTestsGitkeep = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x01\x00\x00\xff\xff\x00\x00\x00\x00\x00\x00\x00\x00")

func migrationsSqlTestsGitkeepBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTestsGitkeep,
		"migrations/sql/tests/.gitkeep",
	)
}

func migrationsSqlTestsGitkeep() (*asset, error) {
	bytes, err := migrationsSqlTestsGitkeepBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/.gitkeep", size: 0, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests1_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x94\x41\xaf\xd3\x30\x0c\xc7\xcf\xeb\xa7\xf0\xad\x9d\x48\xa5\x0d\x24\x2e\x9c\x90\xd8\x61\x12\xda\x24\xb6\xc1\x31\xca\x12\x77\x0d\xeb\x92\xe1\xa4\x0c\x84\xf8\xee\x28\xa4\xed\xba\x76\x1a\xbc\x77\x7e\xa7\xc6\xae\x23\xff\xff\x3f\x5b\xc9\x73\x78\x75\xd2\x07\x12\x1e\x61\x77\x4e\x96\xab\xcd\xe2\xd3\x16\x96\xab\xed\x3a\x99\x94\x3f\x15\x09\x6e\x45\xed\xcb\xd7\x5c\x5a\xe3\xd0\x78\x4e\xf8\xad\x46\xe7\x21\x93\xa5\xa8\x2a\x34\x07\x64\xf0\x1d\x49\x17\x1a\x89\x81\xac\x74\x28\xd2\x8a\x81\xab\xf7\x5f\x51\x7a\x06\xcd\x0d\x5e\x53\xc5\xc0\x1d\xf5\xb9\x4b\xa1\xe2\x4e\xda\x33\x32\x90\x8e\x0a\x06\xa1\x13\x1a\xaf\xa5\x08\xbf\x84\xef\x17\x86\xc8\x6a\x25\x83\x10\x8f\x3f\xfc\x34\xf9\xfc\xfe\xe3\x6e\xb1\x49\x26\x59\x3a\xcf\x3b\x31\x29\x83\x74\x9e\xb7\x82\x62\x14\x45\xc5\x73\xa3\x2a\x06\x84\x4a\x53\x8c\x0a\x51\x39\x8c\x15\x41\x51\x73\xd1\x51\x91\x32\x58\xad\xbf\x64\xd3\xee\x93\xfe\xfa\x9d\x4e\xdf\x25\x0f\x50\xf5\x6c\x68\x6b\x5e\x88\x3d\x99\x98\x43\xe7\xb4\x35\x90\x05\x2a\x63\x8b\x8d\xa2\x5b\x3f\xa1\xec\xda\x39\xea\xfe\x47\xd7\xc1\x4a\xf3\x52\x18\x55\xa1\xba\x19\xd4\x81\x84\xe9\x51\x27\x3c\xe1\x69\x1f\xe6\xd6\x9e\x78\x61\x89\x01\x12\x59\x1a\xc2\x6f\x6c\x70\x21\x25\x3a\xc7\xbd\x3d\xa2\xb9\x66\xb5\x6a\x33\x63\x87\x17\xe1\x78\xed\x50\x3d\x1c\x59\x8b\xdd\x53\x8d\x0c\xde\xbc\x9d\xcd\x22\xea\x1b\xee\x83\xd4\xdf\xa1\x0d\xa9\xc0\x7f\x6c\xef\x5d\x38\xbd\x85\x7d\x88\x45\xc8\x11\x9b\xe7\x7a\xae\xf7\x77\x1d\xa7\xf3\xd1\xd6\x75\x56\xfb\x4f\xdc\x07\x7b\x31\xc9\x9f\x00\x00\x00\xff\xff\x65\x6d\xe2\xe2\xf4\x04\x00\x00")

func migrationsSqlTests1_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests1_testSql,
		"migrations/sql/tests/1_test.sql",
	)
}

func migrationsSqlTests1_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests1_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/1_test.sql", size: 1268, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests2_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x54\x4d\x6f\xd4\x30\x10\x3d\x37\xbf\x62\x6e\xc9\x0a\x47\x2a\x8b\xc4\x85\x13\x12\x3d\x54\x42\x5b\x89\xb6\x70\xb4\xbc\xf6\x64\x63\x9a\xb5\x97\x19\x87\x82\x10\xff\x1d\xb9\xce\xe7\xa6\x5d\xa0\x17\x38\x25\x33\x99\xf1\xbc\xf7\xe6\x39\x65\x09\x2f\xf6\x76\x47\x2a\x20\xdc\x1e\xb2\xcb\xcd\xf5\xc5\x87\x1b\xb8\xdc\xdc\x5c\x65\x67\xf5\x77\x43\x4a\x7a\xd5\x86\x7a\x2d\xb5\x77\x8c\x2e\x48\xc2\x2f\x2d\x72\x80\x42\xd7\xaa\x69\xd0\xed\x50\xc0\x57\x24\x5b\x59\x24\x01\xba\xb1\xb1\xc8\x1a\x01\xdc\x6e\x3f\xa3\x0e\x02\xba\x0e\xd9\x52\x23\x80\xef\xec\x61\x48\xa1\x91\xac\xfd\x01\x05\x68\xa6\x4a\x40\x9c\x84\x2e\x58\xad\xe2\x27\x15\xa6\x85\x31\xf2\xd6\xe8\x08\x24\xe0\xb7\x20\xa0\xf2\xa4\xe3\x09\x69\x8e\xb4\x26\xb6\x46\x18\xab\xec\xe3\xdb\xf7\xb7\x17\xd7\xd9\x59\x91\xaf\xcb\x01\x67\x2e\x20\x5f\x97\x3d\xd6\x14\x25\xbc\xe9\xbd\x3b\x28\x05\x84\xc6\x52\x8a\x2a\xd5\x30\xa6\x8a\x08\xb6\x6b\x64\xaa\x72\x01\x9b\xab\x4f\xc5\x6a\x78\xe4\x3f\x7e\xa6\xaf\x09\x5a\x3c\x31\x5f\xbd\xc9\x4e\xa8\x3a\x61\x6c\xbd\xfb\x6f\xc4\xfd\xa7\x0a\xfe\x95\x62\x8c\xcc\xd6\x3b\x28\xa2\x2a\x4b\x8a\x1d\xa2\x39\x9f\x58\x36\x4e\x4e\xb8\x7f\x33\xf5\xc8\xfd\xb2\x56\xce\x34\x68\x66\x8b\xda\x91\x72\x13\xd5\x09\xf7\xb8\xdf\xc6\xbd\xf5\x6f\xb2\xf2\x24\x00\x89\x3c\x1d\x8b\xdf\xd1\x90\x4a\x6b\x64\x96\xc1\xdf\xa1\x1b\xb3\xd6\xf4\x99\x25\xc3\x7b\xc5\xb2\x65\x34\x27\x57\xd6\xcb\x1e\xa8\x45\x01\xaf\x5e\x9f\x9f\xf7\x66\x9d\x3b\x77\x9a\x7a\x58\xda\x73\xdc\xfb\xa8\x38\x13\xc3\x9e\x94\x45\xe9\x85\x36\x4f\x73\x7e\xee\x3f\x20\xee\xfb\x31\x31\xf2\x97\x0b\x43\x8e\xd6\xfd\xe3\x3b\xed\xb7\x55\xcb\x1d\xd8\x27\xcc\x3a\xa8\xb1\xbc\xd2\x93\xf6\x39\x87\x04\xfa\xe8\xc6\x8d\xc5\x0f\xa0\xa6\x7f\xf3\x77\xfe\xde\x65\xbf\x02\x00\x00\xff\xff\x48\x84\x8d\xc8\xdf\x05\x00\x00")

func migrationsSqlTests2_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests2_testSql,
		"migrations/sql/tests/2_test.sql",
	)
}

func migrationsSqlTests2_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests2_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/2_test.sql", size: 1503, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests3_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x5f\x6f\xd3\x30\x10\x7f\x5e\x3e\xc5\xbd\xa5\x15\x8e\x34\xa8\xc4\x0b\x4f\x48\xec\x61\x12\xea\x24\xb6\xc1\xa3\xe5\xda\x97\xd6\x2c\xb5\xcb\x9d\xc3\x40\x88\xef\x8e\x5c\x27\xa9\x9b\xb0\x20\x78\x61\x4f\xcd\x9d\xef\xec\xdf\x9f\xb3\x5b\x55\xf0\x62\x6f\xb7\xa4\x02\xc2\xfd\xa1\xb8\x5e\xdf\x5e\x7d\xb8\x83\xeb\xf5\xdd\x4d\x71\xb1\xfb\x6e\x48\x49\xaf\xda\xb0\x7b\x25\xb5\x77\x8c\x2e\x48\xc2\x2f\x2d\x72\x80\x85\xde\xa9\xa6\x41\xb7\x45\x01\x5f\x91\x6c\x6d\x91\x04\xe8\xc6\xc6\x22\x6b\x04\x70\xbb\xf9\x8c\x3a\x08\xe8\x3a\x64\x4b\x8d\x00\x7e\xb0\x87\x21\x85\x46\xb2\xf6\x07\x14\xa0\x99\x6a\x01\xf1\x24\x74\xc1\x6a\x15\x97\x54\xc8\x0b\x63\xe4\xad\xd1\x11\x48\xc0\x6f\x41\x40\xed\x49\xc7\x1d\xd2\x39\xd2\x9a\xd8\x9a\x60\x34\x7e\x6b\x9d\x64\x64\xb6\xde\x1d\xd1\xa4\xcc\x80\x79\x59\x7c\x7c\xfb\xfe\xfe\xea\xb6\xb8\x58\x94\xab\x6a\x48\x97\x02\xca\x55\xd5\xd3\x49\x51\xa2\x94\xbe\xbb\xb3\x52\x40\x68\x2c\xa5\xa8\x56\x0d\x63\xaa\x88\x7c\xba\x46\xa6\xba\x14\xb0\xbe\xf9\xb4\x58\x0e\x3f\xe5\x8f\x9f\x69\x35\xa1\x8f\x3b\xa6\x78\x0c\x39\xcb\x66\xf8\x96\x6f\x8a\x19\x8f\x32\xfd\xe2\x26\xcf\xc7\xaa\x31\xb9\xff\x2f\x7f\x12\xb6\x43\x54\x59\xf3\x77\xca\x76\x7d\xb0\x88\xea\x4d\xa5\xe8\x70\x9e\xb3\x8c\x65\x27\x3c\x89\xcd\x1f\x4e\x1d\xdd\x39\xb9\x53\xce\x34\x68\xce\x0c\xdd\x92\x72\x99\x3b\x84\x7b\xdc\x6f\xa2\xbf\xfd\x97\xac\x3d\x09\x40\x22\x4f\x63\x93\x7a\x43\x94\xd6\xc8\x2c\x83\x7f\x40\x77\xca\x5a\xd3\x67\xa6\x0c\x1f\x15\xcb\x96\x71\xde\xc8\xde\x8c\x40\x2d\x0a\x58\xbd\xbe\xbc\xec\x0d\x38\x77\x23\x4f\x1d\xad\xfc\x97\x29\xff\xad\x38\xd9\x60\xcf\xca\xa2\xf4\x44\x9b\xa7\x39\xcf\xbc\x3c\xf3\x72\x1c\xaf\xfa\x54\x8c\xf2\xe5\x64\x4c\x4f\x03\x9d\x3d\x13\xf3\xaa\xf8\x4d\xdd\x72\x07\xf6\x89\x61\x1d\xd4\x98\x5e\xfd\xac\xfd\x9c\xc3\xf0\x3e\xe5\xf7\xf0\x54\x7c\x04\x95\xff\x87\xbc\xf3\x8f\xae\xf8\x15\x00\x00\xff\xff\x74\xe4\x47\x5f\x55\x06\x00\x00")

func migrationsSqlTests3_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests3_testSql,
		"migrations/sql/tests/3_test.sql",
	)
}

func migrationsSqlTests3_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests3_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/3_test.sql", size: 1621, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests4_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\xcd\x6e\xd4\x30\x10\x3e\x37\x4f\x31\xb7\xec\x0a\x47\x2a\x50\x71\xe1\x84\x44\x0f\x95\xd0\x56\xa2\x2d\x1c\x2d\xaf\x3d\xd9\x35\xcd\xda\xcb\xd8\xa6\x20\xc4\xbb\x23\xaf\xf3\xe3\x24\x34\x40\x4f\x3d\x6d\x66\x32\x1e\x7f\x3f\x99\xd9\xaa\x82\x17\x07\xbd\x23\xe1\x11\xee\x8e\xc5\xd5\xe6\xe6\xf2\xe3\x2d\x5c\x6d\x6e\xaf\x8b\xb3\xfd\x0f\x45\x82\x5b\x11\xfc\xfe\x15\x97\xd6\x38\x34\x9e\x13\x7e\x0d\xe8\x3c\xac\xe4\x5e\x34\x0d\x9a\x1d\x32\xf8\x86\xa4\x6b\x8d\xc4\x40\x36\x3a\x16\x69\xc5\xc0\x85\xed\x17\x94\x9e\x41\x7b\x82\x07\x6a\x18\xb8\x7b\x7d\xec\x53\xa8\xb8\x93\xf6\x88\x0c\xa4\xa3\x9a\x41\xbc\x09\x8d\xd7\x52\xc4\x57\xc2\xe7\x85\x31\xb2\x5a\xc9\x08\xc4\xe3\x77\xcf\xa0\xb6\x24\x63\x87\x74\x0f\xd7\x2a\x1e\x4d\x30\x1a\xbb\xd3\x86\x3b\x74\x4e\x5b\x73\x42\x93\x32\x19\xe6\xbc\x33\x17\x41\x69\x34\x12\xd7\xc5\xa7\x77\x1f\xee\x2e\x6f\x8a\xb3\x55\x79\x51\xf5\xd5\x25\x83\xf2\xa2\xea\x58\xa6\x28\x31\x4d\xcf\x2d\x84\x14\x10\x2a\x4d\x29\xaa\x45\xe3\x30\x55\x44\x9a\xed\x41\x47\x75\xc9\x60\x73\xfd\x79\xb5\xee\x7f\xca\x9f\xbf\xd2\xdb\x44\x2a\x76\x4c\xf1\x94\x49\x96\x9d\xe2\x13\x41\x95\xeb\xb7\xc5\x82\x87\x99\xbe\xb1\xdb\xf3\xb1\x72\xee\xd7\x33\xb5\x27\x09\xdf\x02\xad\x3a\x3b\xfe\x5b\xf9\xb6\x01\xac\x22\xd5\xb9\x54\x2d\xe0\x31\xdd\x58\x36\x00\x4b\xb4\xfe\x72\xeb\x64\x66\xf9\x5e\x18\xd5\xa0\x1a\x19\xbe\x23\x61\x32\xf7\x08\x0f\x78\xd8\x46\xff\xbb\x27\x5e\x5b\x62\x80\x44\x96\xa6\x26\x76\x86\x09\x29\xd1\x39\xee\xed\x3d\x9a\x21\xab\x55\x97\x99\x33\x7c\x10\x8e\x07\x87\x6a\xb8\xff\x5f\x6d\xee\xac\xf2\x14\x90\xc1\xeb\x37\xe7\xe7\x9d\x3d\x63\xaf\xf2\xd4\x60\xf4\x53\x67\xe4\x8f\xd2\x65\x63\xb1\x28\x9a\x90\x33\xe5\x96\x14\x79\x74\xaf\x2d\xeb\x72\xda\x18\x73\x55\xca\x97\xb3\xaf\x79\x90\x23\xdb\x36\xcb\xaa\xd8\x6d\x1d\x5c\x0b\xf6\x91\x4f\xb9\x57\x63\xbe\x38\xb2\xe3\x63\x0e\xfd\x9a\xcb\xc7\x75\x28\x3e\x81\xca\xff\xa1\xde\xdb\x07\x53\xfc\x0e\x00\x00\xff\xff\xe2\xb3\x8e\x96\xb3\x06\x00\x00")

func migrationsSqlTests4_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests4_testSql,
		"migrations/sql/tests/4_test.sql",
	)
}

func migrationsSqlTests4_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests4_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/4_test.sql", size: 1715, mode: os.FileMode(420), modTime: time.Unix(1541173634, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests5_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\xcd\x6e\xd4\x30\x10\x3e\x37\x4f\x31\xb7\xec\x0a\x47\x2a\xa0\x72\xe1\x84\x44\x0f\x95\xd0\x56\xa2\x2d\x1c\x2d\xaf\x3d\xd9\x35\xcd\xda\xcb\xd8\xa6\x20\xc4\xbb\x23\xaf\xf3\xe3\x24\x34\x40\x4f\x3d\x6d\x66\x32\x1e\x7f\x3f\x99\xd9\xaa\x82\x17\x07\xbd\x23\xe1\x11\xee\x8e\xc5\xd5\xe6\xe6\xf2\xe3\x2d\x5c\x6d\x6e\xaf\x8b\xb3\xfd\x0f\x45\x82\x5b\x11\xfc\xfe\x15\x97\xd6\x38\x34\x9e\x13\x7e\x0d\xe8\x3c\xac\xe4\x5e\x34\x0d\x9a\x1d\x32\xf8\x86\xa4\x6b\x8d\xc4\x40\x36\x3a\x16\x69\xc5\xc0\x85\xed\x17\x94\x9e\x41\x7b\x82\x07\x6a\x18\xb8\x7b\x7d\xec\x53\xa8\xb8\x93\xf6\x88\x0c\xa4\xa3\x9a\x41\xbc\x09\x8d\xd7\x52\xc4\x57\xc2\xe7\x85\x31\xb2\x5a\xc9\x08\xc4\xe3\x77\xcf\xa0\xb6\x24\x63\x87\x74\x0f\xd7\x2a\x1e\x4d\x30\x1a\xbb\xd3\x86\x3b\x74\x4e\x5b\x73\x42\x93\x32\x19\xe6\xbc\x33\x17\x41\x69\x34\x12\xd7\xc5\xa7\x77\x1f\xee\x2e\x6f\x8a\xb3\x55\x79\x51\xf5\xd5\x25\x83\xf2\xa2\xea\x58\xa6\x28\x31\x4d\xcf\x2d\x84\x14\x10\x2a\x4d\x29\xaa\x45\xe3\x30\x55\x44\x9a\xed\x41\x47\x75\xc9\x60\x73\xfd\x79\xb5\xee\x7f\xca\x9f\xbf\xd2\xdb\x44\x2a\x76\x4c\xf1\x94\x49\x96\x9d\xe2\x13\x41\x95\xeb\xb7\xc5\x82\x87\x99\xbe\xb1\xdb\xf3\xb1\x72\xee\xd7\x33\xb5\x27\x09\xdf\x02\xad\x3a\x3b\xfe\x5b\xf9\xb6\x01\xac\x22\xd5\xb9\x54\x2d\xe0\x31\xdd\x58\x36\x00\x4b\xb4\xfe\x72\xeb\x64\x66\xf9\x5e\x18\xd5\xa0\x1a\x19\xbe\x23\x61\x32\xf7\x08\x0f\x78\xd8\x46\xff\xbb\x27\x5e\x5b\x62\x80\x44\x96\xa6\x26\x76\x86\x09\x29\xd1\x39\xee\xed\x3d\x9a\x21\xab\x55\x97\x99\x33\x7c\x10\x8e\x07\x87\x6a\xb8\xff\x5f\x6d\xee\xac\xf2\x14\x90\xc1\xeb\x37\xe7\xe7\x9d\x3d\x63\xaf\xf2\xd4\x60\xf4\x53\x67\xe4\x8f\xd2\x65\x63\xb1\x28\x9a\x90\x33\xe5\x96\x14\x79\x74\xaf\x2d\xeb\x72\xda\x18\x73\x55\xca\x97\xb3\xaf\x79\x90\x23\xdb\x36\xcb\xaa\xd8\x6d\x1d\x5c\x0b\xf6\x91\x4f\xb9\x57\x63\xbe\x38\xb2\xe3\x63\x0e\xfd\x9a\xcb\xc7\x75\x28\x3e\x81\xca\xff\xa1\xde\xdb\x07\x53\xfc\x0e\x00\x00\xff\xff\x40\x72\xe3\x11\xb3\x06\x00\x00")

func migrationsSqlTests5_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests5_testSql,
		"migrations/sql/tests/5_test.sql",
	)
}

func migrationsSqlTests5_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests5_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/5_test.sql", size: 1715, mode: os.FileMode(420), modTime: time.Unix(1541418964, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests6_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\x5f\x6f\xd3\x30\x10\x7f\x5e\x3e\xc5\xbd\xa5\x15\x8e\x34\x40\xea\x0b\x4f\x48\xec\x61\x12\xea\x24\xb6\xc1\xa3\xe5\xda\x97\xd6\x2c\xb5\xcb\xd9\x66\x20\xc4\x77\x47\xae\x9d\xc4\x6d\x58\x60\x4f\xf0\xd4\xdc\xe5\xee\xfc\xfb\x13\x5f\x9b\x06\x5e\xec\xf5\x96\x84\x47\xb8\x3f\x54\xd7\xeb\xdb\xab\x0f\x77\x70\xbd\xbe\xbb\xa9\x2e\x76\xdf\x15\x09\x6e\x45\xf0\xbb\x57\x5c\x5a\xe3\xd0\x78\x4e\xf8\x25\xa0\xf3\xb0\x90\x3b\xd1\x75\x68\xb6\xc8\xe0\x2b\x92\x6e\x35\x12\x03\xd9\xe9\x58\xa4\x15\x03\x17\x36\x9f\x51\x7a\x06\xb9\x83\x07\xea\x18\xb8\x07\x7d\x18\x52\xa8\xb8\x93\xf6\x80\x0c\xa4\xa3\x96\x41\x3c\x09\x8d\xd7\x52\xc4\x57\xc2\x97\x85\x31\xb2\x5a\xc9\x08\xc4\xe3\x37\xcf\xa0\xb5\x24\xe3\x84\x74\x0e\xd7\x2a\xb6\x26\x18\x9d\xdd\x6a\xc3\x1d\x3a\xa7\xad\x39\xa2\x49\x99\x02\x73\x39\x99\x8b\xa0\x34\x1a\x89\x0c\x84\xa4\x65\xf5\xf1\xed\xfb\xfb\xab\xdb\xea\x62\x51\xaf\x9a\xa1\xa5\x66\x50\xaf\x9a\x9e\x6a\x8a\x12\xdd\xf4\x9c\x71\xa4\x80\x50\x69\x4a\x51\x2b\x3a\x87\xa9\x22\x72\xcd\x8d\x8e\xda\x9a\xc1\xfa\xe6\xd3\x62\x39\xfc\xd4\x3f\x7e\xa6\xb7\x89\x59\x9c\x98\xe2\x73\x3a\x45\xf6\x1c\x9f\x08\xf9\xad\x90\x54\x2f\xdf\x54\x33\x8e\x16\x6a\xc7\xb1\xff\x8f\xb1\x53\xf7\x7e\x6b\xd6\xbf\xf7\x29\x39\x90\x81\x36\xbd\x2f\xd1\x82\x67\x29\x9f\x07\xc0\x22\x52\x9d\x4a\x95\x01\x9f\xd2\x8d\x65\x23\xb0\x44\xeb\x0f\xa7\x9e\xdd\x60\xbe\x13\x46\x75\xa8\x4e\x0c\xdf\x92\x30\x85\x7b\x84\x7b\xdc\x6f\xa2\xff\xfd\x13\x6f\x2d\x31\x40\x22\x4b\xe7\x26\xf6\x86\x09\x29\xd1\x39\xee\xed\x03\x9a\x31\xab\x55\x9f\x99\x32\x7c\x14\x8e\x07\x87\x6a\x3c\xff\x6f\x6d\xee\xad\xf2\x14\x90\xc1\xeb\xd5\xe5\x65\x6f\xcf\xa9\x57\x65\x6a\x34\xfa\xd9\x4e\xcd\x49\x57\x5c\x8b\x59\xd1\x84\x9c\x28\x37\xa7\xc8\x93\x5b\x6e\x5e\x97\xe3\xea\x98\xaa\x52\xbf\x9c\x7c\xcd\xa3\x1c\xc5\xda\x99\x57\xc5\x6e\xda\xe0\x32\xd8\x27\x3e\xe5\x41\x8d\xe9\xe2\x28\xda\x4f\x39\x0c\xfb\xae\xbc\xae\x63\xf1\x11\x54\xf9\x7f\xf5\xce\x3e\x9a\xea\x57\x00\x00\x00\xff\xff\x56\x91\x74\xf2\xc1\x06\x00\x00")

func migrationsSqlTests6_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests6_testSql,
		"migrations/sql/tests/6_test.sql",
	)
}

func migrationsSqlTests6_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests6_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/6_test.sql", size: 1729, mode: os.FileMode(420), modTime: time.Unix(1541664901, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"migrations/sql/shared/.gitkeep":   migrationsSqlSharedGitkeep,
	"migrations/sql/shared/1.sql":      migrationsSqlShared1Sql,
	"migrations/sql/shared/2.sql":      migrationsSqlShared2Sql,
	"migrations/sql/shared/3.sql":      migrationsSqlShared3Sql,
	"migrations/sql/mysql/.gitkeep":    migrationsSqlMysqlGitkeep,
	"migrations/sql/mysql/4.sql":       migrationsSqlMysql4Sql,
	"migrations/sql/mysql/5.sql":       migrationsSqlMysql5Sql,
	"migrations/sql/mysql/6.sql":       migrationsSqlMysql6Sql,
	"migrations/sql/postgres/.gitkeep": migrationsSqlPostgresGitkeep,
	"migrations/sql/postgres/4.sql":    migrationsSqlPostgres4Sql,
	"migrations/sql/postgres/5.sql":    migrationsSqlPostgres5Sql,
	"migrations/sql/postgres/6.sql":    migrationsSqlPostgres6Sql,
	"migrations/sql/tests/.gitkeep":    migrationsSqlTestsGitkeep,
	"migrations/sql/tests/1_test.sql":  migrationsSqlTests1_testSql,
	"migrations/sql/tests/2_test.sql":  migrationsSqlTests2_testSql,
	"migrations/sql/tests/3_test.sql":  migrationsSqlTests3_testSql,
	"migrations/sql/tests/4_test.sql":  migrationsSqlTests4_testSql,
	"migrations/sql/tests/5_test.sql":  migrationsSqlTests5_testSql,
	"migrations/sql/tests/6_test.sql":  migrationsSqlTests6_testSql,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"migrations": &bintree{nil, map[string]*bintree{
		"sql": &bintree{nil, map[string]*bintree{
			"mysql": &bintree{nil, map[string]*bintree{
				".gitkeep": &bintree{migrationsSqlMysqlGitkeep, map[string]*bintree{}},
				"4.sql":    &bintree{migrationsSqlMysql4Sql, map[string]*bintree{}},
				"5.sql":    &bintree{migrationsSqlMysql5Sql, map[string]*bintree{}},
				"6.sql":    &bintree{migrationsSqlMysql6Sql, map[string]*bintree{}},
			}},
			"postgres": &bintree{nil, map[string]*bintree{
				".gitkeep": &bintree{migrationsSqlPostgresGitkeep, map[string]*bintree{}},
				"4.sql":    &bintree{migrationsSqlPostgres4Sql, map[string]*bintree{}},
				"5.sql":    &bintree{migrationsSqlPostgres5Sql, map[string]*bintree{}},
				"6.sql":    &bintree{migrationsSqlPostgres6Sql, map[string]*bintree{}},
			}},
			"shared": &bintree{nil, map[string]*bintree{
				".gitkeep": &bintree{migrationsSqlSharedGitkeep, map[string]*bintree{}},
				"1.sql":    &bintree{migrationsSqlShared1Sql, map[string]*bintree{}},
				"2.sql":    &bintree{migrationsSqlShared2Sql, map[string]*bintree{}},
				"3.sql":    &bintree{migrationsSqlShared3Sql, map[string]*bintree{}},
			}},
			"tests": &bintree{nil, map[string]*bintree{
				".gitkeep":   &bintree{migrationsSqlTestsGitkeep, map[string]*bintree{}},
				"1_test.sql": &bintree{migrationsSqlTests1_testSql, map[string]*bintree{}},
				"2_test.sql": &bintree{migrationsSqlTests2_testSql, map[string]*bintree{}},
				"3_test.sql": &bintree{migrationsSqlTests3_testSql, map[string]*bintree{}},
				"4_test.sql": &bintree{migrationsSqlTests4_testSql, map[string]*bintree{}},
				"5_test.sql": &bintree{migrationsSqlTests5_testSql, map[string]*bintree{}},
				"6_test.sql": &bintree{migrationsSqlTests6_testSql, map[string]*bintree{}},
			}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
