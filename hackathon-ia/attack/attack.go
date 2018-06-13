package attack

import "git.ares-ensiie.eu/hackathon/hackathon-go-client/client"

type Attack struct {
	From *client.Cell
	To   *client.Cell
}

func NewAttack(from, to *client.Cell) *Attack {
	return &Attack{
		From: from,
		To:   to,
	}
}

func (a *Attack) Score() int {
	return a.From.Power - a.To.Power
}

type Attacks []*Attack

func (a Attacks) Len() int {
	return len(a)
}

func (a Attacks) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Attacks) Less(i, j int) bool {
	return a[i].Score() < a[j].Score()
}
