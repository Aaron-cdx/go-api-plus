package jsonutils

import "encoding/json"

// JsonEncode used for marshal object to json string
func JsonEncode(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes) + "\n", nil
}
