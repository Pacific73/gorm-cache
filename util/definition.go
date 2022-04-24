package util

import "errors"

var PrimaryCacheHit = errors.New("primary cache hit")
var SearchCacheHit = errors.New("search cache hit")

var ErrCacheUnmarshal = errors.New("cache hit, but unmarshal error")

type Kv struct {
	Key   string
	Value string
}
