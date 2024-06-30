package model

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CloudSymlinkSync struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Name       string         `json:"name"`
	CloudPath  string         `json:"cloudPath"`
	LocalPath  string         `json:"localPath"`
	Exclude    pq.StringArray `json:"exclude" gorm:"type:text[]"`
	IsCD2      bool           `json:"isCD2"`
	IsTimeScan bool           `json:"isTimeScan"`
}
