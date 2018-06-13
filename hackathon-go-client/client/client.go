// Package client : fournit une interface simple et efficace pour communiquer
// avec le serveur du hackathon. Pas besoin de vous embêter avec le réseau
// ce package s'en charge pour vous.
//
// Le package va automatiquement gérer le réseau et les conditions gagnantes.
// La structure conseillée est la suivante:
//  func main() {
//    c := client.NewClient()
//    err := c.Connect("127.0.0.1:1337", "MonEquipe")
//    if err != nil {
//      // Traiter l'erreur
//    }
//    for c.Status() == ONGOING {
//      err = c.NextTurn()
//      if err != nil {
//        // Traiter l'erreur
//      }
//      if c.Status() != ONGOING {
//        break
//      }
//      // Prendre des descisions...
//    }
//  }
//
// Author: Ares
package client

import (
	"errors"
	"net"
)

// Directions
const (
	TOP = iota
	RIGHT
	BOTTOM
	LEFT
)

// Statuts
const (
	ONGOING = iota
	VICTORY
	DEFEAT
	CONNECTION_LOST
)

// Structure principale pour la gestion du client
type Client struct {
	c            *net.Conn
	unitsToPlace int
	unitCursor   int
	units        []*Cell
	field        *Field
	Name         string
	ID           int
	status       int
}

// NewClient : Permet d'initialiser une instance d'un client. C'est la première
// chose à faire avant de pouvoir communiquer avec le serveur.
func NewClient() *Client {
	return &Client{
		Name:         "",
		ID:           -1,
		status:       ONGOING,
		field:        nil,
		unitsToPlace: 0,
		units:        make([]*Cell, 0),
		unitCursor:   0,
	}
}

// Connect Se connecte au serveur de jeu. Il envoie le nom de l'équipe
// et récupère son identifiant numérique.
//
// L'url doit être de la forme IP:PORT.
//
// Exemple :
//   c := NewClient()
//   c.Connect("127.0.0.1:1337", "MonEquipe")
//
// Une fois appelée l'id et le nom de l'équipe sont disponibles dans les champs
// ID et Name de Client.
//   c.ID // ID De l'équipe
//   c.Name // Nom de l'équipe
func (c *Client) Connect(url string, teamName string) error {
	if len(teamName) > 24 {
		return errors.New("Team name too long")
	}

	conn, err := net.Dial("tcp", url)
	if err != nil {
		return err
	}

	var buffer = make([]byte, len(teamName)+1)
	copy(buffer[:], teamName)
	buffer[len(teamName)] = 0

	_, err = conn.Write(buffer)
	if err != nil {
		c.disconnect()
		return err
	}

	buffer = make([]byte, 1)

	_, err = conn.Read(buffer)

	if err != nil {
		c.disconnect()
		return err
	}

	c.c = &conn
	c.Name = teamName
	c.ID = int(buffer[0])

	return nil
}

func (c *Client) conn() net.Conn {
	return (*c.c)
}

// NextTurn : attends le tour suivant.
//
// Cette fonction va attendre le tour suivant, met à jour la map et calcul si
// l'on a gagné ou perdu. Il est conseillé de faire un test de status après
// l'appel à cette fonction.
//
// Si cette fonction renvoit une erreur, c'est qu'il y a eu un problème lors des
// appels réseau. Il y a de grande chance que l'on ne soit pas en état de
// continuer.
func (c *Client) NextTurn() error {
	err := c.receiveField()

	if err != nil {
		c.disconnect()
		return err
	}

	return nil
}

// Status : renvoie le status actuel du client. Les valeurs possibles sont :
//  ONGOING // La partie est en cours
//  DEFEAT // On a perdu
//  VICTORY // On a gagné
//  CONNECTION_LOST // La connexion au serveur à été perdue
func (c *Client) Status() int {
	return c.status
}

func (c *Client) disconnect() {
	c.conn().Close()
	c.status = CONNECTION_LOST
}

func (c *Client) hasWin() bool {
	return c.countOpponentCell() == 0
}

func (c *Client) hasLost() bool {
	return c.countMyCell() == 0
}
