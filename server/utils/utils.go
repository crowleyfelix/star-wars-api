package utils

import "encoding/json"

func ToJSON(data interface{}) string {
	blob, _ := json.MarshalIndent(data, "", "	")
	return string(blob)
}
