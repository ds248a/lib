package strings

import (
	"github.com/ds248a/lib/strconv"
)

func Copy(s string) string {
	return string(strconv.S2B(s))
}

func CopySlice(slice []string) []string {
	dst := make([]string, len(slice))
	copy(dst, slice)
	return dst
}

func IndexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}

func Include(slice []string, s string) bool {
	return IndexOf(slice, s) != -1
}

func UniqueAppend(slice []string, s ...string) []string {
	for i := range s {
		if IndexOf(slice, s[i]) != -1 {
			continue
		}
		slice = append(slice, s[i])
	}
	return slice
}

func EqualSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func ReverseSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}
