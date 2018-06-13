package network

import (
	"strconv"

	"gopkg.in/errgo.v1"
)

func (c *Client) SendField() error {
	b := make([]byte, 2)
	b[0] = byte(c.Game.Field.SizeX)
	b[1] = byte(c.Game.Field.SizeY)

	c.Logger.Debug("[FIELD] Sending field: " + strconv.Itoa(c.Game.Field.SizeX) + "," + strconv.Itoa(c.Game.Field.SizeY))

	_, err := c.Conn.Write(b)

	if err != nil {
		return errgo.Mask(err)
	}

	for y := 0; y < c.Game.Field.SizeY; y++ {
		buffer := make([]byte, 2*c.Game.Field.SizeX)
		for x := 0; x < c.Game.Field.SizeX; x++ {
			buffer[2*x] = byte(c.Game.Field.Field[x][y].Owner.ID)
			buffer[2*x+1] = byte(c.Game.Field.Field[x][y].Population)
		}

		_, err := c.Conn.Write(buffer)

		if err != nil {
			return errgo.Mask(err)
		}
	}

	return nil
}
