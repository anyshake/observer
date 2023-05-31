package trace

func HasKey(m map[string]any, k []string) bool {
	for _, v := range k {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}
