package models

import (
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Hash   string
	Hash75 string
	Hash50 string
	Hash25 string
	Ext    string
	Name   string
}
