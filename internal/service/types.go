package service

import (
	"context"
	"sync"
	"time"

	"github.com/anyshake/observer/config"
)

type Status struct {
	mu sync.Mutex

	restarts  int
	isRunning bool
	startedAt time.Time
	stoppedAt time.Time
	updatedAt time.Time
}

type Asset struct {
	FilePath   string
	FileName   string
	Size       int64
	ModifiedAt int64
}

type AssetData struct {
	ContentType string
	FileName    string
	Data        []byte
}

type IService interface {
	GetStatus() *Status
	GetName() string
	GetDescription() string
	Init(ctx context.Context) error
	IsEnabled(ctx context.Context) bool
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Restart(ctx context.Context) error
	GetAssetList() ([]Asset, error)
	GetAssetData(string) (*AssetData, error)
	GetConfigConstraint() []config.IConstraint
}
