package service

import (
	"gin-items/dao"
)

type Service struct {
	dao  *dao.Dao
}

// New init
func New() (serv *Service) {
	serv = &Service{
		dao: dao.New(),
	}
	return serv
}
