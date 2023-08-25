package repository

import (
	"monapi/internal/model"
	"time"

	"gorm.io/gorm"
)

type Data interface {
	GetData(id int, from time.Time, to time.Time, limit int, offset int) ([]model.Data_new, error)
	GetLastValue(id int) (model.Data_new, error)
}

type Repository struct {
	Data
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Data: NewDataMssql(db),
	}
}
