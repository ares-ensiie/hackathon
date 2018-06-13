package game

import (
	"testing"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewGame(t *testing.T) {
	Convey("Testing game constructor", t, func() {
		g := NewGame(2, 3, 4)
		So(g.Turn, ShouldEqual, -1)
		So(g.PlayerCount, ShouldEqual, 2)
		So(g.Players, ShouldNotBeNil)
		So(len(g.Players), ShouldEqual, 2)
		So(g.Field, ShouldNotBeNil)
		So(g.Field.SizeX, ShouldEqual, 3)
		So(g.Field.SizeY, ShouldEqual, 4)
	})
}

func TestAddPlayer(t *testing.T) {
	Convey("With value initialized", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		p3 := player.NewPlayer("test3")

		Convey("With only one player", func() {
			_, err := g.AddPlayer(p1)
			So(err, ShouldBeNil)
			So(g.Players[0], ShouldNotBeNil)
			So(g.Players[0].Equals(p1), ShouldBeTrue)
			So(g.Turn, ShouldEqual, -1)
		})

		Convey("Using multiple players", func() {
			_, err := g.AddPlayer(p1)
			So(err, ShouldBeNil)
			_, err = g.AddPlayer(p2)
			So(err, ShouldBeNil)
			So(g.Players[1], ShouldNotBeNil)
			So(g.Players[1].Equals(p2), ShouldBeTrue)
			So(g.Turn, ShouldEqual, 0)
			_, err = g.AddPlayer(p3)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestNextAvailablePlayer(t *testing.T) {
	Convey("With initial values", t, func() {
		g := NewGame(3, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		p3 := player.NewPlayer("test2")
		g.AddPlayer(p1)
		g.AddPlayer(p2)

		Convey("After initialisation", func() {
			p, e := g.nextAvailablePlayer()
			So(e, ShouldBeNil)
			So(p, ShouldEqual, 0)
		})

		Convey("Testing with game realist values", func() {
			g.CurPlayer = 0
			p, e := g.nextAvailablePlayer()
			So(e, ShouldBeNil)
			So(p, ShouldEqual, 1)

			g.CurPlayer = 2
			p, e = g.nextAvailablePlayer()
			So(e, ShouldBeNil)
			So(p, ShouldEqual, 0)
		})

		Convey("Testing when we skip users", func() {
			g.Players[2] = p3
			g.CurPlayer = 1
			p, e := g.nextAvailablePlayer()
			So(e, ShouldBeNil)
			So(p, ShouldEqual, 0)
		})

		Convey("Testing when the game is over", func() {
			g.Players[0] = p3
			g.Players[1] = p3
			g.Players[2] = p3
			g.CurPlayer = 0
			p, e := g.nextAvailablePlayer()
			So(e, ShouldNotBeNil)
			So(p, ShouldEqual, -1)
		})
	})
}

func TestNextTurn(t *testing.T) {
	Convey("With initial values", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		g.AddPlayer(p1)

		Convey("Testing a complete run", func() {
			g.Players[1] = p2
			g.Field.PlacePlayer(p2, 1)
			g.NextTurn(nil)
			So(g.CurPlayer, ShouldEqual, 0)
			So(g.Turn, ShouldEqual, ATTACK)

			g.NextTurn(nil)
			So(g.CurPlayer, ShouldEqual, 0)
			So(g.Turn, ShouldEqual, PLACEMENT)

			g.NextTurn(nil)
			So(g.CurPlayer, ShouldEqual, 1)
			So(g.Turn, ShouldEqual, ATTACK)

			g.NextTurn(nil)
			So(g.CurPlayer, ShouldEqual, 1)
			So(g.Turn, ShouldEqual, PLACEMENT)

			g.NextTurn(nil)
			So(g.CurPlayer, ShouldEqual, 0)
			So(g.Turn, ShouldEqual, ATTACK)
		})

		Convey("Testing with the wrong player", func() {
			g.AddPlayer(p2)
			So(g.NextTurn(p2), ShouldNotBeNil)
			So(g.NextTurn(p1), ShouldBeNil)
		})
	})
}

func TestPlayer(t *testing.T) {
	Convey("Testing Player", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		g.AddPlayer(p1)
		So(g.Player(), ShouldBeNil)
		g.AddPlayer(p2)
		So(g.Player(), ShouldNotBeNil)
		So(g.Player().Equals(p1), ShouldBeTrue)
	})
}

func TestPlaceUnit(t *testing.T) {
	Convey("Testing with invalid parameters", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		g.Field.Field[0][0].Owner = p2
		g.Field.Field[0][0].Population = 2

		_, err := g.PlaceUnit(p1, 0, 0)
		So(err, ShouldNotBeNil)
		g.Turn = PLACEMENT

		g.CurPlayer = 0
		g.Players[0] = p2

		_, err = g.PlaceUnit(p1, 0, 0)
		So(err, ShouldNotBeNil)
		g.Players[0] = p1

		_, err = g.PlaceUnit(p1, 0, 0)
		So(err, ShouldNotBeNil)
		g.RemainingReward = 10

		_, err = g.PlaceUnit(p1, -1, 0)
		So(err, ShouldNotBeNil)
		_, err = g.PlaceUnit(p1, 3, 0)
		So(err, ShouldNotBeNil)
		_, err = g.PlaceUnit(p1, 0, -1)
		So(err, ShouldNotBeNil)
		_, err = g.PlaceUnit(p1, 0, 4)
		So(err, ShouldNotBeNil)

		g.Players[1] = p1
		g.CurPlayer = 1
		_, err = g.PlaceUnit(p1, 0, 0)
		So(err, ShouldNotBeNil)

		g.Players[1] = p2
		g.Field.Field[0][0].Population = config.MAX_POP
		_, err = g.PlaceUnit(p2, 0, 0)
		So(err, ShouldNotBeNil)
	})

	Convey("Testing with valid parameters", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		g.Players[0] = p1
		g.RemainingReward = 1
		g.Field.Field[0][0].Owner = p1
		g.Field.Field[0][0].Population = 1
		g.Turn = PLACEMENT
		g.CurPlayer = 0

		_, err := g.PlaceUnit(p1, 0, 0)
		So(err, ShouldBeNil)
		So(g.RemainingReward, ShouldEqual, 0)
		So(g.Field.Field[0][0].Population, ShouldEqual, 2)
	})
}

func TestAttackGame(t *testing.T) {
	Convey("With initial values", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		g.CurPlayer = 0
		g.Players[0] = p1
		g.Players[1] = p2
		g.Field.Field[0][0].Owner = p1
		g.Field.Field[0][0].Population = 10
		g.Field.Field[1][0].Owner = p2
		g.Field.Field[1][0].Population = 10
		g.Turn = ATTACK
		g.initAttack()

		Convey("Testing with invalid turn", func() {
			g.Turn = PLACEMENT
			_, err := g.Attack(p1, 0, 0, 1, 0)
			So(err, ShouldNotBeNil)
			So(g.RemainingAttacks, ShouldEqual, config.ATTACK_PER_ROUND)
		})

		Convey("Testing with invalid player", func() {
			_, err := g.Attack(p2, 1, 0, 0, 0)
			So(err, ShouldNotBeNil)
			So(g.RemainingAttacks, ShouldEqual, config.ATTACK_PER_ROUND)
		})

		Convey("Testing without any attacks", func() {
			g.RemainingAttacks = 0
			_, err := g.Attack(p1, 0, 0, 1, 0)
			So(err, ShouldNotBeNil)
			So(g.RemainingAttacks, ShouldEqual, 0)
		})

		Convey("Testing with valid values", func() {
			_, err := g.Attack(p1, 0, 0, 1, 0)
			So(err, ShouldBeNil)
			So(g.RemainingAttacks, ShouldEqual, config.ATTACK_PER_ROUND-1)
		})

		Convey("Testing with incorrect values for the core function", func() {
			_, err := g.Attack(p1, 0, 0, 2, 0)
			So(err, ShouldNotBeNil)
			So(g.RemainingAttacks, ShouldEqual, config.ATTACK_PER_ROUND-1)
		})
	})
}

func TestHasLosAndHasWintGame(t *testing.T) {
	Convey("With initial values", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test2")
		g.Field.Field[0][0].Owner = p1
		g.Field.Field[0][0].Population = 5
		g.Turn = ATTACK

		Convey("When in attack mode", func() {
			So(g.HasLost(p2), ShouldBeTrue)
			So(g.HasWin(p1), ShouldBeTrue)
		})

		Convey("When in registration mode", func() {
			g.Turn = REGISTRATIONS
			So(g.HasLost(p2), ShouldBeFalse)
			So(g.HasWin(p1), ShouldBeFalse)
		})
	})
}

func TestDisqualify(t *testing.T) {
	g := NewGame(2, 3, 4)
	p1 := player.NewPlayer("test1")
	p2 := player.NewPlayer("test2")
	g.CurPlayer = 0
	g.Players[0] = p1
	g.Players[1] = p2
	g.Field.Field[0][0].Owner = p1
	g.Field.Field[0][0].Population = 10
	g.Field.Field[1][0].Owner = p2
	g.Field.Field[1][0].Population = 10

	Convey("Test Disqualify", t, func() {
		g.Disqualify(p1)
		So(g.Field.Field[0][0].Owner.IsGaia(), ShouldBeTrue)
		So(g.Field.Field[0][0].Population, ShouldEqual, 10)
		So(g.Field.Field[1][0].Owner.Equals(p2), ShouldBeTrue)
		So(g.Field.Field[1][0].Population, ShouldEqual, 10)
	})
}

func TestNextPlayer(t *testing.T) {
	Convey("Testing next player", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		p2 := player.NewPlayer("test1")
		g.AddPlayer(p1)
		g.AddPlayer(p2)
		g.NextTurn(nil)
		g.Turn = PLACEMENT
		So(g.NextPlayer(p2), ShouldNotBeNil)
		So(g.Player().Equals(p1), ShouldBeTrue)
		So(g.Turn, ShouldEqual, PLACEMENT)

		g.Turn = PLACEMENT
		So(g.NextPlayer(p1), ShouldBeNil)
		So(g.Player().Equals(p2), ShouldBeTrue)
		So(g.Turn, ShouldEqual, ATTACK)

		So(g.NextPlayer(p2), ShouldBeNil)
		So(g.Player().Equals(p1), ShouldBeTrue)
		So(g.Turn, ShouldEqual, ATTACK)

		g.CurPlayer = -1
		So(g.NextPlayer(p1), ShouldNotBeNil)
		So(g.Turn, ShouldEqual, ATTACK)
		So(g.CurPlayer, ShouldEqual, -1)

	})
}

func TestInitPlacement(t *testing.T) {
	Convey("Testing init placement", t, func() {
		g := NewGame(2, 3, 4)
		p1 := player.NewPlayer("test1")
		g.Field.Field[0][0].Owner = p1
		g.Players[0] = p1
		g.CurPlayer = 0
		g.Turn = PLACEMENT
		g.initPlacement()
		So(g.RemainingReward, ShouldEqual, 1)
	})
}

func TestInitAttack(t *testing.T) {
	Convey("Testing init attack", t, func() {
		g := NewGame(2, 3, 4)
		g.initAttack()
		So(g.RemainingAttacks, ShouldEqual, config.ATTACK_PER_ROUND)
	})
}
