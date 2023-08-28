package trace

func isEmpty(m map[string]any, k []string) bool {
	for _, v := range k {
		switch m[v].(type) {
		case string:
			if len(m[v].(string)) == 0 {
				return false
			}
		default:
			continue
		}
	}

	return true
}
