package handlers

import (
	"gopher/internal/logger/logutils"
	"gopher/internal/models"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Error struct {
	Message string `json:"message"`
}

type Response struct {
	Error *Error `json:"error,omitempty"`
}

func newOK() Response {
	return Response{}
}

func newError(message string) Response {
	return Response{
		Error: &Error{
			Message: message,
		},
	}
}

type Saver interface {
	SaveNewAudioFile(models.AudioFile) error
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
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			log.Error("failed to parse form", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))

			return
		}

		// 2. Достаём файл из формы (ключ "file")
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Error("failed to form file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}
		defer file.Close()

		// 3. Читаем в память (для больших файлов лучше стримить)
		data, err := io.ReadAll(file)
		if err != nil {
			log.Error("failed to read file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}

		audioFile := models.NewAudioFile(id, header.Filename, data)

		// 4. Сохраняем через абстракцию
		if err := saver.SaveNewAudioFile(audioFile); err != nil {
			log.Error("failed to save file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}

		log.Info("audiofile saved successfully")
		render.Status(r, http.StatusCreated)
	}
}
