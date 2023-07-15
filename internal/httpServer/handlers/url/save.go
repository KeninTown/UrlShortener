package save

import (
	"errors"
	"fmt"
	"net/http"
	res "urlShortener/internal/lib/api"
	"urlShortener/internal/lib/loggers/sl"
	"urlShortener/internal/lib/random"
	"urlShortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	"gopkg.in/go-playground/validator.v9"
)

// TODO: move aliasLenght to config
const aliasLenght int = 4

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	res.Response
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.url.save.New()"

		log = log.With(
			slog.String("op", op),
			slog.String("requestId", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request", sl.Error(err))

			render.JSON(w, r, res.Error("failed to decode request"))

			return
		}

		log.Info("reqest body decoded", slog.Any("req", req))

		if err = validator.New().Struct(req); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("validation error ocures", sl.Error(err))

			render.JSON(w, r, res.ValidationError(validateError))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLenght)
		}

		// TODO: обработать пересечение alias
		id, err := urlSaver.SaveUrl(req.Url, alias)

		if err != nil {
			if errors.Is(err, storage.ErrAliasExist) {
				log.Info("alias already exist")

				render.JSON(w, r, res.Error(""))

				return
			}

			render.JSON(w, r, res.Error("failed to save url"))
			return
		}

		log.Info(fmt.Sprintf("Save url %s with alias %s", req.Url, req.Url), slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: res.OK(),
			Alias:    alias,
		})
	}
}
