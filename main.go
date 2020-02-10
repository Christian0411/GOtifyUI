package main

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/christian0411/GOtifyUI/spotify"
	"log"
	c "github.com/jroimartin/gocui"
	"math"
	"os"
)

const REDIRECT_URL string  = "http://localhost:8888/callback"

var CLIENT_ID string = os.Getenv("SPOTIFY_CLIENT_ID")
var CLIENT_SECRET string = os.Getenv("SPOTIFY_CLIENT_SECRET")

var spotifyClient = spotify.NewSpotify(CLIENT_ID, CLIENT_SECRET, REDIRECT_URL)
var NowPlayingInfo = spotify.NowPlayingInfo{}


func main() {

	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}
	defer g.Close()


	g.SetManagerFunc(layout)
	err = g.SetKeybinding("", c.KeyCtrlC, c.ModNone, quit)

	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	err = g.SetKeybinding("", c.KeyArrowRight, c.ModNone, next)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	err = g.SetKeybinding("", c.KeyArrowLeft, c.ModNone, previous)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}

	err = g.SetKeybinding("", c.KeySpace, c.ModNone, togglePausePlay)
	if err != nil {
		log.Println("Could not set key binding:", err)
		return
	}
		tw, th := g.Size()
		lv, err := g.SetView("list", 0, 0, tw, int(math.Min(float64(th), 5.0)))
		lv.Title = " Shuffle: OFF | Repeat: OFF "
		fmt.Fprint(lv, formatNowPlaying())
		err = g.MainLoop()
		log.Println("Main loop has finished:", err)
	}


func layout(g *c.Gui) error {
	tw, th := g.Size()

	_, err := g.SetView("list", 0, 0, tw - 5, int(math.Min(float64(th), 2.0)))
	if err != nil {
		return errors.Wrap(err, "Cannot update list view")
	}

	return nil
}

func previous(g *c.Gui, v *c.View) error {
	spotifyClient.Previous()
	v, _ = g.View("list")
	v.Clear()
	fmt.Fprint(v,formatNowPlaying())
	return nil
}


func next(g *c.Gui, v *c.View) error {
	spotifyClient.Next()
	v, _ = g.View("list")
	v.Clear()
	fmt.Fprint(v,formatNowPlaying())
	return nil

}

func togglePausePlay(g *c.Gui, v *c.View) error {
	v,_ = g.View("list")
	NowPlayingInfo.RefreshNowPlaying(spotifyClient)
	if NowPlayingInfo.Playing {
		spotifyClient.Pause()
		v.Title = " Shuffle: OFF | Repeat: OFF "
	} else {
		spotifyClient.Play()
		v.Title = " Shuffle: OFF | Repeat: OFF "
	}
	return nil
}

func formatNowPlaying() string {
	NowPlayingInfo.RefreshNowPlaying(spotifyClient)
	return fmt.Sprintf("\x1b[0;31m" + NowPlayingInfo.SongName + " - " + NowPlayingInfo.ArtistName)
}

func quit(g *c.Gui, v *c.View) error  {
	return c.ErrQuit
}
