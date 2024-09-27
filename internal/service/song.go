package service

import (
	"context"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/lib/logger/sl"
	"log/slog"
)

func (s *SongImpl) CreateSong(
	ctx context.Context, req models.SongCreateReq) (string, error) {

	const op = "service.song.Create"

	log := s.log.With(
		slog.String("op: ", op),
	)

	msg, err := s.songImpl.Save(ctx, req.GroupName, req.SongName)

	if err != nil {
		log.Error("failed create song", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the song was created")

	return msg, nil
}

func (s *SongImpl) UpdateSong(
	ctx context.Context, req models.SongUpdateReq) (models.Song, error) {

	const op = "service.song.UpdateSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	song, err := s.songImpl.Update(ctx, req.Id, req.NewGroupName, req.NewSongName)

	if err != nil {
		log.Error("failed update song", sl.Err(err))

		return models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("the song was updated")

	return song, nil
}

func (s *SongImpl) GetSongByID(
	ctx context.Context, id int64) (models.Song, error) {

	const op = "service.song.GetSongById"

	log := s.log.With(
		slog.String("op: ", op),
	)
	song, err := s.songImpl.GetById(ctx, id)
	if err != nil {
		log.Error("Dont get song by id", sl.Err(err))

		return models.Song{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("The song was found")

	return song, nil
}

func (s *SongImpl) DeleteSong(
	ctx context.Context, id int64) (string, error) {

	const op = "service.song.DeleteSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	msg, err := s.songImpl.Remove(ctx, id)
	if err != nil {
		log.Error("Dont get song by id", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("The song was Remove")
	return msg, err
}

func (s *SongImpl) GetAllSong(
	ctx context.Context) (songs []*models.Song, err error) {

	const op = "service.song.GetAllSong"

	log := s.log.With(
		slog.String("op: ", op),
	)

	songs, err = s.songImpl.GetAll(ctx)

	if err != nil {
		log.Error("failed to get all songs", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Get all songs")
	return songs, nil
}
