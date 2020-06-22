package network

import (
	"net"
	"time"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/game"
	"git.ares-ensiie.eu/hackathon/hackathon-server/logrus_socketio"
	"git.ares-ensiie.eu/hackathon/hackathon-server/player"
	"git.ares-ensiie.eu/hackathon/hackathon-server/plugin"
	"git.ares-ensiie.eu/hackathon/hackathon-server/viewer"

	log "github.com/sirupsen/logrus"
)

func Start() {

	log.Debug("[NET] Starting ...")
	ln, err := net.Listen("tcp", ":"+config.GAME_PORT)
	if err != nil {
		panic(err)
	}

	log.Info("[NET] Listening on :" + config.GAME_PORT)
	var v *viewer.Viewer

	for {
		plugin.Reset()
		player.ResetPlayer()
		g := game.NewGame(config.NB_PLAYERS, config.MAP_SIZE_X, config.MAP_SIZE_Y)
		if v == nil {
			log.Debug("[VIEWER] Launching new viewer")
			v = viewer.NewViewer(g)
			go v.Listen()

			// Logs

			time.Sleep(1 * time.Second)
			m := make(map[string]interface{})
			hook, err := logrus_socketio.NewSocketIOHook("http://localhost:"+config.WEB_PORT+"/socket.io/", "log", m)
			if err != nil {
				panic(err)
			}
			log.AddHook(hook)
		} else {
			log.Debug("[VIEWER] Updating game")
			v.SetGame(g)
		}
		plugin.RegisterPlugin(v)

		done := make(chan bool, 1)

		for i := 0; i < config.NB_PLAYERS; i++ {
			log.Debug("[NET] Waiting for connection")
			conn, err := ln.Accept()
			log.Debug("[NET] New connection from : " + conn.RemoteAddr().String())
			if err != nil {
				i--
				log.Error(err.Error())
			} else {
				go NewClient(conn, g, done).HandleConnection()
			}
		}
		for i := 0; i < config.NB_PLAYERS; i++ {
			<-done
			log.Debug("[NET] Client exited")
		}

		log.Info("[NET] Game over...")
		time.Sleep(10 * time.Second)
		log.Info("[NET] Relaunching game...")
	}
}
