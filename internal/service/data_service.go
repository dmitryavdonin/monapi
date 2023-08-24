package service

import (
	"monapi/internal/model"
	"monapi/internal/repository"
	"time"
)

type DataService struct {
	repo repository.Data
}

func NewDataService(repo repository.Data) *DataService {
	return &DataService{repo: repo}
}

func (s *DataService) GetData(id int, from time.Time, to time.Time, limit int, offset int) ([]model.Data_new, error) {
	return s.repo.GetData(id, from, to, limit, offset)
}
