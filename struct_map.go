package openai

import (
	"bytes"
	"encoding/json"
)

func mapStruct(v any) (map[string]any, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	output := make(map[string]any)
	decoder := json.NewDecoder(buf)
	err = decoder.Decode(&output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
