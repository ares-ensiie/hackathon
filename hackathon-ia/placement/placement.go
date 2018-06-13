package placement

import "git.ares-ensiie.eu/hackathon/hackathon-go-client/client"

type Placement struct {
	Cell  *client.Cell
	Score int
}

func NewPlacement(cell *client.Cell, score int) *Placement {
	return &Placement{
		Cell:  cell,
		Score: score,
	}
}

type Placements []*Placement

func (p Placements) Len() int {
	return len(p)
}

func (p Placements) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Placements) Less(i, j int) bool {
	return p[i].Score < p[j].Score
}
