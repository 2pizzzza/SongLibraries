package handlers

import (
	"context"
	"github.com/2pizzzza/TestTask/internal/domain/models"
	"github.com/2pizzzza/TestTask/internal/service"
	"github.com/2pizzzza/TestTask/internal/utils"
	"net/http"
)

type Handlers struct {
	SongService service.SongService
}

func New(songService service.SongService) *Handlers {
	return &Handlers{SongService: songService}
}

func (h *Handlers) CreateSongHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SongCreateReq
	utils.ReadRequestBody(r, &req)

	songCreateReq := models.SongCreateReq{
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: req.ReleaseDate,
		Lyrics:      req.Lyrics,
		Link:        req.Link,
	}

	msg, err := h.SongService.CreateSong(context.Background(), songCreateReq)
	if err != nil {
		http.Error(w, "Failed create song", http.StatusInternalServerError)
		return
	}

	resp := models.SongCreateResponse{
		Message: msg,
	}

	utils.WriteResponseBody(w, resp)
}

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
		http.Error(w, "Failed to update song", http.StatusInternalServerError)
		return
	}

	utils.WriteResponseBody(w, song)
}
