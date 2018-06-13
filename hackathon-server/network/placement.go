package network

import (
	"strconv"

	"git.ares-ensiie.eu/hackathon/hackathon-server/plugin"

	"gopkg.in/errgo.v1"
)

func (c *Client) Place() error {
	buffer := make([]byte, 1)
	rewards := c.Game.RemainingReward
	buffer[0] = byte(rewards)
	_, err := c.Conn.Write(buffer)

	if err != nil {
		return errgo.Mask(err)
	}

	buffer = make([]byte, 2*rewards)
	_, err = c.Conn.Read(buffer)
	if err != nil {
		return errgo.Mask(err)
	}

	for i := 0; i < rewards; i++ {
		px := int(buffer[2*i])
		py := int(buffer[2*i+1])
		if px != 255 && py != 255 {
			c.Logger.Info("[PLACE] Placing at: (" + strconv.Itoa(px) + " , " + strconv.Itoa(py) + ")")
			placement, err := c.Game.PlaceUnit(c.Player, px, py)
			if err != nil {
				c.Logger.Warn("[PLACE] Error: " + err.Error())
			}
			plugin.OnPlacement(placement)
		}
	}

	err = c.SendField()
	if err != nil {
		return errgo.Mask(err)
	}

	return nil
}
