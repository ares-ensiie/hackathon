package client

import "errors"

func (c *Client) initPlacement(count int) {
	c.unitsToPlace = count
	c.unitCursor = 0
	c.units = make([]*Cell, count)
}

// AddUnit ajoute une unité sur la cellule passée en paramètre.
//
// Si la fonction renvoie une erreur cela signifie que vous avez dépassé
// le nombre maximum d'unités.
//
// Le champ de bataille sera mis à jour uniquement après l'appel de la
// méthode EndAddingUnits.
func (c *Client) AddUnit(cell *Cell) error {
	if c.unitCursor == c.unitsToPlace {
		return errors.New("Cannot place more units")
	}
	c.units[c.unitCursor] = cell
	c.unitCursor++
	return nil
}

// AddUnits ajoute count unités sur la cellule passée en paramètre.
//
// Si la fonction renvoie une erreur cela signifie que vous avez dépassé
// le nombre maximum d'unités.
//
// Le champ de bataille sera mis à jour uniquement après l'appel de la
// méthode EndAddingUnits.
func (c *Client) AddUnits(cell *Cell, count int) error {
	for i := 0; i < count; i++ {
		err := c.AddUnit(cell)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddUnitsList ajoute une unité chacune des cellules de la liste passée en
// paramètre.
//
// Si la fonction renvoie une erreur cela signifie que vous avez dépassé
// le nombre maximum d'unités.
//
// Le champ de bataille sera mis à jour uniquement après l'appel de la
// méthode EndAddingUnits.
func (c *Client) AddUnitsList(cells []*Cell) error {
	for i := 0; i < len(cells); i++ {
		err := c.AddUnit(cells[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// RemainingUnits renvoie le nombre d'unités que vous pouvez encore placer
// pendant ce tour.
func (c *Client) RemainingUnits() int {
	return c.unitsToPlace - c.unitCursor
}

// EndAddingUnits termine la phase de placement des unités. Cette fonction
// va envoyer vos choix au serveur et mettre à jour la map.
//
// Si cette fonction renvoie une erreur c'est qu'il y a eu une erreur réseau
func (c *Client) EndAddingUnits() error {
	if c.unitsToPlace == 0 {
		return nil
	}

	buffer := make([]byte, 2*c.unitsToPlace)
	for i := 0; i < c.unitsToPlace; i++ {
		if i < c.unitCursor {
			buffer[2*i] = byte(c.units[i].X)
			buffer[2*i+1] = byte(c.units[i].Y)
		} else {
			buffer[2*i] = 255
			buffer[2*i+1] = 255
		}
	}

	_, err := c.conn().Write(buffer)

	if err != nil {
		c.disconnect()
		return err
	}

	err = c.receiveField()
	if err != nil {
		c.disconnect()
		return err
	}

	return nil
}
