package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	ThreadID uint
	Thread   Thread `gorm:"foreignKey:ThreadID"`

	User    string
	Content datatypes.JSONType[MessageContent]
}

type MessageContent struct {
	Text            string `json:"text,omitempty"`
	Action          string `json:"action,omitempty"`
	ActionParameter string `json:"action_parameter,omitempty"`
	ActionResult    string `json:"action_result,omitempty"`
}
