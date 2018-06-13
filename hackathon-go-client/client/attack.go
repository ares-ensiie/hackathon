package client

// Attack permet de lancer une attaque. fromX et fromY sont les coordonnées de
// la case qui lance l'attaque et toX et toY sont les coordonnées de la case
// ciblée par l'attaque.
//
// Attention: Si les 4 valeurs sont égales à 255 cela sera interprété comme la
// fin des attaques. Si c'est ce que vous voulez faire utilisez plutot la
// méthode EndAttacks
//
// Si cette méthode renvoie une erreur c'est qu'il y a eu une erreur réseau.
func (c *Client) Attack(fromX, fromY, toX, toY int) error {
	buffer := make([]byte, 4)
	buffer[0] = byte(fromX)
	buffer[1] = byte(fromY)
	buffer[2] = byte(toX)
	buffer[3] = byte(toY)

	_, err := c.conn().Write(buffer)

	if err != nil {
		c.disconnect()
		return err
	}

	if fromX == 255 && fromY == 255 && toX == 255 && toY == 255 {
		return nil
	}

	err = c.receiveField()

	if err != nil {
		c.disconnect()
		return err
	}

	return nil
}

// EndAttacks permet de prévenir le serveur que l'on a fini toutes nos attaques.
// Cette fonction va également attendre le début de la phase de placement et
// initialiser les champs nécessaires à la phase de placement.
func (c *Client) EndAttacks() (int, error) {
	err := c.Attack(255, 255, 255, 255)
	if err != nil {
		return -1, err
	}

	buffer := make([]byte, 1)
	_, err = c.conn().Read(buffer)

	if err != nil {
		c.disconnect()
		return -1, err
	}

	c.initPlacement(int(buffer[0]))
	return c.unitsToPlace, nil
}
