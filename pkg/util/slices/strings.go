package slices

import (
	"regexp"
	"strings"
)

func ContainsStr(arr []string, str string) bool {
	return SearchStr(arr, str) >= 0
}

func DistinctStr(arr []string) []string {
	distinctSlice := make([]string, 0, len(arr))
	m := make(map[string]struct{}, len(arr))
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
	for i := range arr {
		arr[i] = f(arr[i])
	}
	return arr
}

func RemoveEmptyStr(arr []string) []string {
	for i := 0; i < len(arr); {
		if arr[i] == "" {
			arr = append(arr[:i], arr[i+1:]...)
		} else {
			i++
		}
	}
	return arr
}

func RemovePunctuation(str string) string {
	reg := regexp.MustCompile(`\p{P}+`)
	return reg.ReplaceAllString(str, "")
}

func RemoveSpecifiedStr(arr []string, remove string) []string {
	if i := SearchStr(arr, remove); i >= 0 {
		return append(arr[:i], arr[i+1:]...)
	}
	return arr
}

func SearchStr(arr []string, str string) int {
	if len(arr) == 0 {
		return -1
	}
	result := -1
	for index, v := range arr {
		if strings.Compare(v, str) == 0 {
			result = index
			break
		}
	}
	return result
}

//DifferenceStr return the difference set of `a-b`
func DifferenceStr(a []string, b []string) []string {
	ds := make([]string, 0, len(a))
	m := make(map[string]struct{}, len(b))
	for i := range b {
		m[b[i]] = struct{}{}
	}
	for i := range a {
		if _, ok := m[a[i]]; !ok {
			ds = append(ds, a[i])
		}
	}
	return ds
}
