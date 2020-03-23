package cidr

import (
	"strconv"
	"strings"
)

func ParseArrayToRandomValue(ips []string, rdValues []int) []string {
	var IpsResult []string
	for i := range ips {
		ipSegTail := strings.Split(ips[i], ".")[3]
		for j := range rdValues {
			if ipSegTail == strconv.Itoa(rdValues[j]) {
				IpsResult = append(IpsResult, ips[i])
			}
		}
	}
	return IpsResult
}

// 去重
func RemoveRepeatedElement(arr []string) []string {
	var newArr []string
	for i := range arr {
		repeat := false
		for j := range newArr {
			if arr[i] == newArr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}
