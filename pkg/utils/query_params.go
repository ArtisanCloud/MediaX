package utils

import (
	"strings"

	"github.com/ArtisanCloud/MediaXCore/utils/object"
)

// StructToQueryParams converts a struct into a StringMap suitable for query parameters.
// It trims json tag options (e.g. omitempty) and drops empty string values to avoid
// emitting invalid parameters such as "id=" or "chart,omitempty".
func StructToQueryParams(payload interface{}) (*object.StringMap, error) {
	params, err := object.StructToStringMap(payload)
	if err != nil {
		return nil, err
	}
	cleaned := make(object.StringMap, len(*params))
	for rawKey, val := range *params {
		key := strings.TrimSpace(rawKey)
		if key == "" || key == "-" {
			continue
		}
		if idx := strings.Index(key, ","); idx >= 0 {
			key = strings.TrimSpace(key[:idx])
		}
		if key == "" {
			continue
		}
		if strings.TrimSpace(val) == "" {
			continue
		}
		cleaned[key] = val
	}
	return &cleaned, nil
}
