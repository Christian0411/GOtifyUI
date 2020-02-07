package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const SPOTIFY_AUTH_URL string = "https://accounts.spotify.com/authorize"
const SPOTIFY_TOKEN_URL string = "https://accounts.spotify.com/api/token"
const REDIRECT_URL string  = "http://localhost:8888/callback"

var CLIENT_ID string = os.Getenv("SPOTIFY_CLIENT_ID")
var CLIENT_SECRET string = os.Getenv("SPOTIFY_CLIENT_SECRET")

var SCOPES = []string {"user-read-playback-state", "user-modify-playback-state", "user-read-currently-playing"}

func main() {
	fmt.Printf("Please authenticate with Spotify: %s" +
		"?response_type=code" +
		"&client_id=%s" +
		"&scope=%s" +
		"&redirect_uri=%s", SPOTIFY_AUTH_URL, CLIENT_ID, strings.Join(SCOPES, ","), REDIRECT_URL)

	http.HandleFunc("/callback", handleSpotifyAuthentication)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

type SpotifyTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func handleSpotifyAuthentication(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query()["code"][0]

	urlValues := url.Values{}
	urlValues.Add("grant_type", "authorization_code")
	urlValues.Add("redirect_uri", REDIRECT_URL)
	urlValues.Add("code", code)
	urlValues.Add("client_id", CLIENT_ID)
	urlValues.Add("client_secret", CLIENT_SECRET)

	data, err := http.PostForm(SPOTIFY_TOKEN_URL, urlValues)
	defer data.Body.Close()

	if err != nil {
		fmt.Printf("Could not authenticate with Spotify")
		os.Exit(-1)
	}

	var tokenResponse SpotifyTokenResponse
	err = json.NewDecoder(data.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Printf("Could not authenticate with Spotify")
		os.Exit(-1)
	}

	fmt.Printf("\nAuthenticated with Spotify.")

}
