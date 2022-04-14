package common

import (
	"strconv"

	"gopkg.in/yaml.v2"
)

// transfer struct to map[string]string
func StructToMap(s interface{}) (map[string]string, error) {
	m := make(map[interface{}]interface{})
	bytes, err := yaml.Marshal(s)
	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(bytes, &m)

	args := make(map[string]string)
	viewMap(m, args)
	return args, nil
}

// delete item with value=
func CleanMap(m map[string]string) {
	for k, v := range m {
		if len(v) < 1 {
			delete(m, k)
		}
	}
}

// transfer map[interface{}]interface{} to map[string]string
func viewMap(s map[interface{}]interface{}, m map[string]string) {
	for k, v := range s {
		// v is map
		if _, ok := v.(map[interface{}]interface{}); ok {
			viewMap(v.(map[interface{}]interface{}), m)
		} else {

			m[Strval(k)] = Strval(v)
		}
	}
}

func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := yaml.Marshal(value)
		key = string(newValue)
	}

	return key
}
