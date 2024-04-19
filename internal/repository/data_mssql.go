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

	db_from := fmt.Sprintf("%d-%d-%d %d:%d:%d", from.Year(), from.Day(), from.Month(), from.Hour(), from.Minute(), from.Second())
	db_to := fmt.Sprintf("%d-%d-%d %d:%d:%d", to.Year(), to.Day(), to.Month(), to.Hour(), to.Minute(), to.Second())

	result := r.db.Limit(limit).Offset(offset).Where("ID = ? AND DtWr BETWEEN ? AND ?", id, db_from, db_to).Order("DtWr asc").Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, result.Error
}

// get last value
func (r *DataMsssql) GetLastValue(id int) (model.Data_new, error) {
	var item model.Data_new

	now := time.Now()
	db_now := fmt.Sprintf("%d-%d-%d", now.Year(), now.Day(), now.Month())

	one_day := time.Now().AddDate(0, 0, -1)
	db_one_day := fmt.Sprintf("%d-%d-%d", one_day.Year(), one_day.Day(), one_day.Month())

	seven_days := time.Now().AddDate(0, 0, -7)
	db_seven_days := fmt.Sprintf("%d-%d-%d", seven_days.Year(), seven_days.Day(), seven_days.Month())

	result := r.db.Where("ID = ? AND DtWr >= ?", id, db_now).Order("DtWr desc").First(&item)
	if result.Error != nil {
		result := r.db.Where("ID = ? AND DtWr >= ", id, db_one_day).Order("DtWr desc").First(&item)
		if result.Error != nil {
			result := r.db.Where("ID = ? AND DtWr >= ", id, db_seven_days).Order("DtWr desc").First(&item)
			if result.Error != nil {
				return item, result.Error
			}
		}
	}
	return item, nil
}
