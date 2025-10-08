package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kenta-ja8/home-k8s-app/pkg/client"
	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
	"github.com/kenta-ja8/home-k8s-app/pkg/usecase"
)

func run() error {
	cfg := entity.LoadConfig()
	db, err := client.NewPostgresClient(cfg)
	if err != nil {
		return err
	}

	sampleUsecase := usecase.NewSampleUsecase(cfg, db)
	healthcareCollectorUsecase := usecase.NewHealthcareCollectorUsecase(cfg, db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		sampleUsecase.AccessDB(w, r)
	})
	r.Post("/healthcare-collect", func(w http.ResponseWriter, r *http.Request) {
		errorHandler(healthcareCollectorUsecase.Collect(w, r), w)
	})

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}
	return nil
}

func errorHandler(err error, w http.ResponseWriter) {
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

func main() {
	logger.Info("Hello World!")
	defer logger.Info("Goodbye World!")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
