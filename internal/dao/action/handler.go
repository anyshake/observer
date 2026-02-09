package action

import "github.com/anyshake/observer/internal/dao"

type Handler struct {
	daoObj *dao.DAO
}

func NewHandler(daoObj *dao.DAO) *Handler {
	return &Handler{daoObj}
}
