package network

import (
	"strconv"
	"sync"
	"time"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/plugin"

	"gopkg.in/errgo.v1"
)

var IDAttack int
var IDMutex *sync.Mutex = &sync.Mutex{}

func (c *Client) Attack() error {
	IDMutex.Lock()
	curID := IDAttack
	IDAttack++
	IDMutex.Unlock()
	timeout := make(chan int, 1)
	errChan := make(chan error, 1)
	c.Logger.Debug("[ATTACK] Timeout: " + strconv.Itoa(config.ATTACK_TIMEOUT))
	go func(id int) {
		c.Logger.Debug("Launching timeout " + strconv.Itoa(id))
		time.Sleep(time.Duration(config.ATTACK_TIMEOUT) * time.Second)
		c.Logger.Debug("Ending timeout " + strconv.Itoa(id))

		timeout <- id
	}(curID)

	go c.AttackProtocol(errChan)

	select {
	case id := <-timeout:
		if id != curID {
			c.Logger.Debug("Old timeout")
			return nil
		} else {
			return errgo.New("Attack timeout")
		}
	case err := <-errChan:
		if err != nil {
			time.Sleep(1 * time.Second)
			return errgo.Mask(err)
		}
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (c *Client) AttackProtocol(errChan chan error) {
	for c.Game.RemainingAttacks > 0 {
		stop, err := c.AttackOnce()

		if err != nil {
			errChan <- errgo.Mask(err)
		}

		if stop {
			break
		}

		err = c.SendField()
		if err != nil {
			errChan <- errgo.Mask(err)
		}
	}
	errChan <- nil
}

func (c *Client) AttackOnce() (bool, error) {
	buffer := make([]byte, 4)
	_, err := c.Conn.Read(buffer)

	if err != nil {
		return false, errgo.Mask(err)
	}

	if buffer[0] == 255 && buffer[1] == 255 && buffer[2] == 255 && buffer[3] == 255 {
		c.Logger.Info("[ATTACK] End attack packet received")
		return true, nil
	}

	c.Logger.Info("[ATTACK] Attacking coordinates: " +
		"(" + strconv.Itoa(int(buffer[0])) + "," + strconv.Itoa(int(buffer[1])) + ") , " +
		"(" + strconv.Itoa(int(buffer[2])) + "," + strconv.Itoa(int(buffer[3])) + ")")

	result, err := c.Game.Attack(c.Player, int(buffer[0]), int(buffer[1]), int(buffer[2]), int(buffer[3]))

	if err != nil {
		c.Logger.Warn("[ATTACK] FAILED: " + err.Error())
	}

	// Notify all plugins
	plugin.OnAttack(result)

	c.Logger.Debug("[ATTACK] SUCCESS!")

	return false, nil
}
