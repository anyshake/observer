package trace

func IsEmpty(m map[string]interface{}, k []string) bool {
	for _, v := range k {
		if len(m[v].(string)) == 0 {
			return false
		}
	}

	return true
}
