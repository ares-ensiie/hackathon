package field

import (
	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/event"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
)

type Field struct {
	SizeX int        `json:"sizeX"`
	SizeY int        `json:"sizeY"`
	Field [][]*Point `json:"cells"` // Coordinate : x,y
}

type Point struct {
	Owner      *player.Player `json:"owner"`
	Population int            `json:"population"`
}

func NewPoint() *Point {
	return &Point{
		Owner:      player.GAIA,
		Population: 0,
	}
}

func NewField(sx, sy int) *Field {
	field := make([][]*Point, sx)
	for x := 0; x < sx; x++ {
		field[x] = make([]*Point, sy)
		for y := 0; y < sy; y++ {
			field[x][y] = NewPoint()
		}
	}

	return &Field{
		SizeX: sx,
		SizeY: sy,
		Field: field,
	}
}

func (f *Field) PlacePlayer(player *player.Player, initialPopulation int) *event.InitPlacement {
	cont := true
	var initPlacement *event.InitPlacement

	if config.NB_PLAYERS == 2 {
		x := 0
		y := 0
		if !f.Field[0][0].Owner.IsGaia() {
			x = f.SizeX - 1
			y = f.SizeY - 1
		}
		f.Field[x][y].Owner = player
		if x == 0 && y == 0 {
			f.Field[x][y].Population = 2
		} else {
			f.Field[x][y].Population = 3
		}
		initPlacement = &event.InitPlacement{
			Player: player,
			X:      x,
			Y:      y,
			Pop:    f.Field[x][y].Population,
		}

		return initPlacement

	}

	for cont {
		x := config.RNG.Int() % f.SizeX
		y := config.RNG.Int() % f.SizeY
		if f.Field[x][y].Owner.IsGaia() {
			cont = false
			f.Field[x][y].Owner = player
			f.Field[x][y].Population = initialPopulation
			initPlacement = &event.InitPlacement{
				Player: player,
				X:      x,
				Y:      y,
				Pop:    f.Field[x][y].Population,
			}
		}
	}

	return initPlacement
}
