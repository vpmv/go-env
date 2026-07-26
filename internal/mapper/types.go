package mapper

func isMap(m any) bool {
	_, ok := m.(map[string]any)
	return ok
}

func isSlice(m any) bool {
	_, ok := m.([]any)
	return ok
}

func isMapSlice(m any) bool {
	if isSlice(m) && len(m.([]any)) > 0 {
		return isMap(m.([]any)[0])
	}
	return false
}
