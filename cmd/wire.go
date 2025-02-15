package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leehai1107/The-journal/cmd/banner"
	apifx "github.com/leehai1107/The-journal/di"
	"github.com/leehai1107/The-journal/pkg/config"
	"github.com/leehai1107/The-journal/pkg/errors"
	"github.com/leehai1107/The-journal/pkg/graceful"
	"github.com/leehai1107/The-journal/pkg/infra"
	"github.com/leehai1107/The-journal/pkg/logger"
	"github.com/leehai1107/The-journal/pkg/middleware/cors"
	"github.com/leehai1107/The-journal/pkg/recover"
	"github.com/leehai1107/The-journal/pkg/swagger"
	"github.com/leehai1107/The-journal/pkg/utils/ginbuilder"
	"github.com/leehai1107/The-journal/pkg/utils/timeutils"
	"github.com/leehai1107/The-journal/service/blog/delivery/http"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Command of Internal Service",
	Long:  "CLI used to manage internal apis, datas when users access.",
	Run: func(_ *cobra.Command, _ []string) {
		NewServer().Run()
	},
	Version: "1.0.0",
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	app := fx.New(
		fx.Invoke(config.InitConfig),
		fx.Invoke(initLogger),
		fx.Invoke(errors.Initialize),
		fx.Invoke(timeutils.Init),
		fx.Invoke(infra.InitPostgresql),
		//... add module here
		apifx.Module,
		fx.Provide(provideGinEngine),
		fx.Invoke(
			registerService,
			registerSwaggerHandler),
		fx.Invoke(startServer),
		fx.Invoke(banner.Print),
	)
	logger.Infow("Server started!")
	app.Run()
}

func provideGinEngine() *gin.Engine {
	return ginbuilder.BaseBuilder().Build()
}

func registerService(
	g *gin.Engine,
	router http.Router,
) {
	internal := g.Group("/internal")
	internal.Use(
		recover.RPanic,
		cors.CorsCfg(config.ServerConfig().CorsProduction))
	router.Register(internal)
}

func registerSwaggerHandler(g *gin.Engine) {
	swaggerAPI := g.Group("/internal/swagger")
	swag := swagger.NewSwagger()
	swaggerAPI.Use(swag.SwaggerHandler(config.ServerConfig().Production))
	swag.Register(swaggerAPI)
}

func startServer(lifecycle fx.Lifecycle, g *gin.Engine) {
	gracefulService := graceful.NewService(graceful.WithStopTimeout(time.Second), graceful.WithWaitTime(time.Second))

	gracefulService.Register(g)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", config.ServerConfig().HTTPPort)
				fmt.Println("run on port:", port)
				go gracefulService.StartServer(g, port)
				return nil
			},
			OnStop: func(context.Context) error {
				gracefulService.Close()
				infra.ClosePostgresql() // nolint
				return nil
			},
		},
	)
}
func initLogger() {
	logger.Initialize(config.ServerConfig().Logger)
}
