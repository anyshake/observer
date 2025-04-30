package action

import (
	"errors"
	"fmt"
	"time"

	"github.com/anyshake/observer/internal/dao/model"
)

func (h *Handler) SeisRecordsCreate(records ...model.SeisRecord) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	groupedRecords := make(map[int][]model.SeisRecord)

	for _, record := range records {
		t := time.UnixMilli(record.Timestamp).UTC().YearDay()
		groupedRecords[t] = append(groupedRecords[t], record)
	}

	for day, group := range groupedRecords {
		err := h.daoObj.Database.
			Table((&model.SeisRecord{}).GetName(h.daoObj.GetPrefix())).
			Create(&group).
			Error
		if err != nil {
			return fmt.Errorf("failed to insert records for day %d: %w", day, err)
		}
	}

	return nil
}

func (h *Handler) SeisRecordsQuery(startTime, endTime time.Time) ([]model.SeisRecord, error) {
	if h.daoObj == nil {
		return nil, errors.New("database is not opened")
	}

	if startTime.After(endTime) {
		return nil, errors.New("start time is after end time")
	}

	if endTime.Sub(startTime) > time.Hour {
		return nil, errors.New("duration between start time and end time exceeds 1 hour limit")
	}

	var records []model.SeisRecord
	for currentDay := startTime.UTC().YearDay(); currentDay <= endTime.UTC().YearDay(); currentDay++ {
		tableName := fmt.Sprintf("%sseis_records_%d", h.daoObj.GetPrefix(), currentDay)

		var tempRecords []model.SeisRecord
		err := h.daoObj.Database.
			Table(tableName).
			Where("timestamp >= ? AND timestamp <= ?", startTime.UnixMilli(), endTime.UnixMilli()).
			Order("timestamp ASC").
			Find(&tempRecords).
			Error
		if err != nil {
			return nil, fmt.Errorf("failed to query seismic waveform records in table %s: %w", tableName, err)
		}

		records = append(records, tempRecords...)
	}

	return records, nil
}

func (h *Handler) SeisRecordsPurge(startTime, endTime time.Time) error {
	if h.daoObj == nil {
		return errors.New("database is not opened")
	}

	for i := 0; i < model.SEIS_RECORD_SHARDS; i++ {
		tableName := fmt.Sprintf("%sseis_records_%d", h.daoObj.GetPrefix(), i)

		err := h.daoObj.Database.
			Table(tableName).
			Where("timestamp >= ? AND timestamp <= ?", startTime.UnixMilli(), endTime.UnixMilli()).
			Delete(model.SeisRecord{}).
			Error
		if err != nil {
			return fmt.Errorf("failed to purge seismic waveform records in table %s: %w", tableName, err)
		}
	}

	return nil
}
