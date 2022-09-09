package testkit

import (
	"github.com/Pacific73/gorm-cache/cache"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
)

func testFirst(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	var model = new(TestModel)

	result := db.Where("id = ?", 1).First(model)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)
	So(model.ID, ShouldEqual, 1)

	model = new(TestModel)
	result = db.Where("id = ?", 1).First(model)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)

	targetModel := &TestModel{
		ID: 2,
	}

	result = db.First(targetModel)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)

	result = db.First(targetModel)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 2)
}

func testFind(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models := make([]*TestModel, 0)
	result := db.Where("id IN (?)", []int{1, 2}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(len(models), ShouldEqual, 2)
	So(models[0].Value1, ShouldEqual, 1)
	So(models[1].Value1, ShouldEqual, 2)
}

func testPtrFind(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	_prtValue := int64(1)
	model := &TestModel{
		PtrValue1: &_prtValue,
	}
	result := db.Model(model).Find(model)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(model.Value1, ShouldEqual, 1)
}

func testPrimaryFind(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models := make([]*TestModel, 0)
	result := db.Where("id IN (?)", []int{1, 2, 3}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int{1, 2}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(len(models), ShouldEqual, 2)
	So(models[0].Value1, ShouldEqual, 1)
	So(models[1].Value1, ShouldEqual, 2)

	models = make([]*TestModel, 0)
	result = db.Where("id < ?", 3).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)

	models = make([]*TestModel, 0)
	result = db.Where("id IN (?)", []int64{1, 4}).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
}

func testSearchFind(cache *cache.Gorm2Cache, db *gorm.DB) {
	err := cache.ResetCache()
	So(err, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models := make([]*TestModel, 0)
	result := db.Where("id >= ?", 1).Where("id <= ?", 10).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 0)

	models = make([]*TestModel, 0)
	result = db.Where("id >= ?", 1).Where("id <= ?", 10).Find(&models)
	So(result.Error, ShouldBeNil)
	So(cache.GetHitCount(), ShouldEqual, 1)
	So(len(models), ShouldEqual, 10)
}
