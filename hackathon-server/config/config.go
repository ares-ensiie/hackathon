package config

import (
	"flag"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var (
	MAX_REWARD       = 50
	ATTACK_PER_ROUND = 20
	NB_PLAYERS       = 2
	MAP_SIZE_X       = 4
	MAP_SIZE_Y       = 4
	MAX_POP          = 20
	ATTACK_TIMEOUT   = 5
	RNG              = rand.New(rand.NewSource(time.Now().UnixNano()))
	GAME_PORT        = "1337"
	WEB_PORT         = "1338"
)

func InitConfig() {
	debugLevel := flag.Bool("debug", false, "Set log level to debug")
	gamePort := flag.String("game-port", "1337", "Set the game server listenning port")
	webPort := flag.String("web-port", "1338", "Set the web server listenning port")
	mapSize := flag.Int("size", 4, "Set the map size")
	nbPlayers := flag.Int("players", 2, "Number of players")
	attTimeout := flag.Int("attack-timeout", 5, "Max attack timeout")
	//weblogs := flag.Bool("weblogs", false, "Enable logs for the web viewer")

	flag.Parse()

	if *debugLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	/*if *weblogs {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
	}*/
	log.SetFormatter(new(prefixed.TextFormatter))

	GAME_PORT = *gamePort
	WEB_PORT = *webPort
	MAP_SIZE_X = *mapSize
	MAP_SIZE_Y = *mapSize
	NB_PLAYERS = *nbPlayers
	ATTACK_TIMEOUT = *attTimeout
}
