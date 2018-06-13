package game

import (
	"testing"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"

	. "github.com/smartystreets/goconvey/convey"
)

// func Attack(f *field.Field, attacker *player.Player, fromX int, fromY int, toX int, toY int) error {

func TestAttack(t *testing.T) {
	Convey("Testing invalid parameters", t, func() {
		f := field.NewField(3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		Convey("Testing with invalid starting position", func() {
			_, err := Attack(f, nil, -1, 0, 0, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 3, 0, 0, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, -1, 0, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 4, 0, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Testing with invalid ending position", func() {

			_, err := Attack(f, nil, 0, 0, -1, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 0, 3, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 0, 0, -1)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 0, 0, 4)
			So(err, ShouldNotBeNil)
		})

		Convey("Testing when attacking from too far", func() {
			_, err := Attack(f, nil, 0, 0, 2, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 0, 0, 2)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 2, 0, 0, 0)
			So(err, ShouldNotBeNil)
			_, err = Attack(f, nil, 0, 2, 0, 0)
			So(err, ShouldNotBeNil)
		})

		Convey("Testing users and populations", func() {
			f2 := field.NewField(2, 2)
			f2.Field[0][0].Owner = p1
			f2.Field[0][0].Population = 3
			f2.Field[0][1].Owner = p1
			f2.Field[0][1].Population = 3
			f2.Field[1][1].Owner = p2
			f2.Field[1][1].Population = 1

			// Not the owner
			_, err := Attack(f2, p2, 0, 0, 1, 1)
			So(err, ShouldNotBeNil)

			// Attack itself
			_, err = Attack(f2, p1, 0, 0, 0, 1)
			So(err, ShouldNotBeNil)

			// Not enough people
			_, err = Attack(f2, p2, 1, 1, 0, 1)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Testing with correct values", t, func() {
		f := field.NewField(3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		f.Field[0][0].Owner = p1
		f.Field[0][0].Population = 5
		f.Field[1][1].Owner = p2
		f.Field[1][1].Population = 200

		Convey("When attacking gaia", func() {
			_, err := Attack(f, p1, 0, 0, 0, 1)
			So(err, ShouldBeNil)
			So(f.Field[0][1].Owner.Equals(p1), ShouldBeTrue)
			So(f.Field[0][0].Population, ShouldEqual, 1)
			So(f.Field[0][1].Population, ShouldEqual, 4)
		})

		Convey("When attacking other people", func() {
			_, err := Attack(f, p2, 1, 1, 0, 0)
			So(err, ShouldBeNil)

			// Probability based (200 over 5)
			So(f.Field[0][0].Owner.Equals(p2), ShouldBeTrue)
			So(f.Field[0][0].Population, ShouldBeLessThanOrEqualTo, 199)
			So(f.Field[1][1].Population, ShouldEqual, 1)
		})
	})
}

func TestHasLost(t *testing.T) {
	f := field.NewField(3, 4)
	p1 := player.NewPlayer("test1")
	p2 := player.NewPlayer("test2")
	f.Field[0][0].Owner = p1
	f.Field[0][0].Population = 5

	Convey("Testing has lost", t, func() {
		So(HasLost(f, p1), ShouldBeFalse)
		So(HasLost(f, p2), ShouldBeTrue)
	})
}

func TestHasWin(t *testing.T) {
	f := field.NewField(3, 4)
	p1 := player.NewPlayer("test1")
	p2 := player.NewPlayer("test2")
	f.Field[0][0].Owner = p1
	f.Field[0][0].Population = 5

	Convey("Testing has win", t, func() {
		So(HasWin(f, p1), ShouldBeTrue)
		So(HasWin(f, p2), ShouldBeFalse)
		f.Field[1][0].Owner = p2
		f.Field[1][0].Population = 5
		So(HasWin(f, p1), ShouldBeFalse)
	})
}

func TestReward(t *testing.T) {
	f := field.NewField(3, 4)
	p1 := player.NewPlayer("test1")
	p2 := player.NewPlayer("test2")
	p3 := player.NewPlayer("test3")
	f.Field[0][0].Owner = p1
	f.Field[0][0].Population = 5
	f.Field[0][1].Owner = p1
	f.Field[0][1].Population = 5
	f.Field[1][1].Owner = p2
	f.Field[1][1].Population = 200

	Convey("Test reward", t, func() {
		So(Reward(f, p1), ShouldEqual, 2)
		So(Reward(f, p2), ShouldEqual, 1)
		So(Reward(f, p3), ShouldEqual, 0)
	})

	Convey("Test reward limit", t, func() {
		f := field.NewField(config.MAX_REWARD+4, 2)
		p := player.NewPlayer("test1")
		config.NB_PLAYERS = config.MAX_REWARD + 2

		for i := 0; i < config.MAX_REWARD+2; i++ {
			f.PlacePlayer(p, 10)
		}

		So(Reward(f, p), ShouldEqual, config.MAX_REWARD)
	})
}
