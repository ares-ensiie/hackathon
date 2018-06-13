package viewer

import (
	"strconv"

	"git.ares-ensiie.eu/hackathon/hackathon-go-client/client"
	"git.ares-ensiie.eu/hackathon/hackathon-server/field"
	"github.com/andlabs/ui"
)

var labels [][]*ui.Label

func Launch(s *field.Field, c *client.Field) {
	err := ui.Main(func() {
		mainBox := ui.NewHorizontalBox()
		sizeX := 0
		sizeY := 0
		name := ""
		if s == nil {
			sizeX = c.SizeX
			sizeY = c.SizeY
			name = "Client"
		} else {
			sizeX = s.SizeX
			sizeY = s.SizeY
			name = "Server"
		}
		labels = make([][]*ui.Label, sizeX)
		for x := 0; x < sizeX; x++ {
			box := ui.NewVerticalBox()
			labels[x] = make([]*ui.Label, sizeY)
			for y := 0; y < sizeY; y++ {
				if s == nil {
					pos := c.Field[x][y]
					labels[x][y] = ui.NewLabel(strconv.Itoa(pos.Owner) +
						"," + strconv.Itoa(pos.Power) + "  |  ")
				} else {
					pos := s.Field[x][y]
					labels[x][y] = ui.NewLabel(strconv.Itoa(pos.Owner.ID) +
						"," + strconv.Itoa(pos.Population) + "  |  ")
				}
				box.Append(labels[x][y], false)
			}
			mainBox.Append(box, false)
		}

		window := ui.NewWindow(name, 300, 200, false)
		window.SetChild(mainBox)
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		//panic(err)
	}
}

func UpdateLabelsServer(f *field.Field) {
	ui.QueueMain(func() {
		for x := 0; x < f.SizeX; x++ {
			for y := 0; y < f.SizeY; y++ {
				labels[x][y].SetText(strconv.Itoa(f.Field[x][y].Owner.ID) +
					"," + strconv.Itoa(f.Field[x][y].Population) + "  |  ")
			}
		}
	})
}
