package import_logs

import (
	"time"
)

type ImportLogs struct {
	Id         int64     `gorm:"primary_key" json:"id"`
	FileHash   string    `json:"file_hash"`
	Name       string    `json:"name"`
	Mobile     string    `json:"mobile"`
	CardNo     string    `json:"card_no"`
	Status     int       `json:"status"`
	LogContent string    `json:"log_content"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
