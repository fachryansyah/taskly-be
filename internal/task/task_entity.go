package task

import "time"

type Task struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	UserID    string    `gorm:"index" json:"userId"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	Label     string    `json:"label"`
}
