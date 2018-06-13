package events

import "git.ares-ensiie.eu/hackathon/hackathon-server/player"

type Attack struct {
	Attacker    *player.Player `json:"attacker"`
	FromX       int            `json:"fromX"`
	FromY       int            `json:"fromY"`
	ToX         int            `json:"toX"`
	ToY         int            `json:"toY"`
	IsWon       bool           `json:"isWon"`
	AttackerPop int            `json:"attackerPop"`
	DefenderPop int            `json:"defenderPop"`
}
