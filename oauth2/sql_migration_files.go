// Code generated by go-bindata. (@generated) DO NOT EDIT.
// sources:
// migrations/sql/shared/1.sql
// migrations/sql/shared/2.sql
// migrations/sql/shared/3.sql
// migrations/sql/shared/4.sql
// migrations/sql/mysql/.gitkeep
// migrations/sql/mysql/5.sql
// migrations/sql/mysql/6.sql
// migrations/sql/mysql/7.sql
// migrations/sql/postgres/.gitkeep
// migrations/sql/postgres/5.sql
// migrations/sql/postgres/6.sql
// migrations/sql/postgres/7.sql
// migrations/sql/tests/.gitkeep
// migrations/sql/tests/1_test.sql
// migrations/sql/tests/2_test.sql
// migrations/sql/tests/3_test.sql
// migrations/sql/tests/4_test.sql
// migrations/sql/tests/5_test.sql
// migrations/sql/tests/6_test.sql
// migrations/sql/tests/7_test.sql
package oauth2

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

var _migrationsSqlShared1Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x94\x4d\x6b\xc2\x30\x18\xc7\xcf\x09\xe4\x3b\x3c\x47\x65\xf3\x22\x78\xf2\xd4\xcd\x0a\xb2\x4e\xa5\x56\x98\xa7\xf2\x90\x3c\xda\xc0\x9a\xb8\x24\xce\xed\xdb\x8f\x38\xc7\x14\x5f\x18\x03\x6f\xed\xf5\xff\xd2\xf0\x0b\xf9\x77\x3a\x70\x57\xeb\x95\xc3\x40\x30\x5f\x0b\xfe\x98\xa7\x49\x91\x42\x91\x3c\x64\x29\x8c\x86\x30\x9e\x14\x90\xbe\x8c\x66\xc5\x0c\xaa\x4f\xe5\xb0\xb4\xb8\x09\x55\xb7\x44\x29\xc9\x7b\x68\x09\xce\xbc\x5e\x19\x0c\x1b\x47\xb0\xfb\xd8\x3b\x3a\x59\xa1\x6b\x75\x7b\xbd\xf6\x2e\x3f\x9e\x67\x19\x4c\xf3\xd1\x73\x92\x2f\xe0\x29\x5d\xdc\x0b\xce\x1c\xbd\x6d\xc8\x87\x52\x2b\x00\x06\x70\x36\x73\xe0\x23\x55\x62\x00\x60\x41\xd7\xe4\x03\xd6\xeb\xdf\xe2\x41\x3a\x4c\xe6\x59\x01\xc6\x6e\x5b\xed\x18\x91\xaf\x9a\xcc\xbe\x99\x01\x04\xfa\x08\x47\x95\x5e\xda\x35\x45\x8d\xc5\xe3\x9e\xc8\x2b\x87\x26\xfe\xef\xdb\xc6\x4e\xf4\xa5\x75\x75\xa9\x30\xe0\xa5\x7a\xf2\x5e\x5b\xf3\x63\x39\xd2\x05\x6f\xf7\xff\xce\xd8\xd1\xd2\x91\xaf\x1a\xc8\xb7\x84\x2c\xad\xa2\x86\xf0\x2d\x09\x5b\xad\x64\x43\xf8\x3f\x84\x05\x3f\x9c\xe7\x81\xdd\x1a\xc1\x07\xf9\x64\xba\x67\x7e\x66\x90\xfb\x97\x0d\xfb\x35\xb9\xe2\x88\x4f\xe1\x8a\x1c\xef\xb1\x2f\xf8\x57\x00\x00\x00\xff\xff\x80\x97\x61\xb9\x32\x06\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/shared/1.sql", size: 1586, mode: os.FileMode(438), modTime: time.Unix(1540938328, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared2Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\xd0\xc1\x0a\x82\x30\x00\xc6\xf1\xbb\xe0\x3b\x7c\x37\x8b\xf0\x22\x78\xea\xb4\x9a\x9d\x96\x86\x6c\x67\x59\x73\x35\x83\x5a\x6c\x5a\xf4\xf6\x81\x10\x74\x88\x6a\xe0\x03\x7c\xff\x0f\x7e\x69\x8a\xc5\xb9\x3b\x3a\xd9\x6b\x88\x6b\x1c\x11\xc6\x8b\x1a\x9c\xac\x58\x01\xf3\x68\x9d\x6c\xac\x1c\x7a\x93\x35\x52\x29\xed\x3d\x08\xa5\xf0\xc3\xfe\xa4\x55\x8f\x9b\x74\xca\x48\x37\xcb\xf2\x7c\x8e\xb2\xe2\x28\x05\x63\xa0\xc5\x86\x08\xc6\x91\x24\xcb\x2f\x39\xa7\x0f\x4e\x7b\x33\x59\x4f\xd9\x56\x4f\x16\xb3\x5d\xab\x42\x63\x71\xf4\x4e\x49\xed\xfd\xf2\x1b\x13\xb4\xae\x76\x58\x57\x4c\x6c\xcb\xd7\xd7\x3f\x68\x81\xb3\xd1\x26\x70\x33\x12\x7c\xde\x3c\x03\x00\x00\xff\xff\x67\x0b\xde\x9a\x33\x02\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/shared/2.sql", size: 563, mode: os.FileMode(438), modTime: time.Unix(1540940530, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared3Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x74\x91\xc1\x4e\xf3\x30\x10\x84\xcf\xb6\xe4\x77\xd8\x63\xa2\xff\xef\xa5\x52\x4f\x3d\x05\xe2\x4a\x11\xa1\xad\xd2\x44\xa2\xa7\xc8\x24\x4b\x62\x20\x76\xb0\x37\x14\xde\x1e\xa5\x2d\xa2\x28\xc5\xd7\x99\xf9\xc6\xab\x99\xcd\xe0\x5f\xa7\x1b\xa7\x08\xa1\xe8\x05\xbf\xcd\x64\x94\x4b\xc8\xa3\x9b\x54\x42\xb2\x82\xf5\x26\x07\xf9\x90\xec\xf2\x1d\xb4\x9f\xb5\x53\xa5\x55\x03\xb5\xf3\xb2\x7f\xa9\x10\x02\xc1\x99\xd7\x8d\x51\x34\x38\x84\xe3\x63\xef\xca\x55\xad\x72\xc1\x7c\xb1\x08\x8f\xe9\x75\x91\xa6\xb0\xcd\x92\xfb\x28\xdb\xc3\x9d\xdc\xff\x17\x9c\x39\x7c\x1b\xd0\x53\xa9\x6b\x00\x06\x70\x35\x73\xe1\xc3\xba\x54\x04\xc0\x48\x77\xe8\x49\x75\xfd\x0f\x38\x96\xab\xa8\x48\x73\x30\xf6\x10\x84\x63\xa4\x7a\xd5\x68\xce\x64\x06\x40\xf8\x41\xbf\x90\xbe\xb2\x3d\x8e\x1a\x1b\xbf\x3b\x91\x1b\xa7\xcc\xd8\x77\xb2\xb1\x89\xfe\x64\x5d\x57\xd6\x8a\xd4\x5f\x78\xf4\x5e\x5b\xf3\x6d\x99\xea\xc3\xe3\x33\x56\x04\xa7\xfa\xab\x87\x0b\x1e\x2e\x05\x17\xfc\x72\x9a\xd8\x1e\x8c\xe0\x71\xb6\xd9\x9e\xa7\x99\x8c\xb1\x14\xfc\x2b\x00\x00\xff\xff\xc6\x6c\x52\x8f\xcc\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/shared/3.sql", size: 460, mode: os.FileMode(438), modTime: time.Unix(1540938501, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlShared4Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\xd0\xb1\x0e\x82\x30\x10\xc6\xf1\x9d\x84\x77\xf8\x76\xc3\xe2\xea\x54\x6c\x9d\x4e\x6a\x48\x3b\x93\xa6\x9c\x42\x8c\x96\x14\xd4\xf8\xf6\x26\xc4\xc1\xc1\x28\x84\x07\xf8\xfe\x77\xf9\x65\x19\x56\x97\xf6\x14\xdd\xc0\xb0\x5d\x9a\x08\x32\xaa\x84\x11\x39\x29\x34\xcf\x3a\xba\x2a\xb8\xdb\xd0\xac\x2b\xe7\x3d\xf7\x3d\x84\x94\x70\x7e\x68\xef\x8c\x5c\x6b\x42\xa1\x0d\x0a\x4b\x04\xa9\x76\xc2\x92\x81\x29\xad\xda\xfc\xe8\x44\x3e\x46\xee\x9b\xe5\x21\x1f\x6a\x5e\x5e\x09\x6d\xed\x97\x57\xba\xb3\x9f\xfe\x4b\x9a\x7c\xa2\xcb\xf0\xb8\xfe\x67\x97\xa5\x3e\x60\xab\xc9\xee\x8b\xf7\x8d\x29\xc6\xf3\x56\x23\xe8\xbc\xc9\xa8\x37\x6f\x32\x52\x7d\x9d\xbc\x02\x00\x00\xff\xff\x05\xbc\xf0\x06\x8b\x02\x00\x00")

func migrationsSqlShared4SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlShared4Sql,
		"migrations/sql/shared/4.sql",
	)
}

func migrationsSqlShared4Sql() (*asset, error) {
	bytes, err := migrationsSqlShared4SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/shared/4.sql", size: 651, mode: os.FileMode(438), modTime: time.Unix(1540938602, 0)}
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

	info := bindataFileInfo{name: "migrations/sql/mysql/.gitkeep", size: 0, mode: os.FileMode(438), modTime: time.Unix(1540902707, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql5Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\x0e\x72\x75\x0c\x71\x55\x08\xf5\xf3\x0c\x0c\x75\x55\xf0\xf4\x73\x71\x8d\x50\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\x4c\x4e\x4e\x2d\x2e\x8e\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x89\xcf\x4c\x89\xcf\x4c\xa9\x50\xf0\xf7\xc3\xa6\x4a\x41\x03\xa1\x4c\xd3\x9a\xb0\xd9\x45\xa9\x69\x45\xa9\xc5\x19\x84\x0c\x87\x2a\x43\x33\x9d\x0b\xd9\x27\x2e\xf9\xe5\x79\x5c\x2e\x41\xfe\x01\x94\x7a\xc1\x1a\xa7\x29\xa4\x39\xd6\x9a\x0b\x10\x00\x00\xff\xff\x5b\x41\x11\xcf\x69\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/mysql/5.sql", size: 361, mode: os.FileMode(438), modTime: time.Unix(1541448703, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql6Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\x0e\x72\x75\x0c\x71\x55\xf0\xf4\x73\x71\x8d\x50\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\x4c\x4e\x4e\x2d\x2e\x8e\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x49\x4d\x89\x4f\x2c\x89\xcf\x4c\xa9\x50\xf0\xf7\xc3\xa6\x4e\x41\x03\x59\xa1\xa6\x35\x17\x17\xb2\x45\x2e\xf9\xe5\x79\x5c\x2e\x41\xfe\x01\x94\x5b\x64\xcd\x05\x08\x00\x00\xff\xff\xf1\x5f\x27\x42\xc2\x00\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/mysql/6.sql", size: 194, mode: os.FileMode(438), modTime: time.Unix(1541448703, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlMysql7Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x95\xc1\x6a\x83\x40\x10\x86\xef\x81\xbc\xc3\xde\x72\x28\xb9\xf4\x1a\x7a\xb0\x5d\x0b\x05\x13\x43\xba\x42\x7b\x92\x65\x77\x1a\xa5\x54\xd3\x55\x29\x7d\xfb\xc2\x42\xb6\x8b\x6e\x12\x67\xf4\x01\xfe\x9d\x6f\xbe\xdf\x64\xd6\x6b\x76\xf7\x55\x1e\x8d\x6c\x81\x65\xa7\xe5\x22\x4a\x44\x7c\x60\x22\x7a\x4c\x62\x56\xfc\x6a\x23\xf3\x5a\x76\x6d\x71\x9f\x4b\xa5\xa0\x69\x58\xc4\x39\x33\xf0\xdd\x41\xd3\x82\xce\x65\xa7\x4b\xa8\x14\x30\x11\xbf\x09\xb6\xcb\x92\x64\xb3\x5c\x64\x7b\x1e\x89\x70\xfa\x35\x16\x81\xf4\xc3\x6a\xb5\xb9\x3d\x79\x9b\xf2\x97\xe7\xf7\xcb\xc3\x53\x07\x30\x66\x87\xa3\x91\x15\x75\x83\x7e\x16\xc5\x7f\x61\xf0\x3f\xfd\x95\x97\x0c\x7c\x18\x68\x0a\x6a\x09\xe7\x38\xa9\x85\x73\x78\x86\x1a\xfc\x35\xd0\x3d\xf8\x4b\x20\x8b\xe8\xad\x30\xa5\x09\x55\x6b\xa0\xd6\x60\xb3\xa4\x0e\x6c\x72\x86\x02\x1c\x3d\xda\xbe\x63\x47\xaa\xf7\xc9\xa7\x78\xaf\x4b\xad\xa8\xde\x6d\x96\xe4\xdd\x26\x67\xf0\xee\xe8\xd1\xde\x1d\x3b\xd2\xbb\x4f\x3e\xc5\xfb\xe9\x53\x91\xbf\x77\x9b\x25\x79\xb7\xc9\x19\xbc\x3b\x7a\xb4\x77\xc7\x8e\xf4\xee\x93\x8f\xf0\xee\x5f\x61\x5e\xff\x54\xb7\xaf\x09\x3f\xa4\x7b\xf6\x94\x26\xd9\x76\x17\x50\x33\xe2\x1c\xf9\x0f\xf4\x09\xc7\x5d\x21\x3a\x42\xe8\x05\x1c\x83\xfd\x3f\xa1\x03\x0c\xe2\xb8\xe9\xf6\x57\x45\x9f\x3e\x88\xe3\xa6\xdb\x6f\x8b\x3e\x7d\x10\x0f\x4c\xff\x0b\x00\x00\xff\xff\x7a\x4b\x6a\x9a\x16\x0a\x00\x00")

func migrationsSqlMysql7SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlMysql7Sql,
		"migrations/sql/mysql/7.sql",
	)
}

func migrationsSqlMysql7Sql() (*asset, error) {
	bytes, err := migrationsSqlMysql7SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/mysql/7.sql", size: 2582, mode: os.FileMode(438), modTime: time.Unix(1541012702, 0)}
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

	info := bindataFileInfo{name: "migrations/sql/postgres/.gitkeep", size: 0, mode: os.FileMode(438), modTime: time.Unix(1540902707, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres5Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\x0e\x72\x75\x0c\x71\x55\x08\xf5\xf3\x0c\x0c\x75\x55\xf0\xf4\x73\x71\x8d\x50\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\x4c\x4e\x4e\x2d\x2e\x8e\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x89\xcf\x4c\x89\xcf\x4c\xa9\x50\xf0\xf7\xc3\xa6\x4a\x41\x03\xa1\x4c\xd3\x9a\xb0\xd9\x45\xa9\x69\x45\xa9\xc5\x19\x84\x0c\x87\x2a\x43\x33\x9d\x0b\xd9\x27\x2e\xf9\xe5\x79\x5c\x2e\x41\xfe\x01\x44\x7b\xc1\x1a\xa7\x72\xec\xae\xb2\xe6\x02\x04\x00\x00\xff\xff\x35\xd0\xf0\x59\x3a\x01\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/postgres/5.sql", size: 314, mode: os.FileMode(438), modTime: time.Unix(1541448703, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres6Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\xd5\x55\xd0\xce\xcd\x4c\x2f\x4a\x2c\x49\x55\x08\x2d\xe0\x72\x0e\x72\x75\x0c\x71\x55\xf0\xf4\x73\x71\x8d\x50\xc8\xa8\x4c\x29\x4a\x8c\xcf\x4f\x2c\x2d\xc9\x30\x8a\x4f\x4c\x4e\x4e\x2d\x2e\x8e\x2f\x4a\x2d\x2c\x4d\x2d\x2e\x49\x4d\x89\x4f\x2c\x89\xcf\x4c\xa9\x50\xf0\xf7\xc3\xa6\x4e\x41\x03\x59\xa1\xa6\x35\x17\x17\xb2\x45\x2e\xf9\xe5\x79\x5c\x2e\x41\xfe\x01\x24\x58\x64\xcd\x05\x08\x00\x00\xff\xff\xd2\x18\x3e\xa9\xab\x00\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/postgres/6.sql", size: 171, mode: os.FileMode(438), modTime: time.Unix(1541448703, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlPostgres7Sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x93\xb1\x4e\x86\x30\x14\x46\x77\x9e\xe2\x6e\xff\x60\x58\x5c\x9d\xaa\xad\x53\x05\x43\xda\xc4\x8d\x34\xed\x15\x88\x91\x62\x0b\x31\xbe\xbd\x09\x03\x21\x82\xa1\x94\x3e\xc0\x3d\xe7\xa4\xf9\x9a\xe7\x70\xf7\xd9\x35\x4e\x8d\x08\x72\xc8\x08\x17\xac\x02\x41\x1e\x39\x83\xf6\xc7\x38\x55\x5b\x35\x8d\xed\x7d\xad\xb4\x46\xef\x81\x50\x0a\x0e\xbf\x26\xf4\x23\x9a\x5a\x4d\xa6\xc3\x5e\x23\x08\xf6\x26\xa0\x90\x9c\x03\x65\xcf\x44\x72\x01\xb7\xdb\x43\x10\xac\x71\xaa\x0f\x40\xfd\xcf\x72\xf8\xee\xd0\xb7\x89\xca\xd6\xb4\xcb\x69\xda\x1a\x4c\xd4\xb5\xa0\x2e\x47\xd9\xce\xe8\x44\x51\x0b\xea\x72\xd4\xf0\xa1\x53\xbd\xd4\x82\x0a\x8d\xca\xb2\xf5\x1f\xa0\xf6\xbb\x3f\x1c\x2e\xad\xca\x57\x78\x2a\xb9\x7c\x29\x76\x8a\x8f\x87\xbf\xbe\xff\x9b\x19\xb2\xf5\x68\xff\x1e\xe0\x4c\xc0\x3c\xc3\x68\xfb\xe6\xfa\x8c\x7a\x1e\x5b\xb4\x7a\x73\x7d\x46\x3d\x4f\x2a\x5a\xbd\xb9\xde\x51\xff\x06\x00\x00\xff\xff\xf4\xd7\x34\xf2\x86\x05\x00\x00")

func migrationsSqlPostgres7SqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlPostgres7Sql,
		"migrations/sql/postgres/7.sql",
	)
}

func migrationsSqlPostgres7Sql() (*asset, error) {
	bytes, err := migrationsSqlPostgres7SqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/postgres/7.sql", size: 1414, mode: os.FileMode(438), modTime: time.Unix(1540974886, 0)}
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

	info := bindataFileInfo{name: "migrations/sql/tests/.gitkeep", size: 0, mode: os.FileMode(438), modTime: time.Unix(1540902707, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests1_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x93\x31\x6b\xc3\x30\x10\x85\xe7\x08\xf4\x1f\x6e\x93\x4d\xe5\x21\x5d\x3b\x15\x9a\x21\x50\x1c\x68\x92\x76\x14\x87\x74\xb1\x05\x8d\x94\xde\xc9\x94\x52\xfa\xdf\x4b\x6c\xd3\x76\xea\x9e\x0c\x82\x7b\xef\xc1\x7b\x7c\x83\x9a\x06\x6e\x8e\xb1\x63\x2c\x04\xfb\x93\x56\xeb\x76\xbb\x7a\xda\xc1\xba\xdd\x6d\xb4\x5a\xf4\x1f\x81\xd1\x65\x1c\x4a\x7f\xeb\xd0\x7b\x12\x81\x4a\x62\x97\xb0\x0c\x4c\x16\x98\xde\x06\x92\xe2\x62\xf8\xb9\x29\x38\x2c\x16\xfc\x6b\xa4\x34\x05\xe2\xf3\x89\x2c\x74\x8c\xe9\x9c\xce\xf2\x90\xf9\xe8\x02\x16\xb4\x20\x24\x12\x73\x1a\x55\xad\xd5\xf3\xfd\xe3\x7e\xb5\xd5\x6a\x51\x99\x65\x23\xb1\x33\x16\xcc\xb2\x99\xeb\x8d\x85\x76\xf3\x52\xd5\xa3\x37\x8d\x4c\xf9\x58\x3b\x9d\xf3\xd2\xaf\x75\x7e\x9f\x5f\xa6\xbe\xd3\xea\x5f\x42\xa6\x03\x93\xf4\xd7\x8c\xe8\x73\xa0\x6b\xe6\xcb\x31\xf8\x8b\xe7\xfb\xfb\x29\x1f\xf2\x7b\xd2\xea\x3b\x00\x00\xff\xff\x6b\x36\x62\xc8\xa7\x03\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/1_test.sql", size: 935, mode: os.FileMode(438), modTime: time.Unix(1540938906, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests2_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x93\xbf\x6a\xc3\x30\x10\xc6\xe7\x08\xf4\x0e\xb7\xc9\xa6\xf2\xe2\xb5\x53\xa1\x19\x02\xc5\x81\x26\x69\x47\x71\x95\x2e\xb6\x4a\x23\xa5\x3a\x99\x52\x4a\xdf\xbd\x24\x72\xff\x4c\x7d\x80\x78\x10\xdc\x7d\x9f\xd0\x8f\xdf\xa0\xa6\x81\xab\x83\xef\x13\x66\x82\xdd\x51\x8a\x55\xb7\x59\xde\x6f\x61\xd5\x6d\xd7\x52\x2c\x86\x77\x97\xd0\x44\x1c\xf3\xd0\x1a\xb4\x96\x98\xa1\x62\xdf\x07\xcc\x63\x22\x0d\x89\x5e\x47\xe2\x6c\xbc\xfb\x99\xc9\x19\xcc\x1a\xec\x8b\xa7\x50\x0a\xb6\xf1\x48\x1a\xfa\x84\xe1\xd4\x4e\xeb\x3e\xa6\x83\x71\x98\x51\x03\x13\xb3\x8f\xe1\x7b\x1b\x9f\x9e\xc9\xe6\x5a\x8a\x87\x9b\xbb\xdd\x72\x23\xc5\xa2\x52\x6d\xc3\xbe\x57\x1a\x54\xdb\x4c\x1c\xa5\xa1\x5b\x3f\x56\xf5\x39\x2b\xb4\xd2\x9f\xdf\x2f\xe3\x84\xfc\x8d\x4e\xe7\xe3\x73\xba\x57\x38\xaa\xbe\x96\xe2\x5f\xef\x44\xfb\x44\x3c\xcc\x4f\xdc\x46\x47\xf3\xb3\x8e\xde\xd9\x0b\xb5\xfe\xfb\xd9\x6f\xe3\x5b\x90\xe2\x2b\x00\x00\xff\xff\x6e\xd3\x5c\x28\xff\x03\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/2_test.sql", size: 1023, mode: os.FileMode(438), modTime: time.Unix(1540938945, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests3_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x94\x4f\x4b\x33\x31\x10\x87\xcf\x0d\xe4\x3b\xcc\x2d\xbb\xbc\xd9\xcb\xdb\xa3\x27\xc1\x1e\x0a\xb2\x05\xdb\xea\x31\x8c\xc9\x74\x37\x6a\x93\x35\x93\x45\x44\xfc\xee\xd2\x66\xfd\x73\xf2\x03\xb8\x87\xc0\xcc\xfc\x42\x1e\x1e\x02\xd3\x34\xf0\xef\xe8\xbb\x84\x99\x60\x3f\x48\xb1\x6e\xb7\xab\x9b\x1d\xac\xdb\xdd\x46\x8a\x45\xff\xea\x12\x9a\x88\x63\xee\xff\x1b\xb4\x96\x98\xa1\x62\xdf\x05\xcc\x63\x22\x0d\x89\x9e\x47\xe2\x6c\xbc\xfb\xaa\xc9\x19\xcc\x1a\xec\x93\xa7\x50\x02\xb6\x71\x20\x0d\x5d\xc2\x70\x4a\xa7\xf6\x10\xd3\xd1\x38\xcc\xa8\x81\x89\xd9\xc7\xf0\xd9\x8d\xf7\x0f\x64\x73\x2d\xc5\xed\xe5\xf5\x7e\xb5\x95\x62\x51\xa9\x65\xc3\xbe\x53\x1a\xd4\xb2\x99\x38\x4a\x43\xbb\xb9\xab\xea\xf3\xac\xd0\x4a\x7e\x7e\xbf\x94\x13\xf2\x7b\x74\x3a\x6f\xef\xd3\xbd\xc2\x51\xf5\x85\x14\xbf\x7a\x27\x3a\x24\xe2\x7e\x7e\xe2\x36\x3a\x9a\x9f\x75\xf4\xce\xce\xcf\x7a\x78\xb4\x7f\xf5\xaf\x7f\xae\xb8\xab\xf8\x12\xa4\xf8\x08\x00\x00\xff\xff\xcd\x58\x2d\xe2\xf5\x04\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/3_test.sql", size: 1269, mode: os.FileMode(438), modTime: time.Unix(1540938977, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests4_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x94\x4f\x4b\x33\x31\x10\xc6\xcf\x0d\xe4\x3b\xcc\x6d\x77\x79\xb3\x97\x97\xde\x3c\x09\xf6\x50\x90\x2d\xd8\x56\x8f\x61\x4c\xa6\xbb\x51\x9b\xac\x99\x44\x11\xf1\xbb\x4b\x9b\xf5\xcf\xc9\x2f\xb0\x87\xc0\xcc\xf3\x04\x7e\xf0\x23\xa4\x6d\xe1\xdf\xd1\xf5\x11\x13\xc1\x7e\x94\x62\xdd\x6d\x57\x37\x3b\x58\x77\xbb\x8d\x14\x8b\xe1\xcd\x46\xd4\x01\x73\x1a\xfe\x6b\x34\x86\x98\xa1\x66\xd7\x7b\x4c\x39\x92\x82\x48\xcf\x99\x38\x69\x67\xbf\x67\xb2\x1a\x93\x02\xf3\xe4\xc8\x97\x82\x4d\x18\x49\x41\x1f\xd1\x9f\xda\x69\x3d\x84\x78\xd4\x16\x13\x2a\x60\x62\x76\xc1\x7f\x6d\xf9\xfe\x81\x4c\x52\x80\x26\xb9\x17\x6a\xa4\xb8\xbd\xbc\xde\xaf\xb6\x52\x2c\xea\x6a\xd9\xb2\xeb\x2b\x05\xd5\xb2\x9d\x78\x95\x82\x6e\x73\x57\x37\xe7\xac\x50\x4b\x7f\xe6\x94\x71\x42\xff\x44\xa7\xf3\xfe\x31\xdd\x2b\xbc\x4a\x41\x8a\x99\x9a\x0b\x29\xfe\xb4\x10\xe9\x10\x89\x87\xb9\x6b\x30\xc1\xd2\xdc\x1d\x04\x67\xcd\xdc\x1d\x8c\x8f\x66\x1e\xef\xe0\xf7\x47\x79\x15\x5e\xbd\x14\x9f\x01\x00\x00\xff\xff\x76\xf4\x1f\xeb\x3b\x05\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/4_test.sql", size: 1339, mode: os.FileMode(438), modTime: time.Unix(1540939034, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests5_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x94\x4f\x4b\x33\x31\x10\xc6\xcf\x0d\xe4\x3b\xcc\x6d\x77\x79\xb3\x97\x17\x7a\xf2\x24\xd8\x43\x41\xb6\x60\x5b\x3d\x86\x31\x99\xee\x46\x6d\xb2\x66\x12\x45\xc4\xef\x2e\x6d\xd6\x3f\x27\xbf\xc0\x1e\x02\x33\xcf\x13\xf8\xc1\x8f\x90\xb6\x85\x7f\x47\xd7\x47\x4c\x04\xfb\x51\x8a\x75\xb7\x5d\xdd\xec\x60\xdd\xed\x36\x52\x2c\x86\x37\x1b\x51\x07\xcc\x69\xf8\xaf\xd1\x18\x62\x86\x9a\x5d\xef\x31\xe5\x48\x0a\x22\x3d\x67\xe2\xa4\x9d\xfd\x9e\xc9\x6a\x4c\x0a\xcc\x93\x23\x5f\x0a\x36\x61\x24\x05\x7d\x44\x7f\x6a\xa7\xf5\x10\xe2\x51\x5b\x4c\xa8\x80\x89\xd9\x05\xff\xb5\xe5\xfb\x07\x32\x49\x01\x9a\xe4\x5e\xa8\x91\xe2\xf6\xf2\x7a\xbf\xda\x4a\xb1\xa8\xab\x65\xcb\xae\xaf\x14\x54\xcb\x76\xe2\x55\x0a\xba\xcd\x5d\xdd\x9c\xb3\x42\x2d\xfd\x99\x53\xc6\x09\xfd\x13\x9d\xce\xfb\xc7\x74\xaf\xf0\x2a\x05\x29\x66\x6a\x2e\xa4\xf8\xd3\x42\xa4\x43\x24\x1e\xe6\xae\xc1\x04\x4b\x73\x77\x10\x9c\x35\x73\x77\x30\x3e\x9a\x79\xbc\x83\xdf\x1f\xe5\x55\x78\xf5\x52\x7c\x06\x00\x00\xff\xff\x41\x8d\x37\xbe\x3b\x05\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/5_test.sql", size: 1339, mode: os.FileMode(438), modTime: time.Unix(1540939077, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests6_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x94\x4f\x4b\x33\x31\x10\xc6\xcf\x0d\xe4\x3b\xcc\x6d\x77\x79\xb3\x97\xf7\xd0\x8b\x27\xc1\x1e\x0a\xb2\x05\xdb\xea\x31\x8c\xc9\x74\x37\x6a\x93\x35\x93\x28\x22\x7e\x77\x69\xb3\xfe\x39\xf9\x05\xf6\x10\x98\x79\x9e\xc0\x0f\x7e\x84\xb4\x2d\xfc\x3b\xba\x3e\x62\x22\xd8\x8f\x52\xac\xbb\xed\xea\x66\x07\xeb\x6e\xb7\x91\x62\x31\xbc\xd9\x88\x3a\x60\x4e\xc3\x7f\x8d\xc6\x10\x33\xd4\xec\x7a\x8f\x29\x47\x52\x10\xe9\x39\x13\x27\xed\xec\xf7\x4c\x56\x63\x52\x60\x9e\x1c\xf9\x52\xb0\x09\x23\x29\xe8\x23\xfa\x53\x3b\xad\x87\x10\x8f\xda\x62\x42\x05\x4c\xcc\x2e\xf8\xaf\x2d\xdf\x3f\x90\x49\x0a\xd0\x24\xf7\x42\x8d\x14\xb7\x97\xd7\xfb\xd5\x56\x8a\x45\x5d\x2d\x5b\x76\x7d\xa5\xa0\x5a\xb6\x13\xaf\x52\xd0\x6d\xee\xea\xe6\x9c\x15\x6a\xe9\xcf\x9c\x32\x4e\xe8\x9f\xe8\x74\xde\x3f\xa6\x7b\x85\x57\x29\x48\x31\x53\x73\x21\xc5\x9f\x16\x22\x1d\x22\xf1\x30\x77\x0d\x26\x58\x9a\xbb\x83\xe0\xac\x99\xbb\x83\xf1\xd1\xcc\xe3\x1d\xfc\xfe\x28\xaf\xc2\xab\x97\xe2\x33\x00\x00\xff\xff\x18\x06\x4f\x41\x3b\x05\x00\x00")

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

	info := bindataFileInfo{name: "migrations/sql/tests/6_test.sql", size: 1339, mode: os.FileMode(438), modTime: time.Unix(1540939091, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _migrationsSqlTests7_testSql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x95\x4f\x4b\xc3\x40\x10\xc5\xcf\x5d\xd8\xef\x30\xb7\x34\xb8\xb9\x78\xe9\xc1\x93\x60\x0f\x05\x49\xc1\xb6\x7a\x0c\xe3\xee\x34\x59\xb5\xbb\x71\xff\x28\x22\x7e\x77\x49\x37\xd5\xf4\x22\x08\x39\xe6\x10\x98\x79\x2f\xfb\x78\xfc\x58\xd8\xa2\x80\x8b\x83\xae\x1d\x06\x82\x5d\xcb\xd9\xaa\xdc\x2c\xef\xb6\xb0\x2a\xb7\x6b\xce\x66\xcd\x87\x72\x58\x59\x8c\xa1\xb9\xac\x50\x4a\xf2\x1e\xe6\x5e\xd7\x06\x43\x74\x24\xc0\xd1\x6b\x24\x1f\x2a\xad\x7e\x66\x52\x15\x06\x01\xf2\x45\x93\x49\x86\x97\xb6\x25\x01\xb5\x43\xd3\xb9\xfd\xba\xb7\xee\x50\x29\x0c\x28\xc0\x93\xf7\xda\x9a\xd3\x16\x1f\x9f\x48\x06\x01\x28\x83\x7e\xa3\xb3\xe0\xa8\x34\x19\x39\x08\x3b\x29\x39\x67\xf7\xd7\xb7\xbb\xe5\x86\xb3\xd9\x3c\x5b\x14\x5e\xd7\x99\x80\x6c\x51\xf4\x87\x33\x01\xe5\xfa\x61\x9e\x1f\xb5\xd4\x2d\xf9\xc7\x36\x69\xec\x33\x7f\xa5\xee\xfb\xfc\xea\xff\x4b\xad\x32\x01\xc1\x45\x1a\x46\x93\x2a\x30\xaa\xf3\x88\x4e\xc8\xaf\x38\xfb\x93\xa7\xa3\xbd\x23\xdf\x4c\x40\xc7\x02\x2a\xad\xa2\x89\xe6\x58\x34\xad\x56\x72\xa2\x39\x16\xcd\xf6\x59\x4e\x77\xf3\x7f\x34\x87\x4f\xd3\x8d\x7d\x37\x9c\x7d\x07\x00\x00\xff\xff\x0c\xe1\x1b\x35\xad\x06\x00\x00")

func migrationsSqlTests7_testSqlBytes() ([]byte, error) {
	return bindataRead(
		_migrationsSqlTests7_testSql,
		"migrations/sql/tests/7_test.sql",
	)
}

func migrationsSqlTests7_testSql() (*asset, error) {
	bytes, err := migrationsSqlTests7_testSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "migrations/sql/tests/7_test.sql", size: 1709, mode: os.FileMode(438), modTime: time.Unix(1540939928, 0)}
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
	"migrations/sql/shared/1.sql":      migrationsSqlShared1Sql,
	"migrations/sql/shared/2.sql":      migrationsSqlShared2Sql,
	"migrations/sql/shared/3.sql":      migrationsSqlShared3Sql,
	"migrations/sql/shared/4.sql":      migrationsSqlShared4Sql,
	"migrations/sql/mysql/.gitkeep":    migrationsSqlMysqlGitkeep,
	"migrations/sql/mysql/5.sql":       migrationsSqlMysql5Sql,
	"migrations/sql/mysql/6.sql":       migrationsSqlMysql6Sql,
	"migrations/sql/mysql/7.sql":       migrationsSqlMysql7Sql,
	"migrations/sql/postgres/.gitkeep": migrationsSqlPostgresGitkeep,
	"migrations/sql/postgres/5.sql":    migrationsSqlPostgres5Sql,
	"migrations/sql/postgres/6.sql":    migrationsSqlPostgres6Sql,
	"migrations/sql/postgres/7.sql":    migrationsSqlPostgres7Sql,
	"migrations/sql/tests/.gitkeep":    migrationsSqlTestsGitkeep,
	"migrations/sql/tests/1_test.sql":  migrationsSqlTests1_testSql,
	"migrations/sql/tests/2_test.sql":  migrationsSqlTests2_testSql,
	"migrations/sql/tests/3_test.sql":  migrationsSqlTests3_testSql,
	"migrations/sql/tests/4_test.sql":  migrationsSqlTests4_testSql,
	"migrations/sql/tests/5_test.sql":  migrationsSqlTests5_testSql,
	"migrations/sql/tests/6_test.sql":  migrationsSqlTests6_testSql,
	"migrations/sql/tests/7_test.sql":  migrationsSqlTests7_testSql,
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
				"5.sql":    &bintree{migrationsSqlMysql5Sql, map[string]*bintree{}},
				"6.sql":    &bintree{migrationsSqlMysql6Sql, map[string]*bintree{}},
				"7.sql":    &bintree{migrationsSqlMysql7Sql, map[string]*bintree{}},
			}},
			"postgres": &bintree{nil, map[string]*bintree{
				".gitkeep": &bintree{migrationsSqlPostgresGitkeep, map[string]*bintree{}},
				"5.sql":    &bintree{migrationsSqlPostgres5Sql, map[string]*bintree{}},
				"6.sql":    &bintree{migrationsSqlPostgres6Sql, map[string]*bintree{}},
				"7.sql":    &bintree{migrationsSqlPostgres7Sql, map[string]*bintree{}},
			}},
			"shared": &bintree{nil, map[string]*bintree{
				"1.sql": &bintree{migrationsSqlShared1Sql, map[string]*bintree{}},
				"2.sql": &bintree{migrationsSqlShared2Sql, map[string]*bintree{}},
				"3.sql": &bintree{migrationsSqlShared3Sql, map[string]*bintree{}},
				"4.sql": &bintree{migrationsSqlShared4Sql, map[string]*bintree{}},
			}},
			"tests": &bintree{nil, map[string]*bintree{
				".gitkeep":   &bintree{migrationsSqlTestsGitkeep, map[string]*bintree{}},
				"1_test.sql": &bintree{migrationsSqlTests1_testSql, map[string]*bintree{}},
				"2_test.sql": &bintree{migrationsSqlTests2_testSql, map[string]*bintree{}},
				"3_test.sql": &bintree{migrationsSqlTests3_testSql, map[string]*bintree{}},
				"4_test.sql": &bintree{migrationsSqlTests4_testSql, map[string]*bintree{}},
				"5_test.sql": &bintree{migrationsSqlTests5_testSql, map[string]*bintree{}},
				"6_test.sql": &bintree{migrationsSqlTests6_testSql, map[string]*bintree{}},
				"7_test.sql": &bintree{migrationsSqlTests7_testSql, map[string]*bintree{}},
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
