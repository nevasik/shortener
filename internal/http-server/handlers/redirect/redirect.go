package redirect

import (
	resp "Url-shortner/internal/lib/api/response"
	"Url-shortner/internal/lib/logger/slogger"
	"Url-shortner/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct { // запрос
	Alias string `json:"alias" validate:"required,alias"`
}

type Response struct { // ответ
	resp.Response
	URL string `json:"URL"`
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlToGet URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req) // получается тут в req будет хранится результат ответа и если там будет ошибка в err запишется она

		log.Info("request body decoder", slog.Any("request", req))

		alias := r.RequestURI[1:] // убираем /

		url, err := urlToGet.GetURL(alias) // вызов самого метода
		if errors.Is(err, storage.ErrUrlNotFound) {
			log.Info("URL not found for this alias", slog.String("alias", req.Alias))

			render.JSON(w, r, resp.Error("URL not found for this alias"))
			return
		}

		if err != nil {
			log.Error("failed to get URL", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to get URL"))
			return
		}

		log.Info("URL found!", slog.String("URL", url))

		http.Redirect(w, r, url, http.StatusOK)
	}
}

//func responseOK(w http.ResponseWriter, r *http.Request, url string) {
//	render.JSON(w, r, Response{
//		Response: resp.Ok(),
//		URL:      url,
//	})
//}
