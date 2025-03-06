package mytesting

import (
	"context"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/di"
	"github.com/habiliai/agentruntime/internal/db"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	context.Context

	Config *config.RuntimeConfig
	DB     *gorm.DB
}

func (s *Suite) SetupTest() {
	s.Context = context.TODO()
	s.Context = di.WithContainer(s.Context, di.EnvTest)

	s.Config = di.MustGet[*config.RuntimeConfig](s, config.RuntimeConfigKey)
	s.DB = di.MustGet[*gorm.DB](s, db.Key)
}

func (s *Suite) TearDownTest() {
	defer func() {
		if err := db.CloseDB(s.DB); err != nil {
			s.T().Logf("failed to close db: %v", err)
		}
	}()
}
