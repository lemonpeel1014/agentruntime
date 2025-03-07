package db

import (
	"github.com/habiliai/agentruntime/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec("CREATE SCHEMA IF NOT EXISTS agentruntime").Error; err != nil {
		return errors.Wrapf(err, "failed to create schema")
	}

	return errors.WithStack(db.AutoMigrate(
		&entity.Agent{},
		&entity.Message{},
		&entity.Tool{},
		&entity.Thread{},
	))
}

func DropAll(db *gorm.DB) error {
	return errors.WithStack(db.Migrator().DropTable(
		"thread_participants",
		&entity.Thread{},
		"agents_tools",
		&entity.Tool{},
		&entity.Message{},
		&entity.Agent{},
	))
}
