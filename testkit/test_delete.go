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
	result := db.Where("id IN (?)", []int{1, 2, 3, 4, 5}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2, 3}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(len(models), ShouldEqual, 3)

	result = db.Delete(&TestModel{ID: 5})
	So(result.Error, ShouldBeNil)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2, 3, 4}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 2)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2, 3, 4, 5}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 2)

	result = db.Delete([]*TestModel{{ID: 3}, {ID: 4}})
	So(result.Error, ShouldBeNil)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 3)

	result = db.Where("id = 2").Delete(&TestModel{})
	So(result.Error, ShouldBeNil)

	result = db.Where("id IN (?)", []int{1}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 4)
}
