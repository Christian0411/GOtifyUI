package main

import (
	"fmt"
	"github.com/christian0411/GOtifyUI/spotify"
	c "github.com/jroimartin/gocui"
	"github.com/pkg/errors"
	spotify2 "github.com/zmb3/spotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"
)

var CLIENT_SECRET string = os.Getenv("SPOTIFY_CLIENT_SECRET")

var Config *ConfigFile
var spotifyClient *spotify2.Client
var NowPlayingInfo = spotify.NowPlayingInfo{}

func main() {

	Config, err := readConfigFile()
	if err == nil {
		spotifyClient = spotify.NewSpotify(Config.Spotify.ClientID, Config.Spotify.ClientSecret, Config.Spotify.RedirectUrl)
	}
	g, err := c.NewGui(c.OutputNormal)
	if err != nil {
		log.Println("Failed to create a GUI:", err)
		return
	}
	defer g.Close()


	g.SetManagerFunc(layout)

	readConfigFile()
	setupKeyBindings(g)

	tw, th := g.Size()
	lv, err := g.SetView("list", 0, 0, tw, int(math.Min(float64(th), 5.0)))
	lv.Title = " Shuffle: OFF | Repeat: OFF "
	fmt.Fprint(lv, formatNowPlaying())

	go refreshUpdates(g)
	err = g.MainLoop()
	log.Println("Main loop has finished:", err)


}

type ConfigFile struct {
	Spotify struct {
		ClientID string `yaml:"client_id"`
		ClientSecret string `yaml:"client_secret"`
		RedirectUrl string `yaml:"redirect_url"`
	} `yaml:"spotify"`
	Keybinds struct {
		NextSong string `yaml:"next_song"`
		PrevSong string `yaml:"previous_song"`
	} `yaml:"keybinds,omitempty"`
}

func readConfigFile() (ConfigFile, error){
	var cfg ConfigFile
	readFile(&cfg)
	if cfg.Spotify.ClientID == "Fill this out" {
		log.Fatal("client_id must be filled out in config file")
		os.Exit(1)
		return cfg, errors.New("Client ID missing")
	}
	if cfg.Spotify.ClientSecret == "Fill this out" {
		log.Fatal("client_secret must be filled out in config file")
		os.Exit(1)
		return cfg,  errors.New("Client ID missing")
	}
	if cfg.Spotify.RedirectUrl == "Fill this out" {
		log.Fatal("redirect_url must be filled out in config file")
		os.Exit(1)
		return cfg, errors.New("Client ID missing")
	}
	return cfg, nil
}

func readFile(cfg *ConfigFile) {
	f, err := os.Open("config.yml")
	if err != nil {
		createDefaultConfig()
		log.Fatal("Config file not found. Creating default. \nPlease fill out before running again")
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal("Config file is not yml")
	}
}

func createDefaultConfig() {
	var defaultCfg ConfigFile
	defaultCfg.Spotify.ClientID = "Fill this out"
	defaultCfg.Spotify.ClientSecret = "Fill this out"
	defaultCfg.Spotify.RedirectUrl = "Fill this out"
	out, _ := yaml.Marshal(defaultCfg)
	ioutil.WriteFile("config.yml",out,0644)
}

func setupKeyBindings(g *c.Gui) {
	err := g.SetKeybinding("", rune(0x63), c.ModNone, quit)
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

}

func refreshUpdates(g *c.Gui) {
	for {
		time.Sleep(1 * time.Second)
		g.Update(func(g *c.Gui) error {
			out,_ := g.View("list")
			out.Clear()
			NowPlayingInfo.RefreshNowPlaying(spotifyClient)
			if NowPlayingInfo.Playing {
				out.Title = " Shuffle: OFF | Repeat: OFF "
			} else {
				out.Title = " Shuffle: OFF | Repeat: OFF "
			}
			fmt.Fprint(out,formatNowPlaying())
			return nil
		})

	}
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
	currentMinutes, currentSeconds := formatTime(NowPlayingInfo.CurrentTime)
	totalMinutes, totalSeconds := formatTime(NowPlayingInfo.SongLength)
	return fmt.Sprintf("\x1b[0;31m %s  - %s \t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t%02d:%02d / %02d:%02d", NowPlayingInfo.SongName, NowPlayingInfo.ArtistName,
		currentMinutes, currentSeconds,totalMinutes, totalSeconds )
}

func formatTime(ms float64) (int, int){
	minutes := int(ms)
	secondsDecimal := math.Mod(ms, 1)
	seconds := int(secondsDecimal * 60)

	return minutes, seconds
}

func quit(g *c.Gui, v *c.View) error  {
	return c.ErrQuit
}
