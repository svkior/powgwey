package main_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Given some integer is decremented", t, func() {
		x := 1

		Convey("When the integer is decremented", func() {
			x--

			Convey("The value should be lower by one", func() {
				So(x, ShouldEqual, 0)
			})
		})
	})
}
