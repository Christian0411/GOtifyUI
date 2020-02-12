package spotify

import (
	"context"
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var client spotify.Client

func NewSpotify(client_id, client_secret, redirect_url string) *spotify.Client {
	var auth = spotify.NewAuthenticator(redirect_url, spotify.ScopeUserReadPrivate,
		spotify.ScopeUserModifyPlaybackState,
		spotify.ScopeUserReadCurrentlyPlaying)

	auth.SetAuthInfo(client_id, client_secret)
	url := auth.AuthURL("Test")

	fmt.Printf("Please visit %s", url)

	spotifyAuthReciever := http.NewServeMux()

	portRegex, _ := regexp.Compile("\\d{1,5}")
	port := portRegex.FindAllString(redirect_url, 1)[0]

	uri := strings.Split(redirect_url, port)[1]

	s := http.Server{Addr: ":" + port, Handler: spotifyAuthReciever}

	spotifyAuthReciever.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\nRequest received")
		// use the same state string here that you used to generate the URL
		token, err := auth.Token("Test", r)
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

type NowPlayingInfo struct {
	SongName    string
	ArtistName  string
	CurrentTime float64
	SongLength  float64
	Playing     bool
}

func (npi *NowPlayingInfo) RefreshNowPlaying(client *spotify.Client){
	np, _ := client.PlayerCurrentlyPlaying()
	npi.Playing = np.Playing
	npi.SongName = np.Item.Name
	npi.ArtistName = np.Item.Artists[0].Name
	npi.CurrentTime = (float64(np.Progress) / 1000.0) / 60.0
	npi.SongLength = (float64(np.Item.SimpleTrack.Duration) / 1000.0) / 60.0
}
