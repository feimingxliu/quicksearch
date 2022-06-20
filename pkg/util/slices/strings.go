package slices

func ContainsStr(arr []string, str string) bool {
	for i := range arr {
		if arr[i] == str {
			return true
		}
	}
	return false
}

func DistinctStr(arr []string) []string {
	distinctSlice := make([]string, 0)
	m := make(map[string]struct{})
	for i := range arr {
		if _, ok := m[arr[i]]; !ok {
			m[arr[i]] = struct{}{}
			distinctSlice = append(distinctSlice, arr[i])
		}
	}
	return distinctSlice
}
