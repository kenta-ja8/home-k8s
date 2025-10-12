package client

import (
	"fmt"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgresClient(cfg *entity.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable&search_path=main",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_DB,
	)
	if cfg.IS_LOCAL {
		dsn = fmt.Sprintf(
			"postgres://%s@%s:15432/%s?sslmode=disable&search_path=main",
			cfg.POSTGRES_USER,
			cfg.POSTGRES_HOST,
			cfg.POSTGRES_DB,
		)
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open database")
	}
	return db, nil
}
