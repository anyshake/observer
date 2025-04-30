package service

import (
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
	FileName string
	Data     []byte
}

type IService interface {
	GetStatus() *Status
	GetName() string
	GetDescription() string
	Init() error
	IsEnabled() bool
	Start() error
	Stop() error
	Restart() error
	GetAssetList() ([]Asset, error)
	GetAssetData(string) (*AssetData, error)
	GetConfigConstraint() []config.IConstraint
}
