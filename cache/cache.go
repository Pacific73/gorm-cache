package cache

import (
	"gorm.io/gorm"
)

type Gorm2Cache struct {
	db *gorm.DB
}

func (c *Gorm2Cache) AttachToDB() {

}
