package agentruntime

import (
	"fmt"
	"github.com/habiliai/agentruntime/agent"
	"github.com/habiliai/agentruntime/config"
	"github.com/habiliai/agentruntime/internal/db"
	di "github.com/habiliai/agentruntime/internal/di"
	"github.com/habiliai/agentruntime/internal/mylog"
	"github.com/habiliai/agentruntime/runtime"
	"github.com/habiliai/agentruntime/thread"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"gorm.io/gorm"
	"net"
	"os"
)

func newServeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve <agent-file OR agent-files-dir>",
		Short: "Start agent runtime server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.Errorf("agent-file or agent-files-dir is required")
			}

			var agentFiles []string
			if stat, err := os.Stat(args[0]); os.IsNotExist(err) {
				return errors.Wrapf(err, "agent-file or agent-files-dir does not exist")
			} else if stat.IsDir() {
				files, err := os.ReadDir(args[0])
				if err != nil {
					return errors.Wrapf(err, "failed to read agent-files-dir")
				}
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					agentFiles = append(agentFiles, fmt.Sprintf("%s/%s", args[0], file.Name()))
				}
			} else {
				agentFiles = append(agentFiles, args[0])
			}

			ctx := cmd.Context()
			ctx = di.WithContainer(ctx, di.EnvProd)

			agentManager := di.MustGet[agent.Manager](ctx, agent.ManagerKey)
			cfg := di.MustGet[*config.RuntimeConfig](ctx, config.RuntimeConfigKey)
			logger := di.MustGet[*mylog.Logger](ctx, mylog.Key)
			threadManagerServer := di.MustGet[thread.ThreadManagerServer](ctx, thread.ManagerServerKey)
			dbInstance := di.MustGet[*gorm.DB](ctx, db.Key)
			agentManagerServer := di.MustGet[agent.AgentManagerServer](ctx, agent.ManagerServerKey)
			runtimeServer := di.MustGet[runtime.AgentRuntimeServer](ctx, runtime.ServerKey)

			if err := db.AutoMigrate(dbInstance); err != nil {
				return errors.Wrapf(err, "failed to migrate database")
			}

			agentConfigs, err := config.LoadAgentsFromFiles(agentFiles)
			for _, ac := range agentConfigs {
				if _, err := agentManager.SaveAgentFromConfig(ctx, ac); err != nil {
					return err
				}

				logger.Info("Agent loaded", "name", ac.Name)
			}

			lc := net.ListenConfig{}
			listener, err := lc.Listen(ctx, "tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
			if err != nil {
				return errors.Wrapf(err, "failed to listen on %s:%d", cfg.Host, cfg.Port)
			}

			logger.Info("Starting server", "host", cfg.Host, "port", cfg.Port)

			server := grpc.NewServer()
			grpc_health_v1.RegisterHealthServer(server, health.NewServer())
			thread.RegisterThreadManagerServer(server, threadManagerServer)
			agent.RegisterAgentManagerServer(server, agentManagerServer)
			runtime.RegisterAgentRuntimeServer(server, runtimeServer)

			go func() {
				<-ctx.Done()
				server.GracefulStop()
			}()

			return server.Serve(listener)
		},
	}

	return cmd
}
