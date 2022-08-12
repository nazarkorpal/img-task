package models

import "time"

type Image struct {
	Hash      string `gorm:"primaryKey"`
	Hash75    string
	Hash50    string
	Hash25    string
	Ext       string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
