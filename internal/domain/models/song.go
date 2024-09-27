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
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Lyrics      string
	Link        string
}

type SongUpdateReq struct {
	Id           int64
	NewGroupName string
	NewSongName  string
}
