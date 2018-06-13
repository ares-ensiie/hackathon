package client

import (
	"fmt"
	"strconv"
)

// Field décrit le champ de bataille
type Field struct {
	SizeX int
	SizeY int
	Field [][]*Cell
}

// Cell décrit une cellule du champ de bataille
type Cell struct {
	Owner int // identifiant numérique du propriétaire
	Power int // Nombre d'unités disponibles
	X     int // Position X
	Y     int // Position Y
}

func newCell(o, p, x, y int) *Cell {
	return &Cell{
		Owner: o,
		Power: p,
		X:     x,
		Y:     y,
	}
}

func (c *Client) receiveField() error {
	buffer := make([]byte, 2)
	_, err := c.conn().Read(buffer)

	if err != nil {
		return err
	}
	sizeX := int(buffer[0])
	sizeY := int(buffer[1])
	field := make([][]*Cell, sizeX)

	for x := 0; x < sizeX; x++ {
		field[x] = make([]*Cell, sizeY)
	}

	buffer = make([]byte, 2*sizeX)
	for y := 0; y < sizeY; y++ {
		_, err := c.conn().Read(buffer)
		if err != nil {
			return err
		}
		for x := 0; x < sizeX; x++ {
			field[x][y] = newCell(int(buffer[2*x]), int(buffer[2*x+1]), x, y)
		}
	}

	c.field = &Field{
		SizeX: sizeX,
		SizeY: sizeY,
		Field: field,
	}

	if c.hasLost() {
		c.status = DEFEAT
	}

	if c.hasWin() {
		c.status = VICTORY
	}

	return nil
}

func (c *Client) countMyCell() int {
	field := c.GetField()
	count := 0
	for x := 0; x < field.SizeX; x++ {
		for y := 0; y < field.SizeY; y++ {
			if field.Field[x][y].Owner == c.ID {
				count++
			}
		}
	}
	return count
}

func (c *Client) countOpponentCell() int {
	field := c.GetField()
	count := 0
	for x := 0; x < field.SizeX; x++ {
		for y := 0; y < field.SizeY; y++ {
			if field.Field[x][y].Owner != c.ID && field.Field[x][y].Owner != 0 {
				count++
			}
		}
	}
	return count
}

// Get permet de récupérer la cellule aux coordonnées x,y. Si les coordonnées
// sont invalides, alors la fonction renverra nil.
func (f *Field) Get(x, y int) *Cell {
	if x < 0 || x >= f.SizeX {
		return nil
	}

	if y < 0 || y >= f.SizeY {
		return nil
	}

	return f.Field[x][y]
}

// OwnedBy permet de récupérer la liste des cellules appartenant au propriétaire
// passé en paramètre.
func (f *Field) OwnedBy(p int) []*Cell {
	var cells []*Cell

	for x := 0; x < f.SizeX; x++ {
		for y := 0; y < f.SizeY; y++ {
			if f.Get(x, y).Owner == p {
				cells = append(cells, f.Get(x, y))
			}
		}
	}
	return cells
}

// Print the board for debuging purposes
func (f *Field) Print() {
	for y := 0; y < f.SizeY; y++ {
		line := ""
		for x := 0; x < f.SizeX; x++ {
			line += "| " + strconv.Itoa(f.Get(x, y).Owner) + "," + strconv.Itoa(f.Get(x, y).Power) + " "
		}
		fmt.Println(line)
	}
}

// GetMap alias de GetField
func (c *Client) GetMap() *Field {
	return c.GetField()
}

// GetField renvoie le champ de bataille courant
func (c *Client) GetField() *Field {
	return c.field
}

// Get alias c.GetField().Get()
func (c *Client) Get(x, y int) *Cell {
	return c.field.Get(x, y)
}

// MyCells renvoie une liste des cellules qui appartiennent au joueur courant
func (c *Client) MyCells() []*Cell {
	return c.field.OwnedBy(c.ID)
}
