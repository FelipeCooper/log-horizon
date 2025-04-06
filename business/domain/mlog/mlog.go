package mlog

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/felipecooper/log-horizon/foundation/logger"
	"github.com/felipecooper/log-horizon/foundation/transaction"
	"github.com/oklog/ulid/v2"
)

var (
	ErrOnRegisterLog    = errors.New("failed on save log in writer")
	ErrInvalidLevel     = errors.New("unrecognized level")
	ErrInvalidTimeRange = errors.New("invalid time range")
)

type Business struct {
	logger logger.Logger
	store  Store
}

func NewMlog(logger logger.Logger, store Store) *Business {
	return &Business{
		logger: logger,
		store:  store,
	}
}

func (b *Business) NewWithTx(tx transaction.CommitRollbacker) (*Business, error) {
	return b, nil
}

func (b *Business) Register(ctx context.Context, message string, level Level, metadata map[string]string) (Log, error) {
	if !level.IsValid() {
		b.logger.Error(ctx, fmt.Sprintf("unrecognized level: %s", level), "error", ErrInvalidLevel)
		return Log{}, fmt.Errorf("register: %w", ErrInvalidLevel)
	}

	log := Log{
		ID:        ulid.Make(),
		Message:   message,
		Timestamp: time.Now(),
		Level:     level,
		Metadata:  metadata,
	}

	err := b.store.Write(ctx, &log)
	if err != nil {
		b.logger.Error(ctx, "failed to register log", "error", err)
		return Log{}, fmt.Errorf("register: %w", ErrOnRegisterLog)
	}

	return log, nil
}

func (b *Business) Query(ctx context.Context, startTime, endTime time.Time, level Level, page, pageSize int) (SearchResult, error) {
	if !startTime.IsZero() && !endTime.IsZero() && endTime.Before(startTime) {
		b.logger.Error(ctx, "end time before start time", "error", ErrInvalidTimeRange)
		return SearchResult{}, fmt.Errorf("query: %w", ErrInvalidTimeRange)
	}

	if level != "" && !level.IsValid() {
		b.logger.Error(ctx, fmt.Sprintf("invalid level: %s", level), "error", ErrInvalidLevel)
		return SearchResult{}, fmt.Errorf("query: %w", ErrInvalidLevel)
	}

	criteria := SearchCriteria{
		TimeRange: TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Level:    level,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := b.store.Search(ctx, criteria)
	if err != nil {
		b.logger.Error(ctx, "failed to search logs", "error", err)
		return SearchResult{}, fmt.Errorf("query: %w", err)
	}

	return result, nil
}

func (b *Business) ExportToFile(ctx context.Context, startTime, endTime time.Time, level Level) (string, int64, error) {
	if !startTime.IsZero() && !endTime.IsZero() && endTime.Before(startTime) {
		b.logger.Error(ctx, "end time before start time", "error", ErrInvalidTimeRange)
		return "", 0, fmt.Errorf("export: %w", ErrInvalidTimeRange)
	}

	if level != "" && !level.IsValid() {
		b.logger.Error(ctx, fmt.Sprintf("invalid level: %s", level), "error", ErrInvalidLevel)
		return "", 0, fmt.Errorf("export: %w", ErrInvalidLevel)
	}

	criteria := SearchCriteria{
		TimeRange: TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Level: level,
	}

	fileURL, fileSize, err := b.store.ExportToFile(ctx, criteria)
	if err != nil {
		b.logger.Error(ctx, "failed to export logs to file", "error", err)
		return "", 0, fmt.Errorf("export: %w", err)
	}

	return fileURL, fileSize, nil
}

func (b *Business) Count(ctx context.Context, startTime, endTime time.Time, level Level) (int, error) {
	if !startTime.IsZero() && !endTime.IsZero() && endTime.Before(startTime) {
		b.logger.Error(ctx, "end time before start time", "error", ErrInvalidTimeRange)
		return 0, fmt.Errorf("count: %w", ErrInvalidTimeRange)
	}

	if level != "" && !level.IsValid() {
		b.logger.Error(ctx, fmt.Sprintf("invalid level: %s", level), "error", ErrInvalidLevel)
		return 0, fmt.Errorf("count: %w", ErrInvalidLevel)
	}

	criteria := SearchCriteria{
		TimeRange: TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Level: level,
	}

	count, err := b.store.Count(ctx, criteria)
	if err != nil {
		b.logger.Error(ctx, "failed to count logs", "error", err)
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}
