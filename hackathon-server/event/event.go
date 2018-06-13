package event

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

type InitPlacement struct {
	Player *player.Player `json:"player"`
	X      int            `json:"x"`
	Y      int            `json:"y"`
	Pop    int            `json:"pop"`
}

type Placement struct {
	Player *player.Player `json:"player"`
	X      int            `json:"x"`
	Y      int            `json:"y"`
}
