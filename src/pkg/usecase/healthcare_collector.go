package usecase

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
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

type Item struct {
	Value     string `json:"value"`
	StartTime string `json:"startDate"`
	EndTime   string `json:"endDate"`
	Type      string `json:"type"`
	Unit      string `json:"unit"`
	Duration  string `json:"duration"`
	Source    string `json:"source"`
	Name      string `json:"name"`
}
type RequestBody struct {
	RestingEnergyItems []Item `json:"restingEnergyItems"`
	ActiveEnergyItems  []Item `json:"activeEnergyItems"`
	SleepItems         []Item `json:"sleepItems"`
	MindfulMinuteItems []Item `json:"mindfulMinuteItems"`
	HeartRateItems     []Item `json:"heartRateItems"`
}

func convertDurationToSeconds(val string) (string, error) {
	if val == "" {
		return "0", nil
	}

	// hh:mm:ss
	parts := strings.Split(val, ":")

	seconds, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return "", err
	}
	minutes := 0
	if len(parts) > 1 {
		minutes, err = strconv.Atoi(parts[len(parts)-2])
		if err != nil {
			return "", err
		}
	}
	hours := 0
	if len(parts) > 2 {
		hours, err = strconv.Atoi(parts[len(parts)-3])
		if err != nil {
			return "", err
		}
	}

	totalSecontds := hours*60*60 + minutes*60 + seconds
	return strconv.Itoa(totalSecontds), nil
}

func convertStringToFloat(val string) (string, error) {
	if val == "" {
		return "0", nil
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(f, 'f', 1, 64), nil
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

	if len(body.ActiveEnergyItems) > 0 {
		err = u.save(ctx, body.ActiveEnergyItems, "h_active_energies", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			logger.Error("HealthcareCollectorUsecase Collect save healthcare_active_energy", "err: ", err)
			return err
		}
	}
	if len(body.RestingEnergyItems) > 0 {
		err = u.save(ctx, body.RestingEnergyItems, "h_resting_energies", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			logger.Error("HealthcareCollectorUsecase Collect save healthcare_resting_energy", "err: ", err)
			return err
		}
	}
	if len(body.SleepItems) > 0 {
		err = u.save(ctx, body.SleepItems, "h_sleep_records", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			logger.Error("HealthcareCollectorUsecase Collect save healthcare_resting_energy", "err: ", err)
			return err
		}
	}
	if len(body.MindfulMinuteItems) > 0 {
		err = u.save(ctx, body.MindfulMinuteItems, "h_mindful_minutes", convertStringToFloat)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			logger.Error("HealthcareCollectorUsecase Collect save healthcare_resting_energy", "err: ", err)
			return err
		}
	}
	if len(body.HeartRateItems) > 0 {
		err = u.save(ctx, body.HeartRateItems, "h_heart_rates", nil)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			logger.Error("HealthcareCollectorUsecase Collect save healthcare_resting_energy", "err: ", err)
			return err
		}
	}

	_, _ = w.Write([]byte("ok"))

	return nil
}

func (u *HealthcareCollectorUsecase) save(ctx context.Context, items []Item, tableName string, convert func(string) (string, error)) error {
	for _, item := range items {
		logger.Info("HealthcareCollectorUsecase Collect", "item:", item)

		type HealthRecord struct {
			repository.BaseModel
			Value           string
			StartTime       time.Time
			EndTime         time.Time
			TypeCode        string
			Unit            string
			Duration        string
			DurationSeconds string
			Source          string
			Name            string
		}
		heartRate, err := gorm.G[HealthRecord](u.db).Table(tableName).Where("start_time = ?", item.StartTime).First(ctx)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.Error("HealthcareCollectorUsecase Collect GORM First", "err: ", err)
			continue
		}
		if err == gorm.ErrRecordNotFound {
			logger.Info("create", "heartRate:", heartRate)
			startTime, err := time.Parse(time.RFC3339, item.StartTime)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect time Parse", "err: ", err)
				continue
			}
			endTime, err := time.Parse(time.RFC3339, item.EndTime)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect time Parse", "err: ", err)
				continue
			}
			DurationSeconds, err := convertDurationToSeconds(item.Duration)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect Atoi", "err: ", err)
				continue
			}
			value := item.Value
			if convert != nil {
				value, err = convert(item.Value)
				if err != nil {
					logger.Error("HealthcareCollectorUsecase Collect convert", "err: ", err)
					continue
				}
				logger.Info("converted value", "value:", item.Value, "->", value)
			}

			heartRate := HealthRecord{
				BaseModel:       repository.NewBaseModel(),
				Value:           value,
				StartTime:       startTime,
				EndTime:         endTime,
				TypeCode:        item.Type,
				Unit:            item.Unit,
				Duration:        item.Duration,
				DurationSeconds: DurationSeconds,
				Source:          item.Source,
				Name:            item.Name,
			}
			err = gorm.G[HealthRecord](u.db).Table(tableName).Create(ctx, &heartRate)
			if err != nil {
				logger.Error("HealthcareCollectorUsecase Collect GORM Create", "err: ", err)
			}
			continue
		}
		logger.Info("skip ", "heartRate:", heartRate)
	}
	return nil
}
