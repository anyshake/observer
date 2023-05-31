package trace

func IsEmpty(m map[string]any, k []string) bool {
	for _, v := range k {
		if len(m[v].(string)) == 0 {
			return false
		}
	}

	return true
}
