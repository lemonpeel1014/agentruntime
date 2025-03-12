package entity

import (
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model

	Instruction string

	Participants []Agent   `gorm:"many2many:thread_participants;"`
	Messages     []Message `gorm:"foreignKey:ThreadID"`
}
