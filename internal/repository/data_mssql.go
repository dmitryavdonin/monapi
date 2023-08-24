package repository

import (
	"fmt"
	"monapi/internal/model"
	"time"

	"gorm.io/gorm"
)

type DataMsssql struct {
	db *gorm.DB
}

func NewDataMssql(db *gorm.DB) *DataMsssql {
	return &DataMsssql{db: db}
}

func (r *DataMsssql) GetData(id int, from time.Time, to time.Time, limit int, offset int) ([]model.Data_new, error) {
	var items []model.Data_new

	db_from := fmt.Sprintf("%d-%d-%d", from.Year(), from.Day(), from.Month())
	db_to := fmt.Sprintf("%d-%d-%d", to.Year(), to.Day(), to.Month())

	result := r.db.Limit(limit).Offset(offset).Where("ID = ? AND DtWr BETWEEN ? AND ?", id, db_from, db_to).Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}
