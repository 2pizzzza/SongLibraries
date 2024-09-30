package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/storage"
	"log"
	"time"
)

func (s *Storage) Save(
	ctx context.Context, groupName, songName string) (string, error) {

	const op = "postgres.song.Save"

	var exists bool
	err := s.Db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM testtask.public.songs WHERE group_name = $1 AND song_title = $2)",
		groupName, songName).Scan(&exists)

	if err != nil {
		log.Printf("failed to check existence: %v op: %s", err, op)
		return "", fmt.Errorf("%s, %w", op, err)
	}

	if exists {
		log.Printf("Song already exists: group %s, song %s", groupName, songName)
		return "Song already exists", fmt.Errorf("the song '%s' by group '%s' already exists", songName, groupName)
	}

	_, err = s.Db.Exec(
		"INSERT INTO testtask.public.songs (group_name, song_title) VALUES($1, $2)",
		groupName, songName)

	if err != nil {
		log.Printf("failed to create song: %v op: %s", err, op)
		return "", fmt.Errorf("%s, %w", op, err)
	}

	log.Println("Song created successfully: group ", groupName, "song", songName)
	return "Success create song", nil
}

func (s *Storage) GetById(
	ctx context.Context, id int64) (models.Song, error) {

	const op = "postgres.song.GetById"

	var (
		idTemp      int64
		group       string
		song        string
		releaseDate time.Time
		lyrics      string
		link        string
	)

	stmt, err := s.Db.Prepare("SELECT * FROM testtask.public.songs WHERE id = $1")

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	err = stmt.QueryRow(id).Scan(&idTemp, &group, &song, &releaseDate, &lyrics, &link)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Song{}, storage.ErrSongNotFound
		}
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	return models.Song{
		Id:          idTemp,
		GroupName:   group,
		SongName:    song,
		ReleaseDate: releaseDate,
		Lyrics:      lyrics,
		Link:        link,
	}, nil
}

func (s *Storage) Update(
	ctx context.Context, id int64, newGroupName, newSongName string) (models.Song, error) {

	const op = "postgres.song.Update"

	stmt, err := s.Db.Prepare("UPDATE testtask.public.songs SET group_name = $2, song_title = $3 WHERE id = $1")

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	res, err := stmt.Exec(id, newGroupName, newSongName)
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	if rowAffected == 0 {
		return models.Song{}, storage.ErrSongExists
	}

	song, err := s.GetById(ctx, id)
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	return song, nil
}

func (s *Storage) Remove(
	ctx context.Context, id int64) (string, error) {

	const op = "postgres.song.Remove"

	stmt, err := s.Db.Prepare("DELETE FROM testtask.public.songs WHERE id = $1")

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return "", storage.ErrSongNotFound
	}

	return fmt.Sprintf("Successfully deleted song id: %d", id), nil
}

func (s *Storage) GetAll(
	ctx context.Context) (songs []*models.Song, err error) {

	const op = "postgres.song.GetAll"

	var (
		id          int64
		group       string
		song        string
		releaseDate time.Time
		lyrics      string
		link        string
	)

	rows, err := s.Db.Query("SELECT * FROM testtask.public.songs ORDER BY id DESC ")

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	if err != nil {
		log.Printf("failed to query do: %v", err)
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&id, &group, &song, &releaseDate, &lyrics, &link); err != nil {
			log.Fatal(err)
		}
		songs = append(songs, &models.Song{
			Id:          id,
			GroupName:   group,
			SongName:    song,
			ReleaseDate: releaseDate,
			Lyrics:      lyrics,
			Link:        link,
		})
	}

	return songs, nil
}
