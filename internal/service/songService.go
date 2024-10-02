package service

import (
	"context"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"log/slog"
)

type SongRep struct {
	log     *slog.Logger
	songRep SongRepository
}

type SongService interface {
	CreateSong(ctx context.Context, req models.SongCreateReq) (string, error)
	UpdateSong(ctx context.Context, req models.SongUpdateReq) (models.Song, error)
	GetSongByID(ctx context.Context, id int64) (models.Song, error)
	DeleteSong(ctx context.Context, id int64) (string, error)
	GetAllSong(ctx context.Context, filter models.SongFilter, limit, offset int) (songs []*models.Song, err error)
	GetLyricsByIDWithPagination(ctx context.Context, id int64, page, limit int) (models.LyricsResponse, error)
}

type SongRepository interface {
	Save(ctx context.Context, groupName, songName, releaseDate, link string) (string, error)
	GetById(ctx context.Context, id int64) (models.Song, error)
	Update(ctx context.Context, id int64, newGroupName, newSongName string) (models.Song, error)
	Remove(ctx context.Context, id int64) (string, error)
	GetAll(ctx context.Context, filter models.SongFilter, limit, offset int) (songs []*models.Song, err error)
}

func New(
	log slog.Logger,
	song SongRepository) *SongRep {
	return &SongRep{
		log:     &log,
		songRep: song,
	}
}
