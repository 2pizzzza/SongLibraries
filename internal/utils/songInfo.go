package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type TrackInfo struct {
	ReleaseDate string `json:"Release Date"`
	SpotifyURL  string `json:"Spotify URL"`
}

func FetchTrackInfo(songName string, artistName string) (*TrackInfo, error) {
	url := fmt.Sprintf("https://tempserver-chbp.onrender.com/search?song=%s&artist=%s", songName, artistName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var trackInfo TrackInfo
	if err := json.Unmarshal(body, &trackInfo); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return &trackInfo, nil
}
