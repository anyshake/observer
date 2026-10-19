package action

import (
	"context"
	"sync"

	"github.com/anyshake/observer/internal/dao"
	"golang.org/x/sync/singleflight"
)

type Handler struct {
	daoObj          *dao.DAO
	cleanupMu       sync.Mutex
	fdsnAccessGroup singleflight.Group
}

func NewHandler(daoObj *dao.DAO) *Handler {
	return &Handler{daoObj: daoObj}
}

func (h *Handler) CleanupExclusive(ctx context.Context, fn func() error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !h.cleanupMu.TryLock() {
		return ErrCleanupRunning
	}
	defer h.cleanupMu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}
	return fn()
}
