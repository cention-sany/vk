package vk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
)

func ownerOptions(vals url.Values, v ...interface{}) (isGroup bool) {
	if v != nil {
		vSize := len(v)
		if vSize > 0 {
			if id, ok := v[0].(string); ok && id != "" {
				if vSize > 1 {
					if isCommunity, ok := v[1].(bool); ok {
						if isCommunity {
							id = fmt.Sprint("-", id)
							isGroup = true
						}
					}
				}
				vals.Set("owner_id", id)
			} else if id, ok := v[0].(int); ok && id > 0 {
				if vSize > 1 {
					if isCommunity, ok := v[1].(bool); ok {
						if isCommunity {
							id = -1 * id
							isGroup = true
						}
					}
				}
				vals.Set("owner_id", strconv.Itoa(id))
			}
		}
	}
	return
}

func unmarshaler(v interface{}, r io.Reader) error {
	unmarshal := json.NewDecoder(r)
	if err := unmarshal.Decode(v); err != nil {
		return err
	}
	return nil
}

// ElemInSlice checks if element is in the slice
func ElemInSlice(elem string, slice []string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}
