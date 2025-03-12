package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MessageExample struct {
	User string `json:"user"`
	Text string `json:"text"`
}

type Agent struct {
	gorm.Model

	Name            string `gorm:"index:idx_agent_name_uniq,unique,where:deleted_at IS NULL"`
	ModelName       string
	System          string
	Role            string
	Bio             datatypes.JSONSlice[string]
	Lore            datatypes.JSONSlice[string]
	MessageExamples datatypes.JSONSlice[MessageExample]

	Tools []Tool `gorm:"many2many:agents_tools"`

	Metadata datatypes.JSONType[map[string]string]
}
