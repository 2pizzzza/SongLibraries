package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/storage"
	"log"
)

func (s *Storage) Save(
	ctx context.Context, groupName, songName, releaseDate, link string) (string, error) {

	const op = "postgres.song.Save"

	var exists bool
	err := s.Db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM songs WHERE group_name = $1 AND song_title = $2)",
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
		"INSERT INTO songs (group_name, song_title, release_date, link) VALUES($1, $2, $3 ,$4)",
		groupName, songName, releaseDate, link)

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

	var song models.Song

	stmt, err := s.Db.Prepare("SELECT id, group_name, song_title, release_date, lyrics, link FROM songs WHERE id = $1")
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("failed to close statement: %v", closeErr)
		}
	}()

	err = stmt.QueryRow(id).Scan(&song.Id, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Lyrics, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Song{}, storage.ErrSongNotFound
		}
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	return song, nil
}

func (s *Storage) Update(
	ctx context.Context, id int64, newGroupName, newSongName string) (models.Song, error) {

	const op = "postgres.song.Update"

	stmt, err := s.Db.Prepare("UPDATE songs SET group_name = $2, song_title = $3 WHERE id = $1")
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("failed to close statement: %v", closeErr)
		}
	}()

	res, err := stmt.Exec(id, newGroupName, newSongName)
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	if rowsAffected == 0 {
		return models.Song{}, fmt.Errorf("%s: song with id %d not found", op, id)
	}

	song, err := s.GetById(ctx, id)
	if err != nil {
		return models.Song{}, fmt.Errorf("%s, %w", op, err)
	}

	return song, nil
}

func (s *Storage) Remove(ctx context.Context, id int64) (string, error) {
	const op = "postgres.song.Remove"

	stmt, err := s.Db.Prepare("DELETE FROM songs WHERE id = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("failed to close statement: %v", closeErr)
		}
	}()

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

func (s *Storage) GetAll(ctx context.Context, filter models.SongFilter, limit, offset int) (songs []*models.Song, err error) {
	const op = "postgres.song.GetAll"

	query := "SELECT * FROM songs WHERE TRUE"
	args := []interface{}{}
	argCount := 1

	if filter.GroupName != "" {
		query += fmt.Sprintf(" AND group_name ILIKE $%d", argCount)
		args = append(args, "%"+filter.GroupName+"%")
		argCount++
	}
	if filter.SongName != "" {
		query += fmt.Sprintf(" AND song_title ILIKE $%d", argCount)
		args = append(args, "%"+filter.SongName+"%")
		argCount++
	}
	if filter.ReleaseDate != "" {
		query += fmt.Sprintf(" AND release_date = $%d", argCount)
		args = append(args, filter.ReleaseDate)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := s.Db.Query(query, args...)
	if err != nil {
		log.Printf("failed to query songs: %v", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("failed to close rows: %v", closeErr)
		}
	}()

	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.Id, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Lyrics, &song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		songs = append(songs, &song)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return songs, nil
}
