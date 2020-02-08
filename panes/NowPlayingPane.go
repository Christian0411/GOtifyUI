package panes

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"strings"
)

type NowPlayingPane struct {
	name string
	x, y int
	w, h int
	body string
}

func NewNowPlayingPane(name string, x, y int, body string) *NowPlayingPane {
	lines := strings.Split(body, "\n")

	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &NowPlayingPane{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (w *NowPlayingPane) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, w.body)
	}
	return nil
}