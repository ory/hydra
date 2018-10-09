package mapx

import (
	"errors"
	"time"
)

var ErrKeyDoesNotExist = errors.New("key is not present in map")
var ErrKeyCanNotBeTypeAsserted = errors.New("key could not be type asserted")

func GetString(values map[interface{}]interface{}, key interface{}) (string, error) {
	if v, ok := values[key]; !ok {
		return "", ErrKeyDoesNotExist
	} else if sv, ok := v.(string); !ok {
		return "", ErrKeyCanNotBeTypeAsserted
	} else {
		return sv, nil
	}
}

func GetStringSlice(values map[interface{}]interface{}, key interface{}) ([]string, error) {
	if v, ok := values[key]; !ok {
		return []string{}, ErrKeyDoesNotExist
	} else if sv, ok := v.([]string); ok {
		return sv, nil
	} else if sv, ok := v.([]interface{}); ok {
		vs := make([]string, len(sv))
		for k, v := range sv {
			if vv, ok := v.(string); !ok {
				return []string{}, ErrKeyCanNotBeTypeAsserted
			} else {
				vs[k] = vv
			}
		}
		return vs, nil
	} else {
		return []string{}, ErrKeyCanNotBeTypeAsserted
	}
}

func GetTime(values map[interface{}]interface{}, key interface{}) (time.Time, error) {
	v, ok := values[key]
	if !ok {
		return time.Time{}, ErrKeyDoesNotExist
	}

	if sv, ok := v.(time.Time); ok {
		return sv, nil
	} else if sv, ok := v.(int64); ok {
		return time.Unix(sv, 0), nil
	} else if sv, ok := v.(int32); ok {
		return time.Unix(int64(sv), 0), nil
	} else if sv, ok := v.(int); ok {
		return time.Unix(int64(sv), 0), nil
	} else if sv, ok := v.(float64); ok {
		return time.Unix(int64(sv), 0), nil
	} else if sv, ok := v.(float32); ok {
		return time.Unix(int64(sv), 0), nil
	}

	return time.Time{}, ErrKeyCanNotBeTypeAsserted
}

func GetInt64(values map[interface{}]interface{}, key interface{}) (int64, error) {
	if v, ok := values[key]; !ok {
		return 0, ErrKeyDoesNotExist
	} else if sv, ok := v.(int64); !ok {
		return 0, ErrKeyCanNotBeTypeAsserted
	} else {
		return sv, nil
	}
}

func GetInt32(values map[interface{}]interface{}, key interface{}) (int32, error) {
	if v, ok := values[key]; !ok {
		return 0, ErrKeyDoesNotExist
	} else if sv, ok := v.(int32); ok {
		return sv, nil
	} else if sv, ok := v.(int); ok {
		return int32(sv), nil
	} else {
		return 0, ErrKeyCanNotBeTypeAsserted
	}
}

func GetInt(values map[interface{}]interface{}, key interface{}) (int, error) {
	if v, ok := values[key]; !ok {
		return 0, ErrKeyDoesNotExist
	} else if sv, ok := v.(int32); ok {
		return int(sv), nil
	} else if sv, ok := v.(int); ok {
		return sv, nil
	} else {
		return 0, ErrKeyCanNotBeTypeAsserted
	}
}

func GetFloat32(values map[interface{}]interface{}, key interface{}) (float32, error) {
	if v, ok := values[key]; !ok {
		return 0, ErrKeyDoesNotExist
	} else if sv, ok := v.(float32); !ok {
		return 0, ErrKeyCanNotBeTypeAsserted
	} else {
		return sv, nil
	}
}

func GetFloat64(values map[interface{}]interface{}, key interface{}) (float64, error) {
	if v, ok := values[key]; !ok {
		return 0, ErrKeyDoesNotExist
	} else if sv, ok := v.(float64); !ok {
		return 0, ErrKeyCanNotBeTypeAsserted
	} else {
		return sv, nil
	}
}

func GetStringDefault(values map[interface{}]interface{}, key interface{}, defaultValue string) string {
	if s, err := GetString(values, key); err == nil {
		return s
	}
	return defaultValue
}

func GetStringSliceDefault(values map[interface{}]interface{}, key interface{}, defaultValue []string) []string {
	if s, err := GetStringSlice(values, key); err == nil {
		return s
	}
	return defaultValue
}
