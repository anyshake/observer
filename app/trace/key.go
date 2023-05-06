package trace

func HasKey(m map[string]interface{}, k []string) bool {
	for _, v := range k {
		if _, ok := m[v]; !ok {
			return false
		}
	}

	return true
}
