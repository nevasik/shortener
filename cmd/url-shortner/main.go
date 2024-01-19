package main

import (
	"Url-shortner/internal/config"
	"Url-shortner/internal/http-server/handlers/url/save"
	"Url-shortner/internal/lib/logger/slogger"
	"Url-shortner/internal/storage/postgresql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// init config: cleanenv
	cnf := config.MustLoad()
	log := config.SetupLogger()

	log.Info(
		"starting url-shortener",
		slog.String("env", cnf.Env),
		slog.String("version", "123"),
	)

	// init storage: PostgresSQL
	storage, err := postgresql.New()

	if err != nil {
		log.Error("failed to init storage", slogger.Err(err))
		os.Exit(1)
	}

	// init router: chi, chi render
	router := chi.NewRouter()

	router.Use(middleware.Logger)    // логгирование запросов, время выполнения и ошибок
	router.Use(middleware.RequestID) // идентификатор запроса
	router.Use(middleware.RealIP)    // ip адрес запроса
	router.Use(middleware.Recoverer) // на случай, если случится panic, мы смогли бы его обработать и чтобы приложение не падалось
	router.Use(middleware.URLFormat) // чтобы в приложении были красивые URL

	router.Post("/url", save.New(log, storage))

	fmt.Println()
	log.Info("starting server", slog.String("address", cnf.Address))

	// run server
	server := &http.Server{
		Addr:         cnf.Address,
		Handler:      router,
		ReadTimeout:  cnf.TimeOut, // время, за которое мы читаем запросы с сервера
		WriteTimeout: cnf.TimeOut, // время, за которое мы отправляем ответ на сервер
		IdleTimeout:  cnf.IdleTimeout,
	}

	err = server.ListenAndServe() // сам запуск сервера
	if err != nil {
		log.Error("Failed to start server")
	}

	log.Error("Server stopped")
}
