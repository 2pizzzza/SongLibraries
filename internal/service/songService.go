package service

import (
	"context"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"log/slog"
)

type SongImpl struct {
	log      *slog.Logger
	songImpl SongRepository
}

type SongRepository interface {
	Save(ctx context.Context, groupName, songName string) (string, error)
	GetById(ctx context.Context, id int64) (models.Song, error)
	Update(ctx context.Context, id int64, newGroupName, newSongName string) (models.Song, error)
	Remove(ctx context.Context, id int64) (string, error)
	GetAll(ctx context.Context) (songs []*models.Song, err error)
}

func New(
	log slog.Logger,
	song SongRepository) *SongImpl {
	return &SongImpl{
		log:      &log,
		songImpl: song,
	}
}
