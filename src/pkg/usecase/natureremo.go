package usecase

import (
	"context"
	"time"

	"github.com/kenta-ja8/home-k8s-app/pkg/client"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
	"github.com/kenta-ja8/home-k8s-app/pkg/repository"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Natureremo struct {
	natureremoClient *client.NatureremoClient
	db               *gorm.DB
}

func NewNatureremo(
	natureremoClient *client.NatureremoClient,
	db *gorm.DB,
) *Natureremo {
	return &Natureremo{
		natureremoClient: natureremoClient,
		db:               db,
	}
}

type NatureremoEvent struct {
	repository.BaseModel
	DeviceID        string
	DeviceName      string
	EventType       string
	SensorCreatedAt time.Time
	Value           float64
}

func (u *Natureremo) Collect(ctx context.Context) error {
	devices, err := u.natureremoClient.GetDevices(ctx)
	if err != nil {
		return err
	}
	tx := u.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	var tableName = "natureremo_events"
	for _, device := range devices {
		for eventType, event := range device.NewestEvents {
			query := tx.WithContext(ctx).
				Table(tableName).
				Where(
					"device_id = ? AND event_type = ? AND sensor_created_at = ?",
					device.ID,
					eventType,
					event.CreatedAt,
				)

			var existing NatureremoEvent
			if err := query.Take(&existing).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.Wrap(err, "failed to query")
				}

				record := &NatureremoEvent{
					BaseModel:       repository.NewBaseModel(),
					DeviceID:        device.ID,
					DeviceName:      device.Name,
					EventType:       eventType,
					SensorCreatedAt: event.CreatedAt,
					Value:           event.Val,
				}
				if err := tx.WithContext(ctx).Table(tableName).Create(record).Error; err != nil {
					return errors.Wrap(err, "failed to created")
				}

				logger.Info("created natureremo event record: %+v", record)
			}
			logger.Info("skip existing natureremo event: %+v", &existing)
		}
	}

	tx.Commit()

	return nil
}
