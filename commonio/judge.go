package commonio

import (
	"encoding/json"
	"strconv"
)

func IsInt(str string) bool {
	if _, err := strconv.Atoi(str); err != nil {
		return false
	}
	return true
}

func IsJson(bs []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(bs, &js) == nil
}
