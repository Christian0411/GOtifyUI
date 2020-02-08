package main

import (
	"fmt"
	"github.com/christian0411/GOtifyUI/spotify"
	"github.com/christian0411/GOtifyUI/panes"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

const SPOTIFY_AUTH_URL string = "https://accounts.spotify.com/authorize"
const SPOTIFY_TOKEN_URL string = "https://accounts.spotify.com/api/token"
const REDIRECT_URL string  = "http://localhost:8888/callback"

var CLIENT_ID string = os.Getenv("SPOTIFY_CLIENT_ID")
var CLIENT_SECRET string = os.Getenv("SPOTIFY_CLIENT_SECRET")

var spotifyClient = spotify.NewSpotify(CLIENT_ID, CLIENT_SECRET, REDIRECT_URL)
var NowPlaying, _ = spotifyClient.PlayerCurrentlyPlaying()

func main() {
	fmt.Printf(NowPlaying.Item.Name)
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	
	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, next); err != nil {
		log.Panicln(err)
	}


	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	fmt.Printf("\nGood Bye!")
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+100, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		NowPlaying, _ = spotifyClient.PlayerCurrentlyPlaying()
		fmt.Fprintln(v, NowPlaying.Item.Name)
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func next(g *gocui.Gui, v *gocui.View) error {
	spotifyClient.Next()
	return nil
}