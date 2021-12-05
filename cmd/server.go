package cmd

import (
	"github.com/wyy-go/go-web-template/internal/dao"
	"github.com/wyy-go/go-web-template/internal/jobs"
	"github.com/wyy-go/go-web-template/internal/routers"
	"github.com/wyy-go/go-web-template/pkg/env"
	"github.com/wyy-go/go-web-template/pkg/logger"
	"github.com/wyy-go/go-web-template/pkg/logger/zap"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/wyy-go/go-web-template/internal/config"
)

var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "run server",
	Example: "go-web-template server",
	PreRun: func(cmd *cobra.Command, args []string) {
		log.Println("PreRun")
		setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Run")
		run()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func setup() {
	// load config
	if err := config.LoadConfig(cfgFile); err != nil {
		log.Fatal(err)
	}

	d := env.ToDeploy(config.GetConfig().Server.DeployEnv)
	if d == env.DeployUnknown {
		log.Fatal("未设置发布模式")
	}
	env.SetDeploy(d)

	// logger
	lvl, err := logger.GetLevel(viper.GetString("logger.level"))
	if err != nil {
		log.Fatal(err)
	}

	w := &lumberjack.Logger{
		Filename:   config.GetConfig().LoggerConfig.Filename,
		MaxSize:    config.GetConfig().LoggerConfig.MaxSize,
		MaxBackups: config.GetConfig().LoggerConfig.MaxBackups,
		MaxAge:     config.GetConfig().LoggerConfig.MaxAge,
		Compress:   config.GetConfig().LoggerConfig.Compress,
	}
	if env.IsDeployRelease() {
		l, err := zap.NewLogger(logger.WithLevel(lvl), zap.WithCallerSkip(2), logger.WithWriter([]io.Writer{w}))
		if err != nil {
			log.Fatal(err)
		}
		logger.DefaultLogger = l
	} else {
		l, err := zap.NewLogger(logger.WithLevel(lvl), zap.WithCallerSkip(2), logger.WithWriter([]io.Writer{os.Stderr, w}))
		if err != nil {
			log.Fatal(err)
		}
		logger.DefaultLogger = l
	}

	daoConfig := dao.Config{}
	daoConfig.DriverName = "mysql"
	daoConfig.DataSourceName = config.GetConfig().MysqlConfig.DataSource
	daoConfig.MaxIdleConn = config.GetConfig().MysqlConfig.MaxIdle
	daoConfig.MaxOpenConn = config.GetConfig().MysqlConfig.MaxOpen
	//dao.Setup(daoConfig)

	// schedule
	jobs.Setup()

}

func run() {
	mode := viper.GetString("mode")
	gin.SetMode(mode)
	engine := gin.New()
	routers.Setup(engine)
	port := config.GetConfig().Server.Port
	engine.Run(":" + port)
}
