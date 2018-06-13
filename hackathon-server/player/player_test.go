package player

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewPlayer(t *testing.T) {
	Convey("Testing with only one player", t, func() {
		p1 := NewPlayer("test1")
		So(p1.ID, ShouldBeGreaterThanOrEqualTo, 0)
		So(p1.Name, ShouldEqual, "test1")
	})

	Convey("Testing with two players", t, func() {
		p1 := NewPlayer("test1")
		p2 := NewPlayer("test2")
		So(p1.ID, ShouldNotEqual, p2.ID)
	})
}

func TestEquals(t *testing.T) {
	p1 := &Player{
		Name: "test1",
		ID:   1,
	}

	p2 := &Player{
		Name: "test1",
		ID:   1,
	}

	p3 := &Player{
		Name: "test2",
		ID:   2,
	}

	Convey("Testing with two similar player", t, func() {
		So(p1.Equals(p2), ShouldBeTrue)
		So(p2.Equals(p1), ShouldBeTrue)
	})

	Convey("Testing with two different player", t, func() {
		So(p1.Equals(p3), ShouldBeFalse)
	})
}

func TestGaia(t *testing.T) {
	Convey("Testing with Gaia", t, func() {
		So(GAIA.IsGaia(), ShouldBeTrue)
	})

	Convey("Testing with a normal player", t, func() {
		So(NewPlayer("test").IsGaia(), ShouldBeFalse)
	})
}
