package handlers

import (
	"context"
	"gopher/internal/logger/logutils"
	"gopher/internal/models"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const maxUploadSize = 10 << 20 // 10 MB

type SaverService interface {
	Save(ctx context.Context, file models.File, r io.Reader, size int64) (int, error)
}

type MetaDataSaver interface {
	SaveRequest(context.Context) (int64, error)
	SaveFile(context.Context, int64, string) (int64, error)
}

func Post(logger *slog.Logger, saverService SaverService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Post"

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
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

		idx, err := saverService.Save(
			r.Context(),
			models.File{FileName: header.Filename},
			file,
			header.Size,
		)
		if err != nil {
			log.Error("failed to save file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("internal server error"))
			return
		}

		log.Info("file saved successfully")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, RequestIdReponse{Id: idx})
	}
}
