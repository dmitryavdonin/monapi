package model

import (
	"time"
)

type Data_new struct {
	ID          int       `gorm:"type:integer" json:"id"`
	Temperature float32   `gorm:"type:numeric" json:"temperature"`
	Humidity    float32   `gorm:"type:numeric" json:"humidity"`
	DtWr        time.Time `gorm:"not null" json:"dt_wr"`
}
