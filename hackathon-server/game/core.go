package game

import (
	"fmt"
	"math"
	"strconv"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/event"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"

	"gopkg.in/errgo.v1"
)

func Reward(f *field.Field, p *player.Player) int {
	result := 0
	for x := 0; x < f.SizeX; x++ {
		for y := 0; y < f.SizeY; y++ {
			if f.Field[x][y].Owner.Equals(p) {
				result++
			}
		}
	}
	if result > config.MAX_REWARD {
		result = config.MAX_REWARD
	}
	return result
}

func HasLost(f *field.Field, p *player.Player) bool {
	for x := 0; x < f.SizeX; x++ {
		for y := 0; y < f.SizeY; y++ {
			if f.Field[x][y].Owner.Equals(p) {
				return false
			}
		}
	}
	return true
}

func HasWin(f *field.Field, p *player.Player) bool {
	for x := 0; x < f.SizeX; x++ {
		for y := 0; y < f.SizeY; y++ {
			if !f.Field[x][y].Owner.IsGaia() && !f.Field[x][y].Owner.Equals(p) {
				return false
			}
		}
	}
	return true
}

func Attack(f *field.Field, attacker *player.Player, fromX int, fromY int, toX int, toY int) (*event.Attack, error) {
	if fromX < 0 || fromX >= f.SizeX {
		return nil, errgo.New("Invalid starting X position : " + strconv.Itoa(fromX))
	}
	if fromY < 0 || fromY >= f.SizeY {
		return nil, errgo.New("Invalid starting Y position : " + strconv.Itoa(fromY))
	}

	if toX < 0 || toX >= f.SizeX {
		return nil, errgo.New("Invalid end X position : " + strconv.Itoa(toX))
	}
	if toY < 0 || toY >= f.SizeY {
		return nil, errgo.New("Invalid end Y position : " + strconv.Itoa(toY))
	}

	if math.Abs(float64(fromX-toX)) > 1 || math.Abs(float64(fromY-toY)) > 1 {
		return nil, errgo.New("Cannot attack from this long")
	}

	if !f.Field[fromX][fromY].Owner.Equals(attacker) {
		return nil, errgo.New("You are not the owner of this position")
	}

	if f.Field[fromX][fromY].Owner.Equals(f.Field[toX][toY].Owner) {
		return nil, errgo.New("You vannot attack yourself")
	}

	if f.Field[fromX][fromY].Population < 2 {
		return nil, errgo.New("You cannot attack with less than 2 people")
	}

	attackers := f.Field[fromX][fromY].Population - 1
	defenders := f.Field[toX][toY].Population

	for attackers > 0 && defenders > 0 {
		if config.RNG.Int()%2 == 0 {
			attackers--
		} else {
			defenders--
		}
	}

	f.Field[fromX][fromY].Population = attackers + 1
	f.Field[toX][toY].Population = defenders

	// Si l'attaquant à gagné
	if defenders == 0 {
		f.Field[toX][toY].Owner = attacker
		// Si l'attaquant à gagné mais qu'il ne lui reste plus qu'une personne.
		// Je suis quasi certain que c'est impossible mais au cas ou ...
		if f.Field[fromX][fromY].Population == 1 {
			fmt.Println("CA ARRIVE !!!!!!!")
			f.Field[fromX][fromY].Population = 1
			f.Field[toX][toY].Population = 1
		} else {
			f.Field[toX][toY].Population = f.Field[fromX][fromY].Population - 1
			f.Field[fromX][fromY].Population = 1
		}
	}

	result := &event.Attack{
		Attacker:    attacker,
		FromX:       fromX,
		FromY:       fromY,
		ToX:         toX,
		ToY:         toY,
		IsWon:       f.Field[toX][toY].Owner.Equals(attacker),
		AttackerPop: f.Field[fromX][fromY].Population,
		DefenderPop: f.Field[toX][toY].Population,
	}

	return result, nil
}
