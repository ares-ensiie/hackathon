package network

import (
	"net"
	"strconv"

	"git.ares-ensiie.eu/hackathon/hackathon-server/game"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
	"git.ares-ensiie.eu/hackathon/hackathon-server/plugin"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	Player *player.Player
	Game   *game.Game
	Conn   net.Conn
	Closed bool
	Done   chan bool
	Logger *log.Entry
	Turn   int
}

func NewClient(c net.Conn, g *game.Game, done chan bool) *Client {
	return &Client{
		Player: nil,
		Game:   g,
		Conn:   c,
		Closed: false,
		Done:   done,
		Logger: log.NewEntry(log.New()),
		Turn:   0,
	}
}

func (c *Client) HandleConnection() {
	// --- CLIENT REGISTRATION
	c.Logger.Debug("[NET] Starting process for: " + c.Conn.RemoteAddr().String())
	c.RegisterPlayer()

	c.Logger = log.WithFields(log.Fields{
		"name": c.Player.Name,
		"id":   c.Player.ID,
	})

	c.Logger.Debug("[NET] Player registred")
	// --- MAIN CLIENT LOOP
	for {
		c.Logger.Info("[NET] Waiting for our turn")

		// --- WAIT FOR OUR TURN
		done := WaitForStep(c.Game, c.Player, game.ATTACK)
		if done {
			if c.Game.HasLost(c.Player) {
				c.SendField()
				c.Logger.Info("[NET] Lost game.")
				c.Abort(false)
				break
			}

			if c.Game.HasWin(c.Player) {
				c.SendField()
				c.Logger.Info("[NET] Win game.")
				c.Abort(false)
				break
			}
		}
		c.Logger.Info("[NET] launching attack")

		// --- SEND FIELD
		err := c.SendField()
		if err != nil {
			c.Logger.Error("[NET] Unable to send field: " + err.Error())
			c.Abort(true)
			break
		}

		// // viewer.UpdateLabelsServer(c.Game.Field)

		if c.Game.HasLost(c.Player) {
			c.Logger.Info("[NET] Lost game.")
			// send the map
			plugin.OnField(c.Game.Field)
			c.Abort(false)
			break
		}

		if c.Game.HasWin(c.Player) {
			c.Logger.Info("[NET] Win game.")
			// send the map
			plugin.OnField(c.Game.Field)
			c.Abort(false)
			break
		}

		// --- ATTACK

		err = c.Attack()
		if err != nil {
			c.Logger.Error("[NET] Unable to attack: " + err.Error())
			c.Abort(true)
			break
		}
		c.Logger.Info("[NET] All attack done. Waiting for our placement turn.")
		// send the map
		plugin.OnField(c.Game.Field)
		c.Game.NextTurn(c.Player)

		// --- PLACEMENT

		done = WaitForStep(c.Game, c.Player, game.PLACEMENT)
		if done {
			if c.Game.HasLost(c.Player) {
				c.SendField()
				c.Logger.Info("[NET] Lost game.")
				c.Abort(false)
				break
			}

			if c.Game.HasWin(c.Player) {
				c.SendField()
				c.Logger.Info("[NET] Win game.")
				c.Abort(false)
				break
			}
		}
		c.Logger.Info("[NET] Launching placement phase")

		err = c.Place()
		if err != nil {
			c.Logger.Error("[NET] Error while placing: " + err.Error())
			c.Abort(true)
			break
		}
		c.Logger.Info("[NET] Placement phase done.")

		// send the map
		plugin.OnField(c.Game.Field)
		c.Game.NextTurn(c.Player)
		c.Turn++

	}

	c.Logger.Debug("[NET] Client: " + c.Conn.RemoteAddr().String() + " exitting ...")
}

func (c *Client) RegisterPlayer() {
	buffer := make([]byte, 25)
	count, err := c.Conn.Read(buffer)

	if err != nil {
		c.Logger.Error("[NET] Error while getting team name (client: " + c.Conn.RemoteAddr().String() + ")")
		c.Abort(true)
	}

	name := string(buffer[:count])

	c.Player = player.NewPlayer(name)

	c.Logger.Info("[NET] Client: " + c.Conn.RemoteAddr().String() + " has been registred as: " + c.Player.Name)

	initPlacement, err := c.Game.AddPlayer(c.Player)
	plugin.OnInitPlacement(initPlacement)

	if err != nil {
		c.Logger.Error("[NET] Unable to register " + name + ": " + err.Error())
		c.Abort(true)
	} else {
		buffer := make([]byte, 1)
		buffer[0] = byte(c.Player.ID)
		_, err := c.Conn.Write(buffer)

		if err != nil {
			c.Logger.Error("[NET] Unable to send ID to " + name + ": " + err.Error())
			panic("Could not continue")
		}
	}
}

func (c *Client) Abort(disqualified bool) {
	if !c.Closed {
		c.Conn.Close()
		c.Closed = true
		if disqualified {
			c.Game.Disqualify(c.Player)
		}
		err := c.Game.NextPlayer(c.Player)
		if err != nil && !c.Game.HasLost(c.Player) && !c.Game.HasWin(c.Player) {
			c.Logger.Error("[ABORT] Error: " + err.Error())
		}
		c.Logger.Info("[NET] Player removed!")
		c.Logger.Info("[GAME] Ended after " + strconv.Itoa(c.Turn) + " turns.")
		c.Done <- true
		c.Logger.Debug("[NET] Exit signal sent")
	}
}
