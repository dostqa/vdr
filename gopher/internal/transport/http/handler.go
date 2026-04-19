package http

import (
	"fmt"
	"gopher/internal/service"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const maxUploadSize = 50 << 20 // 50 MB

type Handler struct {
	requestService *service.RequestService
}

func NewHandler(requestService *service.RequestService) *Handler {
	return &Handler{
		requestService: requestService,
	}
}

func (h *Handler) GetByRequestID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetByRequestID"

	rawID := chi.URLParam(r, "id")
	if rawID == "" {
		slog.Error("%s: %w", op, fmt.Errorf("request id was empty"))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("request id should not be empty"))
		return
	}

	requestID, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Error("%s: %w", op, err)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("request id is not an int"))
		return
	}

	isReady, err := h.requestService.IsReady(r.Context(), requestID)
	if err != nil {
		slog.Error("%s: %w", op, err)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("no such request"))
		return
	}

	if isReady {
		msg, err := h.requestService.OutputMessage(r.Context(), requestID)
		if err != nil {
			slog.Error("%s: %w", op, err)

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, newError("internal error"))
			return
		}
		render.Status(r, http.StatusOK)
		render.JSON(w, r, &msg)
	} else {
		render.Status(r, http.StatusProcessing)
		render.JSON(w, r, newResponse("request is not ready"))
	}
}

func (h *Handler) SaveFile(w http.ResponseWriter, r *http.Request) {
	const op = "handler.SaveFile"

	// 1. Парсим multipart form
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		slog.Error("%s: failed to parse form: %w", op, err)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("invalid multipart form"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		slog.Error("%s: failed to get file from form: %w", op, err)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, newError("file is required"))
		return
	}
	defer file.Close()

	idx, err := h.requestService.Save(
		r.Context(),
		header.Filename,
		file,
		header.Size,
	)
	if err != nil {
		slog.Error("%s: failed to save file: %w", op, err)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, newError("internal server error"))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, RequestIdReponse{Id: idx})
}

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	fileStream, err := h.requestService.File(r.Context(), name)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer fileStream.Close()

	w.Header().Set("Content-Type", "audio/webm")

	io.Copy(w, fileStream)
}
