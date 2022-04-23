package callback

import (
	"github.com/Pacific73/gorm-cache/cache"
	"gorm.io/gorm"
)

func AfterRow(cache *cache.Gorm2Cache) func (db *gorm.DB) {
	return func (db *gorm.DB) {
		return
	}
}