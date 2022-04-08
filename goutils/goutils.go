/**
* @author mnunez
 */

package goutils

import (
	"encoding/json"
	"strings"
)

func ToJSONString(value interface{}) (string, error) {
	bytes, error := json.Marshal(value)

	return string(bytes), error
}

func ToJSON(value string) (interface{}, error) {
	var jsonResult interface{}

	decoder := json.NewDecoder(strings.NewReader(value))
	decoder.UseNumber()

	if error := decoder.Decode(&jsonResult); error != nil {
		return nil, error
	} else {
		return jsonResult, nil
	}
}

func FromJSONTo(value string, instance interface{}) error {
	return json.Unmarshal([]byte(value), instance)
}
