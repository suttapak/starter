package helpers

import "encoding/json"

func ParseJson[T any](src any, dns *T) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, dns); err != nil {
		return err
	}
	return nil
}
