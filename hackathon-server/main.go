package main

import (
	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/network"
)

func main() {
	config.InitConfig()
	network.Start()
}
