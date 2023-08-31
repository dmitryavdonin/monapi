package model

import (
	"time"
)

type Data_new struct {
	ID          int       `gorm:"type:integer" json:"id"`
	Temperature float32   `gorm:"type:numeric" json:"temperature"`
	Humidity    float32   `gorm:"type:numeric" json:"humidity"`
	DtWr        time.Time `gorm:"not null" json:"dt_wr" time_format:"2006-01-02 15:04:05"`
}

type DTO struct {
	Temperature float32 `json:"T"`
	Humidity    float32 `json:"H"`
	Time        string  `json:"t"`
}
