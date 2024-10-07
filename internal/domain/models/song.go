package models

import "time"

type Song struct {
	Id          int64
	GroupName   Group
	SongName    string
	ReleaseDate time.Time
	Lyrics      string
	Link        string
}

type Group struct {
	Id        int64
	GroupName string
}

type SongCreateReq struct {
	GroupName string `json:"group_name"`
	SongName  string `json:"song_name"`
}

type SongUpdateReq struct {
	Id           int64  `json:"id"`
	NewGroupName string `json:"new_group_name"`
	NewSongName  string `json:"new_song_name"`
}

type SongCreateResponse struct {
	Message string `json:"message"`
}

type SongFilter struct {
	GroupName   string
	SongName    string
	ReleaseDate string
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type LyricsResponse struct {
	SongID   int64    `json:"song_id"`
	Title    string   `json:"title"`
	Group    string   `json:"group"`
	Page     int      `json:"page"`
	Limit    int      `json:"limit"`
	Total    int      `json:"total"`
	Couplets []string `json:"couplets"`
}
