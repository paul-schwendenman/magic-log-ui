package shared

import "encoding/json"

func MustJson(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
