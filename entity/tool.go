package entity

import "gorm.io/gorm"

type Tool struct {
	gorm.Model

	Name          string `gorm:"index:idx_tool_name_uniq,unique,where:deleted_at IS NULL"`
	Description   string
	LocalToolName string
}
