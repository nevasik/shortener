package save

import (
	"Url-shortner/internal/config"
	resp "Url-shortner/internal/lib/api/response"
	"Url-shortner/internal/lib/logger/slogger"
	"Url-shortner/internal/lib/random"
	"Url-shortner/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

type Request struct { // запрос
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"` // omitempty - означает, что если будет null, то оно не будет сериализоваться
}

type Response struct { // ответ
	resp.Response
	Alias string `json:"alias,omitempty"` // на случай, если пользователь не указал alias, то он будет генерироваться самостоятельно случ. образом
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc { // обработчик сохранения урла
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request body"))
			return
		}

		if err != nil {
			log.Error("failed to decode request body", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoder", slog.Any("request", req))

		// проверяем ошибки на валидацию для возврата человекочитаемых ошибок
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", slogger.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(config.AliasLength) // генератор случайных букв
			// здесь может быть такая ситуация, что сгенерированный alias уже используется нами в другом url, подумать как это можно обработать в будущем
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrUrlExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", slogger.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.Ok(), Alias: alias,
	})
}
