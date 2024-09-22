package utils

func AddItemToMap(m map[string]interface{}, key string, value any) {
	if value != nil {
		switch v := value.(type) {
		case int:
			if v != 0 {
				m[key] = value
			}
		case uint32:
			if v != 0 {
				m[key] = value
			}
		case uint64:
			if v != 0 {
				m[key] = value
			}
		case bool:
			if v {
				m[key] = value
			}
		case string:
			if v != "" {
				m[key] = value
			}
		case float64:
			if v != 0 {
				m[key] = value
			}
		case []any:
			if len(v) > 0 {
				m[key] = value
			}
		default:
			if v != nil {
				m[key] = value
			}
		}
	}
}
