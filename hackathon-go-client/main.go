package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"git.ares-ensiie.eu/hackathon/hackathon-go-client/client"
)

func ce(e error) {
	if e != nil {
		log.Println("--- ERROR ---")
		panic(e)
	}
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := client.NewClient()
	ce(c.Connect("127.0.0.1:1337", "Slt"+strconv.Itoa(rng.Int()%10)))

	log.Println("Connected ! Name: " + c.Name + ", ID: " + strconv.Itoa(c.ID))

	for c.Status() == client.ONGOING {
		log.Println("Waiting for our turn")
		ce(c.NextTurn())

		if c.Status() != client.ONGOING {
			break
		}

		log.Println("Attacking")
		for i := 0; i < 10; i++ {
			mycell := c.MyCells()[rng.Int()%len(c.MyCells())]
			if mycell.Power >= 2 {
				dx := mycell.X + rng.Int()%3 - 1
				dy := mycell.Y + rng.Int()%3 - 1
				if dx >= 0 && dx < c.GetField().SizeX && dy >= 0 && dy < c.GetField().SizeY {
					dest := c.Get(dx, dy)
					if dest != nil && dest.Owner != c.ID {
						ce(c.Attack(mycell.X, mycell.Y, dx, dy))
					}
				}
			}
		}
		log.Println("End attacks")
		_, err := c.EndAttacks()
		ce(err)
		if c.Status() != client.ONGOING {
			break
		}

		log.Println("Adding units")
		log.Println(c.RemainingUnits())

		for c.RemainingUnits() > 0 {
			mycell := c.MyCells()[rng.Int()%len(c.MyCells())]
			c.AddUnit(mycell)
		}
		ce(c.EndAddingUnits())
	}

	c.GetField().Print()

	switch c.Status() {
	case client.DEFEAT:
		log.Println("WE LOST :(")
		break
	case client.VICTORY:
		log.Println("WE WON !")
		break
	case client.CONNECTION_LOST:
		log.Println("CONNECTION LOST : WE PROBABLY LOST :(")
		break
	case client.ONGOING:
		log.Println("WTF !!")
		break
	}
}
