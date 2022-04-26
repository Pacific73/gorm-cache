package testkit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPrimaryCacheFunctionality(t *testing.T) {
	Convey("test primary cache functionality", t, func() {
		//testFirst(primaryCache, primaryDB)
		testFind(allCache, allDB)
	})

}
