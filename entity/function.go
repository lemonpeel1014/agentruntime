package entity

import "gorm.io/gorm"

type Function struct {
	gorm.Model

	Name        string `gorm:"index:idx_function_name_uniq,unique,where:deleted_at IS NULL"`
	Description string
}
