package usecase

import (
	"fmt"
	"net/http"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/helper"
	"gorm.io/gorm"
)

type SampleUsecase struct {
	cfg *entity.Config
	db  *gorm.DB
}

func NewSampleUsecase(cfg *entity.Config, db *gorm.DB) *SampleUsecase {
	return &SampleUsecase{
		cfg: cfg,
		db:  db,
	}
}

func (u *SampleUsecase) AccessDB(w http.ResponseWriter, r *http.Request) {
	var employee entity.Employee
	u.db.First(&employee)
	fmt.Println("First employee:", employee, employee.ID)
	fmt.Printf("First employee: %+v\n", employee)
	w.Write([]byte(
		"AAAwelcomeX" + u.cfg.BUILD_DATE + employee.Name,
	))
}

func (u *SampleUsecase) InsertRecord() error {
	tx := u.db.Begin()
	defer tx.Rollback()

	emp := entity.Employee{
		ID:   helper.NewUUID(),
		Name: "Tanaka",
	}
	u.db.Create(&emp)
	tx.Commit()
	return nil
}
