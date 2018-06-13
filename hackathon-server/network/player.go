package network

import (
	"time"

	"git.ares-ensiie.eu/hackathon/hackathon-server/game"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
)

func WaitForStep(g *game.Game, p *player.Player, step int) bool {
	for {
		time.Sleep(1 * time.Second)
		if g.HasWin(p) || g.HasLost(p) {
			return true
		}
		curPlayer := g.Player()
		if curPlayer != nil {
			if curPlayer.Equals(p) && g.Turn == step {
				return false
			}
		}
	}
}
