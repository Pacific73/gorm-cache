package cache

import (
	"fmt"
	"strings"
)

func GenPrimaryCacheKey(tableName string, primaryKey string) string {
	return fmt.Sprintf("p:%s:%s", tableName, primaryKey)
}

func GenPrimaryCachePrefix(tableName string) string {
	return "p:" + tableName
}

func GenSearchCacheKey(tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("s:%s:%s", tableName, buf.String())
}

func GenSearchCachePrefix(tableName string) string {
	return "s:" + tableName
}
