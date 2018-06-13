package main

import (
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"git.ares-ensiie.eu/hackathon/hackathon-go-client/client"
	"git.ares-ensiie.eu/hackathon/hackathon-ia/attack"
	"git.ares-ensiie.eu/hackathon/hackathon-ia/placement"
)

func ce(e error) {
	if e != nil {
		log.Println("--- ERROR ---")
		log.Println(e.Error())
	}
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	c := client.NewClient()
	ce(c.Connect("192.168.1.4:1337", "Slt"+strconv.Itoa(rng.Int()%10)))

	log.Println("Connected ! Name: " + c.Name + ", ID: " + strconv.Itoa(c.ID))

	for c.Status() == client.ONGOING {
		log.Println("Waiting for our turn")
		ce(c.NextTurn())

		if c.Status() != client.ONGOING {
			break
		}

		log.Println("Attacking")

		for i := 0; i < 20; i++ {
			attacks := ListAttacks(c)
			log.Println("Possible attacks: " + strconv.Itoa(len(attacks)))
			if len(attacks) == 0 {
				break
			}
			sort.Sort(attacks)
			a := attacks[len(attacks)-1]
			if a.Score() >= 0 {
				ce(c.Attack(a.From.X, a.From.Y, a.To.X, a.To.Y))
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
		placements := ListPlacements(c)
		sort.Sort(placements)
		i := len(placements) - 1
		log.Println("Possible placements: " + strconv.Itoa(len(placements)))

		for c.RemainingUnits() > 0 {
			if i == -1 {
				i = len(placements) - 1
			}
			c.AddUnit(placements[i].Cell)
			i--
		}
		ce(c.EndAddingUnits())
		if c.Status() != client.ONGOING {
			break
		}
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

func ListAttacks(c *client.Client) attack.Attacks {
	var attacks attack.Attacks
	myCells := c.MyCells()
	for _, cell := range myCells {
		if cell.Power > 1 {
			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if cell.X+dx >= 0 && cell.X+dx < c.GetField().SizeX && cell.Y+dy >= 0 && cell.Y+dy < c.GetField().SizeY {
						if c.Get(cell.X+dx, cell.Y+dy).Owner != c.ID {
							attacks = append(attacks, attack.NewAttack(cell, c.Get(cell.X+dx, cell.Y+dy)))
							log.Println("POSSIBLE")
						} else {
							log.Println("ID REJECTED")
							log.Println(cell.Owner)
							log.Println(c.ID)
						}
					} else {
						log.Println("Position rejected")
					}
				}
			}
		} else {
			log.Println("Power rejected")
		}
	}
	return attacks
}

func ListPlacements(c *client.Client) placement.Placements {
	var placements placement.Placements
	myCells := c.MyCells()
	for _, cell := range myCells {
		score := 0
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if cell.X+dx >= 0 && cell.X+dx < c.GetField().SizeX && cell.Y+dy >= 0 && cell.Y+dy < c.GetField().SizeY {
					if cell.Owner != c.ID {
						score += cell.Power
					}
				}
			}
		}
		placements = append(placements, placement.NewPlacement(cell, score))
	}
	return placements
}
