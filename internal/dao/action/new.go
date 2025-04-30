package action

import "github.com/anyshake/observer/internal/dao"

func New(daoObj *dao.DAO) *Handler {
	return &Handler{daoObj}
}
