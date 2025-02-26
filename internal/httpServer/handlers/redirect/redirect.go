package redirect

import (
	resp "api/internal/httpServer/apiResp"
	"api/internal/storage"
	"context"
	"errors"
	"net/http"

	_ "api/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type Request struct {
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.0 --name=URLGetter
type URLGetter interface {
	GetUrl(alias string, ctx context.Context) (string, error)
}

// GetUrl return list url
// @Summary Redirect URL
// @Description Redirect GET
// @Tag URLRedirect
// @Accept json
// @Produce json
// @Param alias path string true "Alias of the URL to redirect to"
// @Success 200 {object} Response
// @Failure 400 {object} resp.Response
// @Failure 500 {object} resp.Response
// @Failure 404 {object} resp.Response
// @Router /{alias} [get]
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
