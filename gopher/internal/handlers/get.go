package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type RequestChecker interface {
	IsRequestReady(context.Context, int64) (bool, error)
}

func GetByRequestID(logger *slog.Logger, requestChecker RequestChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetByRequestID"

		log := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		rawID := chi.URLParam(r, "id")
		if rawID == "" {
			log.Error("%s: %w", op, fmt.Errorf("request id was empty"))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("request id should not be empty"))
			return
		}

		requestID, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			log.Error("%s: %w", op, err)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("request id is not an int"))
			return
		}

		isReady, err := requestChecker.IsRequestReady(r.Context(), requestID)
		if err != nil {
			log.Error("%s: %w", op, err)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, newError("no such request"))
			return
		}

		if isReady {
			log.Info("request status is checked, request is proccessed")
			render.Status(r, http.StatusOK)
			render.JSON(w, r, newError("request is ready"))
		} else {
			log.Info("request status is checked, request is not proccessed")
			render.Status(r, http.StatusOK)
			render.JSON(w, r, newError("request is not ready"))
		}
	}
}
