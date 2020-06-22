package viewer

import (
	"encoding/json"
	"net/http"

	"git.ares-ensiie.eu/hackathon/hackathon-server/config"
	"git.ares-ensiie.eu/hackathon/hackathon-server/event"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"
	"git.ares-ensiie.eu/hackathon/hackathon-server/game"

	socketio "github.com/googollee/go-socket.io"
	log "github.com/sirupsen/logrus"
)

// Viewer represent the spectator server.
type Viewer struct {
	Game   *game.Game
	Server *socketio.Server
}

func newSocketIoServer(game *game.Game) *socketio.Server {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	setOnConnectionEvent(server, game)
	server.On("log", func(msg string) {
		server.BroadcastTo("game", "game log", msg)
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Debug("[Viewer] error:", err)
	})

	return server
}

func setOnConnectionEvent(server *socketio.Server, game *game.Game) {
	server.On("connection", func(so socketio.Socket) {
		log.Debug("on connection")
		so.Join("game")

		gameJSON, _ := json.Marshal(game)
		so.Emit("game init", string(gameJSON))
	})
}

func NewViewer(game *game.Game) *Viewer {
	return &Viewer{
		Game:   game,
		Server: newSocketIoServer(game),
	}
}

func (v *Viewer) SetGame(g *game.Game) {
	v.Game = g
	gameJSON, _ := json.Marshal(v.Game)
	setOnConnectionEvent(v.Server, v.Game)
	v.Server.BroadcastTo("game", "game init", string(gameJSON))
}

func (v *Viewer) OnAttack(attack *event.Attack) {
	attackJSON, _ := json.Marshal(attack)
	v.Server.BroadcastTo("game", "game attack", string(attackJSON))
}

func (v *Viewer) OnInitPlacement(placement *event.InitPlacement) {
	placementJSON, _ := json.Marshal(placement)
	v.Server.BroadcastTo("game", "game initplacement", string(placementJSON))
}

func (v *Viewer) OnPlacement(placement *event.Placement) {
	placementJSON, _ := json.Marshal(placement)
	v.Server.BroadcastTo("game", "game placement", string(placementJSON))
}

func (v *Viewer) OnField(field *field.Field) {
	fieldJSON, _ := json.Marshal(field)
	v.Server.BroadcastTo("game", "game field", string(fieldJSON))
}

func (v *Viewer) Listen() {
	http.Handle("/socket.io/", v.Server)
	http.Handle("/", http.FileServer(http.Dir("webui")))
	log.Info("[VIEWER] Serving at localhost:" + config.WEB_PORT + "...")
	log.Fatal(http.ListenAndServe(":"+config.WEB_PORT, nil))
}
