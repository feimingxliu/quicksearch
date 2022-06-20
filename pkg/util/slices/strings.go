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

//FilterStr applies f to each element in arr.
func FilterStr(arr []string, f func(string) string) []string {
	fArr := make([]string, 0)
	for i := range arr {
		fArr = append(fArr, f(arr[i]))
	}
	return fArr
}

func RemoveEmptyStr(arr []string) []string {
	a := make([]string, 0)
	for i := range arr {
		if len(arr[i]) > 0 {
			a = append(a, arr[i])
		}
	}
	return a
}
