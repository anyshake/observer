package model

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/anyshake/observer/internal/dao"
	"github.com/anyshake/observer/internal/hardware/explorer"
	"gorm.io/gorm"
	"gorm.io/sharding"
)

const SEIS_RECORD_SHARDS = 366 // By days of the year

type SeisRecord struct {
	dao.BaseTable
	RecordTime  int64  `gorm:"column:record_time;index;not null;unique"`
	SampleRate  int    `gorm:"column:sample_rate;not null"`
	ChannelData []byte `gorm:"column:channel_data;not null"`
}

func (t *SeisRecord) GetModel() any {
	return &SeisRecord{}
}

func (t *SeisRecord) GetName(tablePrefix string) string {
	return fmt.Sprintf("%s%s", tablePrefix, "seis_records")
}

func (t *SeisRecord) UseAutoMigrate() bool {
	return false
}

func (t *SeisRecord) AddPlugins(dbObj *gorm.DB, tablePrefix string) ([]gorm.Plugin, error) {
	tableName := t.GetName(tablePrefix)
	for i := 0; i < SEIS_RECORD_SHARDS; i++ {
		if err := dbObj.Table(fmt.Sprintf("%s_%d", tableName, i)).AutoMigrate(&SeisRecord{}); err != nil {
			return nil, fmt.Errorf("failed to auto migrate table %s_%d: %w", tableName, i, err)
		}
	}

	shard := sharding.Register(sharding.Config{
		ShardingKey:    "record_time",
		NumberOfShards: SEIS_RECORD_SHARDS,
		ShardingAlgorithm: func(columnValue any) (string, error) {
			return "", nil // We will handle sharding ourselves in SeisRecordsCreate
		},
		PrimaryKeyGenerator: sharding.PKSnowflake,
	}, tableName)

	return []gorm.Plugin{shard}, nil
}

func (t *SeisRecord) Encode(recordTime time.Time, sampleRate int, channelData []explorer.ChannelData) error {
	t.RecordTime = recordTime.UnixMilli()
	t.SampleRate = sampleRate

	buf := bytes.Buffer{}
	encoder := gob.NewEncoder(&buf)

	if err := encoder.Encode(channelData); err != nil {
		return fmt.Errorf("failed to encode channel data: %w", err)
	}
	t.ChannelData = buf.Bytes()

	return nil
}

func (t *SeisRecord) Decode() (recordTime time.Time, sampleRate int, channelData []explorer.ChannelData, err error) {
	buf := bytes.Buffer{}
	decoder := gob.NewDecoder(&buf)

	if _, err = buf.Write(t.ChannelData); err != nil {
		return time.Time{}, 0, nil, fmt.Errorf("failed to write channel data: %w", err)
	}
	if err = decoder.Decode(&channelData); err != nil {
		return time.Time{}, 0, nil, fmt.Errorf("failed to decode channel data: %w", err)
	}

	return time.UnixMilli(t.RecordTime), t.SampleRate, channelData, nil
}
