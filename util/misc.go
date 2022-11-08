package util

import "math/rand"

func ShouldCache(tableName string, tables []string) bool {
	if len(tables) == 0 {
		return true
	}
	return ContainString(tableName, tables)
}

func ContainString(target string, slice []string) bool {
	for _, s := range slice {
		if target == s {
			return true
		}
	}
	return false
}

func RandFloatingInt64(v int64) int64 {
	randNum := rand.Float64()*0.2 + 0.9
	return int64(float64(v) * randNum)
}
