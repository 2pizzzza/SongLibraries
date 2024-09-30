package handlers

import (
	"context"
	"fmt"
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
		http.Error(w, fmt.Sprintf("%e", err, msg), http.StatusInternalServerError)
		return
	}

	resp := models.SongCreateResponse{
		Message: msg,
	}

	utils.WriteResponseBody(w, resp)
}
