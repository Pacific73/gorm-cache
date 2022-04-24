package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GenInstanceId() string {
	charList := []byte("1234567890")
	rand.Seed(time.Now().Unix())
	length := 10
	str := make([]byte, 0)
	for i := 0; i < length; i++ {
		str = append(str, charList[rand.Intn(len(charList))])
	}
	return string(str)
}

func GenPrimaryCacheKey(instanceId string, tableName string, primaryKey string) string {
	return fmt.Sprintf("p:%s:%s:%s", instanceId, tableName, primaryKey)
}

func GenPrimaryCachePrefix(instanceId string, tableName string) string {
	return "p:" + instanceId + ":" + tableName
}

func GenSearchCacheKey(instanceId string, tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		buf.WriteString(fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("s:%s:%s:%s", instanceId, tableName, buf.String())
}

func GenSearchCachePrefix(instanceId string, tableName string) string {
	return "s:" + instanceId + ":" + tableName
}
