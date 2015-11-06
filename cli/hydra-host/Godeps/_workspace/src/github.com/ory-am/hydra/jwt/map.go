package jwt

import "encoding/json"

type Map map[string]interface{}

func (m Map) Marshall() ([]byte, error) {
	return json.Marshal(m)
}

func Unmarshal(data []byte) (Map, error) {
	var m Map
	return m, json.Unmarshal(data, &m)
}
