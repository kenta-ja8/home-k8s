package helper

import (
	"github.com/google/uuid"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
)

func NewUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		logger.Error("Failed to generate UUID")
		panic(err)
	}

	return id.String()
}
