package vk

import (
	"strings"
	"fmt"
)

// ElemInSlice checks if element is in the slice
func ElemInSlice(elem string, slice []string) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func JoinIntArr(arr []int) string {
	var str = make([]string,len(arr))
	for i, v := range arr {
		str[i] = fmt.Sprint(v)
	}
	return strings.Join(str,",")
}
