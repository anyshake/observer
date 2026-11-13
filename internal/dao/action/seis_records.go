package action

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/sharding"
)

var ErrCleanupRunning = errors.New("cleanup task is already running")

func (h *Handler) SeisRecordsGetQueryWindow() time.Duration {
	return time.Hour
}

// SeisRecordsCreate inserts records and skips those whose record_time already
// exists. Writers retry batches after timeouts whose commit outcome is unknown,
// so a retry must not fail because part of the batch was already stored.
func (h *Handler) SeisRecordsCreate(ctx context.Context, records ...model.SeisRecord) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	groupedRecords := make(map[int][]model.SeisRecord)
	var days []int

	for _, record := range records {
		t := time.UnixMilli(record.RecordTime).UTC().YearDay()
		if _, ok := groupedRecords[t]; !ok {
			days = append(days, t)
		}
		groupedRecords[t] = append(groupedRecords[t], record)
	}

	for _, day := range days {
		group := groupedRecords[day]
		tableName := fmt.Sprintf("%sseis_records_%d", h.daoObj.GetPrefix(), day%model.SEIS_RECORD_SHARDS)

		err := h.daoObj.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			recordTimes := make([]int64, len(group))
			for i, record := range group {
				recordTimes[i] = record.RecordTime
			}
			var existing []int64
			if err := tx.Table(tableName).Where("record_time IN ?", recordTimes).Pluck("record_time", &existing).Error; err != nil {
				return err
			}
			stored := make(map[int64]struct{}, len(existing)+len(group))
			for _, recordTime := range existing {
				stored[recordTime] = struct{}{}
			}
			pending := make([]model.SeisRecord, 0, len(group))
			for _, record := range group {
				if _, ok := stored[record.RecordTime]; ok {
					continue
				}
				stored[record.RecordTime] = struct{}{}
				pending = append(pending, record)
			}
			if len(pending) == 0 {
				return nil
			}
			return tx.Table(tableName).Create(&pending).Error
		})
		if err != nil {
			return fmt.Errorf("failed to insert records for day %d: %w", day, err)
		}
	}

	return nil
}

func (h *Handler) SeisRecordsQuery(ctx context.Context, startTime, endTime time.Time) ([]model.SeisRecord, error) {
	var records []model.SeisRecord
	err := h.SeisRecordsQueryEach(ctx, startTime, endTime, func(record model.SeisRecord) error {
		records = append(records, record)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (h *Handler) SeisRecordsQueryEach(ctx context.Context, startTime, endTime time.Time, callback func(model.SeisRecord) error) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}
	if ctx == nil {
		return errors.New("query context is nil")
	}
	if callback == nil {
		return errors.New("query callback is nil")
	}

	if startTime.After(endTime) {
		return errors.New("start time is after end time")
	}

	queryWindowLimit := h.SeisRecordsGetQueryWindow()
	if endTime.Sub(startTime) > queryWindowLimit {
		return fmt.Errorf("duration between start time and end time exceeds %.0f minutes limit", queryWindowLimit.Minutes())
	}

	startTimeUTC := startTime.UTC()
	endTimeUTC := endTime.UTC()
	currentDate := time.Date(startTimeUTC.Year(), startTimeUTC.Month(), startTimeUTC.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(endTimeUTC.Year(), endTimeUTC.Month(), endTimeUTC.Day(), 0, 0, 0, 0, time.UTC)
	for !currentDate.After(endDate) {
		tableName := fmt.Sprintf("%sseis_records_%d", h.daoObj.GetPrefix(), currentDate.YearDay()%model.SEIS_RECORD_SHARDS)
		if err := h.seisRecordsQueryTable(ctx, tableName, startTime.UnixMilli(), endTime.UnixMilli(), callback); err != nil {
			return err
		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return nil
}

// seisRecordsQueryPageSize bounds how many records are read per statement.
// Every page is fully read and its cursor closed before callbacks run, so slow
// consumers never hold the database connection that SQLite shares with writers.
const seisRecordsQueryPageSize = 60

func (h *Handler) seisRecordsQueryTable(ctx context.Context, tableName string, startTime, endTime int64, callback func(model.SeisRecord) error) error {
	cursor := startTime
	for cursor <= endTime {
		if err := ctx.Err(); err != nil {
			return err
		}
		page, err := h.seisRecordsQueryPage(ctx, tableName, cursor, endTime)
		if err != nil {
			return err
		}
		for _, record := range page {
			if err := ctx.Err(); err != nil {
				return err
			}
			if err := callback(record); err != nil {
				return err
			}
		}
		if len(page) < seisRecordsQueryPageSize {
			return nil
		}
		last := page[len(page)-1].RecordTime
		if last == math.MaxInt64 {
			return nil
		}
		cursor = last + 1
	}
	return nil
}

func (h *Handler) seisRecordsQueryPage(ctx context.Context, tableName string, startTime, endTime int64) ([]model.SeisRecord, error) {
	rows, err := h.daoObj.Database.
		WithContext(ctx).
		Table(tableName).
		Select("id", "created_at", "record_time", "sample_rate", "channel_data").
		Where("record_time >= ? AND record_time <= ?", startTime, endTime).
		Order("record_time ASC").
		Limit(seisRecordsQueryPageSize).
		Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query seismic waveform records in table %s: %w", tableName, err)
	}
	defer rows.Close()

	page := make([]model.SeisRecord, 0, seisRecordsQueryPageSize)
	for rows.Next() {
		var record model.SeisRecord
		if err := rows.Scan(
			&record.PrimaryKey,
			&record.CreatedAt,
			&record.RecordTime,
			&record.SampleRate,
			&record.ChannelData,
		); err != nil {
			return nil, fmt.Errorf("failed to scan seismic waveform record in table %s: %w", tableName, err)
		}
		page = append(page, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate seismic waveform records in table %s: %w", tableName, err)
	}
	return page, nil
}

func (h *Handler) SeisRecordsPurge(ctx context.Context, startTime, endTime time.Time) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	if startTime.After(endTime) {
		return errors.New("start time is after end time")
	}

	return h.CleanupExclusive(ctx, func() error {
		return h.seisRecordsPurgeRange(ctx, startTime, endTime)
	})
}

func (h *Handler) SeisRecordsPurgeAll(ctx context.Context) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	return h.CleanupExclusive(ctx, func() error {
		return h.seisRecordsForEachShard(ctx, func(tableName string) error {
			if err := h.seisRecordsEnsureTable(ctx, tableName); err != nil {
				return fmt.Errorf("failed to ensure seismic waveform records table %s: %w", tableName, err)
			}

			if err := h.seisRecordsDeleteAll(ctx, tableName); err != nil {
				return fmt.Errorf("failed to purge seismic waveform records table %s: %w", tableName, err)
			}

			return nil
		})
	})
}

func (h *Handler) seisRecordsEnsureTable(ctx context.Context, tableName string) error {
	tx := h.daoObj.Database.
		WithContext(ctx).
		Session(&gorm.Session{}).
		Set(sharding.ShardingIgnoreStoreKey, nil).
		Table(tableName)

	if shardingDialector, ok := tx.Dialector.(sharding.ShardingDialector); ok {
		tx.Config.Dialector = shardingDialector.Dialector
	}

	return tx.AutoMigrate(&model.SeisRecord{})
}

func (h *Handler) seisRecordsDeleteAll(ctx context.Context, tableName string) error {
	quotedTableName, err := h.quoteSQLIdentifier(tableName)
	if err != nil {
		return err
	}

	return h.daoObj.Database.
		WithContext(ctx).
		Session(&gorm.Session{}).
		Set(sharding.ShardingIgnoreStoreKey, nil).
		Exec(fmt.Sprintf("DELETE FROM %s", quotedTableName)).
		Error
}

func (h *Handler) seisRecordsPurgeRange(ctx context.Context, startTime, endTime time.Time) error {
	return h.seisRecordsForEachShard(ctx, func(tableName string) error {
		if err := h.seisRecordsEnsureTable(ctx, tableName); err != nil {
			return fmt.Errorf("failed to ensure seismic waveform records table %s: %w", tableName, err)
		}

		if err := h.seisRecordsDeleteRange(ctx, tableName, startTime.UnixMilli(), endTime.UnixMilli()); err != nil {
			return fmt.Errorf("failed to purge seismic waveform records in table %s: %w", tableName, err)
		}

		return nil
	})
}

// seisRecordsPurgeBatchSize bounds each purge statement. A single DELETE of a
// whole day can occupy SQLite's only connection long enough to time out
// concurrent archive inserts.
const seisRecordsPurgeBatchSize = 500

func (h *Handler) seisRecordsDeleteRange(ctx context.Context, tableName string, startTime, endTime int64) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		var ids []uint64
		err := h.daoObj.Database.
			WithContext(ctx).
			Session(&gorm.Session{}).
			Set(sharding.ShardingIgnoreStoreKey, nil).
			Table(tableName).
			Where("record_time >= ? AND record_time <= ?", startTime, endTime).
			Order("record_time ASC").
			Limit(seisRecordsPurgeBatchSize).
			Pluck("id", &ids).
			Error
		if err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		err = h.daoObj.Database.
			WithContext(ctx).
			Session(&gorm.Session{}).
			Set(sharding.ShardingIgnoreStoreKey, nil).
			Table(tableName).
			Where("id IN ?", ids).
			Delete(model.SeisRecord{}).
			Error
		if err != nil {
			return err
		}
		if len(ids) < seisRecordsPurgeBatchSize {
			return nil
		}
	}
}

func (h *Handler) quoteSQLIdentifier(identifier string) (string, error) {
	if identifier == "" {
		return "", errors.New("identifier is empty")
	}

	for _, ch := range identifier {
		if ch >= 'a' && ch <= 'z' {
			continue
		}
		if ch >= 'A' && ch <= 'Z' {
			continue
		}
		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch == '_' {
			continue
		}
		return "", fmt.Errorf("identifier %q contains unsupported character %q", identifier, ch)
	}

	var quoted strings.Builder
	stmt := &gorm.Statement{DB: h.daoObj.Database}
	stmt.QuoteTo(&quoted, clause.Table{Name: identifier})
	return quoted.String(), nil
}

func (h *Handler) seisRecordsForEachShard(ctx context.Context, fn func(tableName string) error) error {
	sem := make(chan struct{}, h.seisRecordsPurgeConcurrency())
	errs := make(chan error, model.SEIS_RECORD_SHARDS)

	var wg sync.WaitGroup
shards:
	for i := 0; i < model.SEIS_RECORD_SHARDS; i++ {
		if ctx.Err() != nil {
			break
		}
		tableName := fmt.Sprintf("%sseis_records_%d", h.daoObj.GetPrefix(), i)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			break shards
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			if err := fn(tableName); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)

	joinedErr := ctx.Err()
	for err := range errs {
		joinedErr = errors.Join(joinedErr, err)
	}
	return joinedErr
}

func (h *Handler) seisRecordsPurgeConcurrency() int {
	if h.daoObj == nil || h.daoObj.Database == nil {
		return 1
	}

	sqlDB, err := h.daoObj.Database.DB()
	if err != nil {
		return 1
	}

	maxOpenConnections := sqlDB.Stats().MaxOpenConnections
	if maxOpenConnections == 1 {
		return 1
	}

	maxConcurrency := min(runtime.NumCPU(), 8)
	if maxOpenConnections > 1 {
		maxConcurrency = min(maxConcurrency, maxOpenConnections)
	}
	if maxConcurrency < 1 {
		return 1
	}

	return maxConcurrency
}

// DatabaseStats reports connection pool statistics for diagnostics.
func (h *Handler) DatabaseStats() (sql.DBStats, bool) {
	if h.daoObj == nil || h.daoObj.Database == nil {
		return sql.DBStats{}, false
	}
	sqlDB, err := h.daoObj.Database.DB()
	if err != nil {
		return sql.DBStats{}, false
	}
	return sqlDB.Stats(), true
}
