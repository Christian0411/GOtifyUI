package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
)

const SPOTIFY_AUTH_URL string = "https://accounts.spotify.com/authorize"
const SPOTIFY_TOKEN_URL string = "https://accounts.spotify.com/api/token"
const REDIRECT_URL string  = "http://localhost:8888/callback"

var state string = "Test"
var auth = spotify.NewAuthenticator(REDIRECT_URL, spotify.ScopeUserReadPrivate, spotify.ScopeUserModifyPlaybackState)
var client spotify.Client

func NewSpotify(client_id, client_secret, redirect_url string) *spotify.Client {

	auth.SetAuthInfo(client_id, client_secret)
	url := auth.AuthURL(state)

	fmt.Printf("Please visit %s", url)

	spotifyAuthReciever := http.NewServeMux()

	s := http.Server{Addr: ":8888", Handler: spotifyAuthReciever}

	spotifyAuthReciever.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nRequest received")
		// use the same state string here that you used to generate the URL
		token, err := auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		// create a client using the specified token
		client  = auth.NewClient(token)
		s.Shutdown(context.Background())
	})

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return &client
}