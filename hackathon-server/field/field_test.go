package field

import (
	"testing"

	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPoint(t *testing.T) {
	Convey("Testing a new point", t, func() {
		point := NewPoint()
		So(point.Population, ShouldEqual, 0)
		So(point.Owner.IsGaia(), ShouldBeTrue)
	})
}

func TestNewField(t *testing.T) {
	Convey("Testing basic fields", t, func() {
		field := NewField(2, 3)
		So(field.Field, ShouldNotBeNil)
		So(field.SizeX, ShouldEqual, 2)
		So(field.SizeY, ShouldEqual, 3)
	})

	Convey("Testing field array", t, func() {
		field := NewField(2, 3)
		So(len(field.Field), ShouldEqual, 2)
		for i := 0; i < 2; i++ {
			So(field.Field[i], ShouldNotBeNil)
			So(len(field.Field[i]), ShouldEqual, 3)
			for j := 0; j < 3; j++ {
				So(field.Field[i][j], ShouldNotBeNil)
				So(field.Field[i][j].Owner.IsGaia(), ShouldBeTrue)
			}
		}
	})
}

func TestPlacePlayer(t *testing.T) {
	Convey("Testing with a single cell", t, func() {
		f := NewField(1, 1)
		p1 := player.NewPlayer("test1")
		f.PlacePlayer(p1, 8)
		So(f.Field[0][0].Owner.Equals(p1), ShouldBeTrue)
		So(f.Field[0][0].Population, ShouldEqual, 8)
	})

	Convey("With two cells", t, func() {
		f := NewField(2, 1)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		f.PlacePlayer(p1, 8)
		f.PlacePlayer(p2, 4)

		So(f.Field[0][0].Owner.IsGaia(), ShouldBeFalse)
		So(f.Field[1][0].Owner.IsGaia(), ShouldBeFalse)

		if f.Field[0][0].Owner.Equals(p1) {
			So(f.Field[0][0].Population, ShouldEqual, 8)
			So(f.Field[1][0].Population, ShouldEqual, 4)
		} else {
			So(f.Field[0][0].Population, ShouldEqual, 4)
			So(f.Field[1][0].Population, ShouldEqual, 8)
		}
	})
}
