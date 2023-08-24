package service

import (
	"monapi/internal/model"
	"monapi/internal/repository"
	"time"
)

type Data interface {
	GetData(id int, from time.Time, to time.Time, limit int, offset int) ([]model.Data_new, error)
}

type Services struct {
	Data
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Data: NewDataService(repos.Data),
	}
}
