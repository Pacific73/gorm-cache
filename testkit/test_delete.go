package testkit

import (
	"github.com/Pacific73/gorm-cache/cache"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func testPrimaryDelete(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models := make([]*TestModel, 0)
	result := db.Where("id IN (?)", []int{101, 102, 103, 104, 105}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{101, 102, 103}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(len(models), ShouldEqual, 3)

	result = db.Delete(&TestModel{ID: 105})
	So(result.Error, ShouldBeNil)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{101, 102, 103, 104}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 2)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{101, 102, 103, 104, 105}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 2)

	result = db.Delete([]*TestModel{{ID: 103}, {ID: 104}})
	So(result.Error, ShouldBeNil)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{101, 102}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 3)

	result = db.Where("id = 102").Delete(&TestModel{})
	So(result.Error, ShouldBeNil)

	result = db.Where("id IN (?)", []int{101}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 4)
}
