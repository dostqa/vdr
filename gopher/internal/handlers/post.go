package handlers

import (
	"context"
	"gopher/internal/logger/logutils"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const maxUploadSize = 10 << 20 // 10 MB

type FileSaver interface {
	SaveFile(context.Context, string, io.Reader, int64) (string, error)
}

type MetaDataSaver interface {
	SaveRequest(context.Context) (int64, error)
	SaveFile(context.Context, int64, string) (int64, error)
}

func Post(logger *slog.Logger, metaDataSaver MetaDataSaver, fileSaver FileSaver) http.HandlerFunc {
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

		// 3. Сохраняем нужную информацию
		// 3.1 Сохраняем request
		id, err := metaDataSaver.SaveRequest(r.Context())
		if err != nil {
			log.Error("failed to create request", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}

		// 3.2 Сохраняем файл в S3 хранилище
		filepath, err := fileSaver.SaveFile(r.Context(), strconv.FormatInt(id, 10)+"_"+header.Filename, file, header.Size)
		if err != nil {
			log.Error("failed to save file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
			return
		}

		// 3.3 сохраняем файл в postgress
		_, err = metaDataSaver.SaveFile(r.Context(), id, filepath)
		if err != nil {
			log.Error("failed to create file", logutils.Err(err))

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("Internal server error"))
		}

		log.Info("audiofile saved successfully")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, newOK())
	}
}
