package plugin

import (
	"encoding/json"

	"git.ares-ensiie.eu/hackathon/hackathon-server/event"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"

	log "github.com/sirupsen/logrus"
)

var plugins []Plugin

type Plugin interface {
	OnAttack(*event.Attack)
	OnInitPlacement(*event.InitPlacement)
	OnPlacement(*event.Placement)
	OnField(*field.Field)
}

func RegisterPlugin(p Plugin) {
	plugins = append(plugins, p)
}

func OnAttack(attack *event.Attack) {
	for _, p := range plugins {
		p.OnAttack(attack)
	}
}

func OnInitPlacement(placement *event.InitPlacement) {
	log.Debug(json.Marshal(placement))
	for _, p := range plugins {
		p.OnInitPlacement(placement)
	}
}

func OnPlacement(placement *event.Placement) {
	log.Debug(json.Marshal(placement))
	for _, p := range plugins {
		p.OnPlacement(placement)
	}
}

func OnField(field *field.Field) {
	log.Debug(json.Marshal(field))
	for _, p := range plugins {
		p.OnField(field)
	}
}

func Reset() {
	plugins = make([]Plugin, 0)
}
