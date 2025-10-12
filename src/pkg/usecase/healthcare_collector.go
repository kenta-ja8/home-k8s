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
	"github.com/pkg/errors"
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

type valueConverter func(string) (string, error)

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

func identityConverter(val string) (string, error) {
	return val, nil
}

func buildHealthRecord(item Item, convert valueConverter) (*HealthRecord, error) {
	if convert == nil {
		convert = identityConverter
	}
	startTime, err := time.Parse(time.RFC3339, item.StartTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse start time")
	}
	endTime, err := time.Parse(time.RFC3339, item.EndTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse end time")
	}
	durationSeconds, err := convertDurationToSeconds(item.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse duration")
	}
	value, err := convert(item.Value)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert value")
	}
	return &HealthRecord{
		BaseModel:       repository.NewBaseModel(),
		Value:           value,
		StartTime:       startTime,
		EndTime:         endTime,
		TypeCode:        item.Type,
		Unit:            item.Unit,
		Duration:        item.Duration,
		DurationSeconds: durationSeconds,
		Source:          item.Source,
		Name:            item.Name,
	}, nil
}

func (u *HealthcareCollectorUsecase) Collect(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.WithStack(err)
	}
	var body RequestBody
	err = json.Unmarshal(b, &body)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() { _ = r.Body.Close() }()

	if len(body.ActiveEnergyItems) > 0 {
		if err := u.save(ctx, body.ActiveEnergyItems, "h_active_energies", nil); err != nil {
			return err
		}
	}
	if len(body.RestingEnergyItems) > 0 {
		if err := u.save(ctx, body.RestingEnergyItems, "h_resting_energies", nil); err != nil {
			return err
		}
	}
	if len(body.SleepItems) > 0 {
		if err := u.save(ctx, body.SleepItems, "h_sleep_records", nil); err != nil {
			return err
		}
	}
	if len(body.MindfulMinuteItems) > 0 {
		if err := u.save(ctx, body.MindfulMinuteItems, "h_mindful_minutes", convertStringToFloat); err != nil {
			return err
		}
	}
	if len(body.HeartRateItems) > 0 {
		if err := u.save(ctx, body.HeartRateItems, "h_heart_rates", nil); err != nil {
			return err
		}
	}

	_, _ = w.Write([]byte("ok"))

	return nil
}

func (u *HealthcareCollectorUsecase) save(ctx context.Context, items []Item, tableName string, convert valueConverter) error {
	converter := convert
	if converter == nil {
		converter = identityConverter
	}
	tx := u.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	for _, item := range items {
		logger.Info("HealthcareCollectorUsecase Collect item: %+v", item)

		record, err := buildHealthRecord(item, converter)
		if err != nil {
			return err
		}

		query := tx.WithContext(ctx).
			Table(tableName).
			Where("start_time = ? AND value = ? AND source = ?", record.StartTime, record.Value, record.Source)

		var existing HealthRecord
		if err := query.Take(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.Wrap(err, "failed to query")
			}

			if err := tx.WithContext(ctx).Table(tableName).Create(record).Error; err != nil {
				return errors.Wrap(err, "failed to query")
			}

			logger.Info("created healthcare record: %+v", record)
			continue
		}

		logger.Info("record already existed: %+v", &existing)
	}

	tx.Commit()

	return nil
}
