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

func exec() error {
	cfg := entity.LoadConfig()
	db, err := client.NewPostgresClient(cfg)
	if err != nil {
		return err
	}

	sampleUsecase := usecase.NewSampleUsecase(cfg, db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		sampleUsecase.AccessDB(w, r)
	})

	var employee entity.Employee
	db.First(&employee)
	fmt.Println("First employee:", employee, employee.ID)
	fmt.Printf("First employee: %+v\n", employee)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	logger.Info("Hello World!")
	defer logger.Info("Goodbye World!")

	if err := exec(); err != nil {
		log.Fatal(err)
	}
}
