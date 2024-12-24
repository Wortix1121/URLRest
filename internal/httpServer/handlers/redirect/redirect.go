package redirect

import (
	resp "Api/internal/httpServer/apiResp"
	"Api/internal/storage"
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=URLGetter
type URLGetter interface {
	GetUrl(alias string, ctx context.Context) (string, error)
}

func GetHand(log *zap.Logger, urlGet URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.GetHand"

		log = log.With(
			zap.String("op", op),
			zap.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGet.GetUrl(alias, context.Background())
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", zap.Error(err))

			render.JSON(w, r, resp.Error("not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", zap.Error(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", zap.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}

}
