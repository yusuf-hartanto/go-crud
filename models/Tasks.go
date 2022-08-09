package Model

import (
	"gorm.io/gorm"
)

type Tasks struct {
	gorm.Model
	ID           int `gorm:"primaryKey;autoIncrement:true"`
	Task         string
	Assignee     string
	Deadline     string
	Description  string
	ProfileImage string
}
