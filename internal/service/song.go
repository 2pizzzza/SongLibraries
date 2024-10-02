package service

import (
	"context"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/lib/logger/sl"
	"log/slog"
	"strings"
)

func (s *SongRep) CreateSong(
	ctx context.Context, req models.SongCreateReq) (string, error) {

	const op = "service.song.Create"

	log := s.log.With(
		slog.String("op: ", op),
	)

	msg, err := s.songRep.Save(ctx, req.GroupName, req.SongName)

	if err != nil {
		log.Error(msg, sl.Err(err))

		return msg, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the song was created")

	return msg, nil
}

func (s *SongRep) UpdateSong(
	ctx context.Context, req models.SongUpdateReq) (models.Song, error) {

	const op = "service.song.UpdateSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	song, err := s.songRep.Update(ctx, req.Id, req.NewGroupName, req.NewSongName)

	if err != nil {
		log.Error("failed update song", sl.Err(err))

		return models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the song was updated")

	return song, nil
}

func (s *SongRep) GetSongByID(
	ctx context.Context, id int64) (models.Song, error) {

	const op = "service.song.GetSongById"

	log := s.log.With(
		slog.String("op: ", op),
	)
	song, err := s.songRep.GetById(ctx, id)
	if err != nil {
		log.Error("Dont get song by id", sl.Err(err))

		return models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("The song was found")

	return song, nil
}

func (s *SongRep) DeleteSong(
	ctx context.Context, id int64) (string, error) {

	const op = "service.song.DeleteSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	msg, err := s.songRep.Remove(ctx, id)
	if err != nil {
		log.Error("Dont get song by id", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("The song was Remove")
	return msg, err
}

func (s *SongRep) GetAllSong(ctx context.Context, filter models.SongFilter, limit, offset int) (songs []*models.Song, err error) {
	const op = "service.song.GetAllSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	songs, err = s.songRep.GetAll(ctx, filter, limit, offset)
	if err != nil {
		log.Error("failed to get songs", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("fetched all songs")
	return songs, nil
}

func (s *SongRep) GetLyricsByIDWithPagination(
	ctx context.Context, id int64, page, limit int) (models.LyricsResponse, error) {

	const op = "service.song.GetLyricsByIDWithPagination"

	log := s.log.With(
		slog.String("op", op),
	)

	song, err := s.songRep.GetById(ctx, id)
	if err != nil {
		log.Error("Failed to get song by id", sl.Err(err))
		return models.LyricsResponse{}, fmt.Errorf("%s: %w", op, err)
	}

	couplets := strings.Split(song.Lyrics, "\n\n")
	totalCouplets := len(couplets)

	start := (page - 1) * limit
	if start > totalCouplets {
		return models.LyricsResponse{}, fmt.Errorf("page out of range")
	}
	end := start + limit
	if end > totalCouplets {
		end = totalCouplets
	}

	resp := models.LyricsResponse{
		SongID:   song.Id,
		Title:    song.SongName,
		Group:    song.GroupName,
		Page:     page,
		Limit:    limit,
		Total:    totalCouplets,
		Couplets: couplets[start:end],
	}

	log.Info("Lyrics retrieved successfully")

	return resp, nil
}
