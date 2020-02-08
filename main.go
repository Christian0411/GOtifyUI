package main

import (
	"fmt"
	"github.com/christian0411/GOtifyUI/spotify"
	"os"
)

const SPOTIFY_AUTH_URL string = "https://accounts.spotify.com/authorize"
const SPOTIFY_TOKEN_URL string = "https://accounts.spotify.com/api/token"
const REDIRECT_URL string  = "http://localhost:8888/callback"

var CLIENT_ID string = os.Getenv("SPOTIFY_CLIENT_ID")
var CLIENT_SECRET string = os.Getenv("SPOTIFY_CLIENT_SECRET")

func main() {

	var spotifyClient = spotify.NewSpotify(CLIENT_ID, CLIENT_SECRET, REDIRECT_URL)
	fmt.Printf("\nGood Bye!")
}

