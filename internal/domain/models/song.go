package models

import "time"

type Song struct {
	Id          int64
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Lyrics      string
	Link        string
}

type SongCreateReq struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"release_date"`
	Lyrics      string    `json:"lyrics"`
	Link        string    `json:"link"`
}

type SongUpdateReq struct {
	Id           int64  `json:"id"`
	NewGroupName string `json:"new_group_name"`
	NewSongName  string `json:"new_song_name"`
}

type SongCreateResponse struct {
	Message string `json:"message"`
}
