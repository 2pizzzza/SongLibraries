package service

import (
	"context"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"log/slog"
)

type Songs struct {
	log      *slog.Logger
	songImpl SongService
}

type SongService interface {
	CreateSong(ctx context.Context, groupName, songName string) (string, error)
	GetSongById(ctx context.Context, id int64) (models.Song, error)
	UpdateSong(ctx context.Context, id int64, newGroupName, newSongName string) (models.Song, error)
	RemoveSong(ctx context.Context, id int64) (string, error)
	GetAllSongs(ctx context.Context) (songs []*models.Song, err error)
}

func New(
	log slog.Logger,
	song SongService) *Songs {
	return &Songs{
		log:      &log,
		songImpl: song,
	}
}
