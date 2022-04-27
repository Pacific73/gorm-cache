package testkit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrimaryCacheFunctionality(t *testing.T) {
	Convey("test primary cache functionality", t, func() {
		testFirst(primaryCache, primaryDB)

		testFind(primaryCache, primaryDB)

		testPrimaryFind(primaryCache, primaryDB)

		testPrimaryUpdate(primaryCache, primaryDB)

		testPrimaryDelete(primaryCache, primaryDB)
	})
}

func TestSearchCacheFunctionality(t *testing.T) {
	Convey("test search cache functionality", t, func() {
		testFirst(searchCache, searchDB)

		testFind(searchCache, searchDB)

		testSearchFind(searchCache, searchDB)

		testSearchCreate(searchCache, searchDB)

		testSearchDelete(searchCache, searchDB)

		testSearchUpdate(searchCache, searchDB)
	})
}
