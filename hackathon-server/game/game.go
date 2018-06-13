package game

import (
	"strconv"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/event"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
	"gopkg.in/errgo.v1"
)

var (
	REGISTRATIONS = -1
	ATTACK        = 0
	PLACEMENT     = 1
)

type Game struct {
	PlayerCount      int              `json:"playerCount"`
	Players          []*player.Player `json:"players"`
	Field            *field.Field     `json:"field"`
	Turn             int
	CurPlayer        int
	RemainingReward  int
	RemainingAttacks int
}

func NewGame(pc, sx, sy int) *Game {
	return &Game{
		PlayerCount: pc,
		Players:     make([]*player.Player, pc),
		Field:       field.NewField(sx, sy),
		Turn:        REGISTRATIONS,
		CurPlayer:   -1,
	}
}

func (g *Game) AddPlayer(p *player.Player) (*event.InitPlacement, error) {
	if g.Turn != REGISTRATIONS {
		return nil, errgo.New("Cannot add a player now.")
	}

	var initPlacement *event.InitPlacement
	for i := 0; i < g.PlayerCount; i++ {
		if i == g.PlayerCount-1 {
			g.NextTurn(nil)
		}

		if g.Players[i] == nil {
			g.Players[i] = p
			initPlacement = g.Field.PlacePlayer(p, 1)
			break
		}
	}
	return initPlacement, nil
}

func (g *Game) NextTurn(p *player.Player) error {

	if p != nil && !g.Player().Equals(p) {
		return errgo.New("Not your turn")
	}

	if g.Turn == REGISTRATIONS {
		var err error
		g.CurPlayer, err = g.nextAvailablePlayer()
		if err != nil {
			panic(err)
		}
		g.Turn = ATTACK
	} else {
		if g.Turn != PLACEMENT {
			g.Turn++
		} else {
			g.Turn = ATTACK
			g.CurPlayer, _ = g.nextAvailablePlayer()
		}
	}

	if g.Turn == PLACEMENT {
		g.initPlacement()
	}

	if g.Turn == ATTACK {
		g.initAttack()
	}

	return nil
}

func (g *Game) nextAvailablePlayer() (int, error) {
	for i := 1; i < g.PlayerCount; i++ {
		playerID := (g.CurPlayer + i) % g.PlayerCount

		if !HasLost(g.Field, g.Players[playerID]) {
			return playerID, nil
		}
	}

	return -1, errgo.New("The game is over")
}

func (g *Game) initPlacement() {
	g.RemainingReward = Reward(g.Field, g.Player())
}

func (g *Game) initAttack() {
	g.RemainingAttacks = config.ATTACK_PER_ROUND
}

func (g *Game) HasWin(p *player.Player) bool {
	if g.Turn == REGISTRATIONS {
		return false
	}
	return HasWin(g.Field, p)
}

func (g *Game) HasLost(p *player.Player) bool {
	if g.Turn == REGISTRATIONS {
		return false
	}
	return HasLost(g.Field, p)
}

func (g *Game) Player() *player.Player {
	if g.Turn == REGISTRATIONS {
		return nil
	}
	return g.Players[g.CurPlayer]
}

func (g *Game) Attack(p *player.Player, fromX, fromY, toX, toY int) (*event.Attack, error) {
	if g.Turn != ATTACK {
		return nil, errgo.New("You cannot attack now.")
	}

	if !g.Player().Equals(p) {
		return nil, errgo.New("It's not your turn")
	}

	if g.RemainingAttacks <= 0 {
		return nil, errgo.New("You cannot attack anymore")
	}

	g.RemainingAttacks--
	return Attack(g.Field, g.Player(), fromX, fromY, toX, toY)
}

func (g *Game) NextPlayer(p *player.Player) error {
	if g.CurPlayer == -1 {
		return errgo.New("Cannot process current player")
	}
	if !g.Player().Equals(p) {
		return errgo.New("Not your turn")
	}

	g.Turn = ATTACK
	g.CurPlayer, _ = g.nextAvailablePlayer()
	return nil
}

func (g *Game) PlaceUnit(p *player.Player, px, py int) (*event.Placement, error) {
	if g.Turn != PLACEMENT {
		return nil, errgo.New("You cannot place units now.")
	}

	if !g.Player().Equals(p) {
		return nil, errgo.New("It's not your turn")
	}

	if g.RemainingReward <= 0 {
		return nil, errgo.New("You cannot place any new rewards.")
	}

	if px < 0 || px >= g.Field.SizeX {
		return nil, errgo.New("Invalid starting X position : " + strconv.Itoa(px))
	}
	if py < 0 || py >= g.Field.SizeY {
		return nil, errgo.New("Invalid starting Y position : " + strconv.Itoa(py))
	}

	if !g.Field.Field[px][py].Owner.Equals(g.Player()) {
		return nil, errgo.New("You cannot give units to your opponents")
	}

	if g.Field.Field[px][py].Population >= config.MAX_POP {
		return nil, errgo.New("This cell is full")
	}
	g.Field.Field[px][py].Population++
	g.RemainingReward--

	placement := &event.Placement{
		X: px,
		Y: py,
		Player: p,
	}
	return placement, nil
}

func (g *Game) Disqualify(p *player.Player) {
	for x := 0; x < g.Field.SizeX; x++ {
		for y := 0; y < g.Field.SizeY; y++ {
			if g.Field.Field[x][y].Owner.Equals(p) {
				g.Field.Field[x][y].Owner = player.GAIA
			}
		}
	}
}
