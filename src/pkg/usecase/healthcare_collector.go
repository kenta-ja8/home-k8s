package usecase

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kenta-ja8/home-k8s-app/pkg/entity"
	"github.com/kenta-ja8/home-k8s-app/pkg/logger"
	"github.com/kenta-ja8/home-k8s-app/pkg/repository"
	"gorm.io/gorm"
)

type HealthcareCollectorUsecase struct {
	cfg *entity.Config
	db  *gorm.DB
}

func NewHealthcareCollectorUsecase(
	cfg *entity.Config,
	db *gorm.DB,
) *HealthcareCollectorUsecase {
	return &HealthcareCollectorUsecase{
		cfg: cfg,
		db:  db,
	}
}

type RequestBody struct {
	Items []struct {
		Value     string `json:"value"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Type      string `json:"type"`
		Unit      string `json:"unit"`
		Duration  string `json:"duration"`
		Source    string `json:"source"`
		Name      string `json:"name"`
	} `json:"items"`
}

func (u *HealthcareCollectorUsecase) Collect(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		logger.Error("HealthcareCollectorUsecase Collect ReadAll", "err: ", err)
		return err
	}
	var body RequestBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		logger.Error("HealthcareCollectorUsecase Collect Unmarshal", "err: ", err)
		return err
	}
	defer func() { _ = r.Body.Close() }()

	for _, item := range body.Items {
		logger.Info("HealthcareCollectorUsecase Collect", "item:", item)

		type HeartRate struct {
			repository.BaseModel
			Value     float64
			StartDate time.Time
		}
		heartRate, err := gorm.G[HeartRate](u.db).Where("start_date = ?", item.StartDate).First(ctx)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("HealthcareCollectorUsecase Collect GORM First", "err: ", err)
			continue
		}
		if err == gorm.ErrRecordNotFound {
			logger.Info("create", "heartRate:", heartRate)
			startDate, err := time.Parse(time.RFC3339, item.StartDate)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect time Parse", "err: ", err)
				continue
			}
			value, err := strconv.ParseFloat(item.Value, 64)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect strconv ParseFloat", "err: ", err)
				continue
			}
			heartRate := HeartRate{
				BaseModel: repository.NewBaseModel(),
				Value:     value,
				StartDate: startDate,
			}
			err = gorm.G[HeartRate](u.db).Create(ctx, &heartRate)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect GORM Create", "err: ", err)
			}
			continue
		}
		logger.Info("skip ", "heartRate:", heartRate)
	}

	_, _ = w.Write([]byte("ok"))

	return nil
}
