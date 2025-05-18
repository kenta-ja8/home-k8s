package usecase

import (
	"context"
	"strings"

	"github.com/kenta-ja8/home-k8s-app/pkg/client"
	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
)

type PantryOrderReminderUsecase struct{}

func NewPantryOrderReminderUsecase() *PantryOrderReminderUsecase {
	return &PantryOrderReminderUsecase{}
}

const msg = `
コープの注文は入力しましたか？
from home-k8s
`

func (pu *PantryOrderReminderUsecase) Run(ctx context.Context) error {
	logger.Info("Start pantry order reminder usecase")
	cfg := entity.LoadConfig()
	lineMessageClient := client.NewLineMessageClient(cfg)
	err := lineMessageClient.SendMessage(strings.TrimSpace(msg))
	if err != nil {
		return err
	}
	logger.Info("End pantry order reminder usecase")
	return nil
}
