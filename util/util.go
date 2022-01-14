package util

import "fmt"

func FormatRedisKey(prefixKey string, targetSuffix interface{}) string {
	return fmt.Sprintf("%s%s", prefixKey, targetSuffix)
}

func Includes(len int, f func(int) bool) bool {
	for i := 0; i < len; i++ {
		if f(i) {
			return true
		}
	}
	return false
}
