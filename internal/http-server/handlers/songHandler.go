package handlers

import (
	"context"
	"errors"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/service"
	"github.com/2pizzzza/TestTask/internal/storage"
	"github.com/2pizzzza/TestTask/internal/utils"
	"log"
	"net/http"
	"strconv"
)

type Handlers struct {
	SongService service.SongService
}

func New(songService service.SongService) *Handlers {
	return &Handlers{SongService: songService}
}

// CreateSong godoc
// @Summary Create a new song
// @Description Create a new song in the library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.SongCreateReq true "Song data"
// @Success 201 {object} models.SongCreateResponse "Song created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/create [post]
func (h *Handlers) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SongCreateReq
	if err := utils.ReadRequestBody(r, &req); err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Invalid request body"}, http.StatusBadRequest)
		return
	}

	songCreateReq := models.SongCreateReq{
		GroupName: req.GroupName,
		SongName:  req.SongName,
	}

	msg, err := h.SongService.CreateSong(context.Background(), songCreateReq)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Failed to create song"}, http.StatusInternalServerError)
		return
	}

	resp := models.SongCreateResponse{
		Message: msg,
	}

	utils.WriteResponseBody(w, resp, 201)
}

// UpdateSong godoc
// @Summary Update an existing song
// @Description Update the details of a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.SongUpdateReq true "Updated song data"
// @Success 200 {object} models.Song "Song updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 404 {object} models.ErrorResponse "Song not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/update [put]
func (h *Handlers) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SongUpdateReq
	utils.ReadRequestBody(r, &req)

	songUpdateReq := models.SongUpdateReq{
		Id:           req.Id,
		NewGroupName: req.NewGroupName,
		NewSongName:  req.NewSongName,
	}

	song, err := h.SongService.UpdateSong(context.Background(), songUpdateReq)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Failed to update song"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, song, 200)
}

// GetSongByID godoc
// @Summary Get a song by ID
// @Description Retrieve the details of a song by its ID
// @Tags songs
// @Produce json
// @Param id query int true "Song ID"
// @Success 200 {object} models.Song "Song found"
// @Failure 404 {object} models.ErrorResponse "Song not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/info [get]
func (h *Handlers) GetSongByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Missing song ID"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Invalid song ID"}, http.StatusBadRequest)
		return
	}

	song, err := h.SongService.GetSongByID(context.Background(), id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			utils.WriteResponseBody(w, models.ErrorResponse{Message: "Song not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Failed to get song"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, song, 200)
}

// DeleteSong godoc
// @Summary Delete a song by ID
// @Description Remove a song from the library by its ID
// @Tags songs
// @Produce json
// @Param id query int true "Song ID"
// @Success 200 {string} string "Successfully deleted song"
// @Failure 404 {object} models.ErrorResponse "Song not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs/delete [delete]
func (h *Handlers) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	log.Print("asdsa", idStr)

	if idStr == "" {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Missing song ID"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Invalid song ID"}, http.StatusBadRequest)
		return
	}

	msg, err := h.SongService.DeleteSong(context.Background(), id)
	if err != nil {
		if errors.Is(err, storage.ErrSongNotFound) {
			utils.WriteResponseBody(w, models.ErrorResponse{Message: "Song not found"}, http.StatusNotFound)
			return
		}
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Failed to delete song"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, map[string]string{"message": msg}, 200)
}

// GetAllSong godoc
// @Summary Get all songs
// @Description Retrieve all songs with optional filtering and pagination
// @Tags songs
// @Produce json
// @Param group_name query string false "Group name filter"
// @Param song_name query string false "Song title filter"
// @Param release_date query string false "Release date filter"
// @Param limit query int false "Number of songs to return"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} models.Song "List of songs"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /songs [get]
func (h *Handlers) GetAllSongsHandler(w http.ResponseWriter, r *http.Request) {
	filter := models.SongFilter{
		GroupName:   r.URL.Query().Get("group_name"),
		SongName:    r.URL.Query().Get("song_name"),
		ReleaseDate: r.URL.Query().Get("release_date"),
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	songs, err := h.SongService.GetAllSong(context.Background(), filter, limit, offset)
	if err != nil {
		utils.WriteResponseBody(w, models.ErrorResponse{Message: "Failed to fetch songs"}, http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, songs, 200)
}
