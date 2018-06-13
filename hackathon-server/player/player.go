package player

import "sync"

var curID = 1
var IDMutex *sync.Mutex = &sync.Mutex{}

type Player struct {
	ID   int    `json:"ID"`
	Name string `json:"name"`
}

var GAIA = &Player{
	ID:   0,
	Name: "Gaia",
}

func NewPlayer(name string) *Player {
	IDMutex.Lock()
	id := curID
	curID++
	IDMutex.Unlock()

	return &Player{
		ID:   id,
		Name: name,
	}
}

func (p *Player) IsGaia() bool {
	return p.Equals(GAIA)
}

func (p *Player) Equals(p1 *Player) bool {
	return p.ID == p1.ID
}

func ResetPlayer() {
	curID = 1
}
