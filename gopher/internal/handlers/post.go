package handlers

import (
	"gopher/internal/logger/logutils"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const maxUploadSize = 10 << 20 // 10 MB

type Response struct {
	Message string `json:"message"`
}

func newOK() Response {
	return Response{
		Message: "OK",
	}
}

func newError(message string) Response {
	return Response{
		Message: message,
	}
}

type Saver interface {
	SaveNewAudioFile(string, string, io.Reader) error
}

func Post(logger *slog.Logger, saver Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Post"
		id := middleware.GetReqID(r.Context())

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", id),
		)

		// 1. Парсим multipart form
		r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize) // ограничение на размер файла
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			log.Error("failed to parse form", logutils.Err(err))

			// превысили лимит
			if strings.Contains(err.Error(), "request body too large") {
				render.Status(r, http.StatusRequestEntityTooLarge)
				render.JSON(w, r, newError("file too large"))
				return
			}

			// что-то ещё пошло не так
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("invalid multipart form"))
			return
		}

		// 2. Достаём файл из формы (ключ "file")
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Error("failed to get file from form", logutils.Err(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("file is required"))
			return
		}
		defer file.Close()

		if err := saver.SaveNewAudioFile(id, header.Filename, file); err != nil {
			log.Error("failed to save file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}

		log.Info("audiofile saved successfully")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, newOK())
	}
}
