package del

import (
	"context"
	"errors"
	"io"
	"net/http"

	resp "api/internal/httpServer/apiResp"
	"api/internal/storage"

	_ "api/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type Request struct {
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=URLDeletter
type URLDeletter interface {
	DeleteUrl(alias string, ctx context.Context) (int64, error)
}

// DeleteUrl return list url
// @Summary Delete URL
// @Description Delete URL
// @Tag URLDelete
// @Accept json
// @Produce json
// @Param input body Request true "Alias of the URL to delete"
// @Success 200 {object} Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Router /url/urldel/{alias} [delete]
func DelHand(log *zap.Logger, urlDelletter URLDeletter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.DelHand"

		log = log.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", zap.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", zap.Error(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias

		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		id, err := urlDelletter.DeleteUrl(alias, context.Background())
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", zap.Error(err))

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete url", zap.Error(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url delete", zap.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}

}
