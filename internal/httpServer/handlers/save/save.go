package save

import (
	"Api/internal/config"
	resp "Api/internal/httpServer/apiResp"
	random "Api/internal/lib/random"
	"Api/internal/storage"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

const aliaslength = 6

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=URLSaver
type URLSaver interface {
	SaveUrl(ctx context.Context, alias string, urlToSave string) (int64, error)
}

func SaveHand(log *zap.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.SaveHand"

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

		aL := config.MustLoad()

		//Добавить - обработку ошибки если пришёл уже существующий alias
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aL.AliasLength)
		}

		id, err := urlSaver.SaveUrl(context.Background(), alias, req.URL)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", zap.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", zap.Error(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", zap.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})

	}
}
