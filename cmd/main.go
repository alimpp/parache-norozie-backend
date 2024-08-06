package main

import (
	"ecom/api"
	"ecom/config"
	"ecom/pkg/constants"
	"ecom/pkg/services"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

var (
	configPath     string
	verbosityLevel int
)

func init() {
	flag.StringVar(&configPath, "c", "config.yml", "config file path")
	flag.Parse()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if !config.LoadConfig(configPath) {
		os.Exit(-1)
	}

	if config.Config.Mode.Debug {
		fmt.Printf("\033[1;33m%s\033[0m", "Running in debug mode\n")
	}

	if verbosityLevel < 0 {
		verbosityLevel = config.Config.Log.Level
	}
	zerolog.SetGlobalLevel(zerolog.Level(verbosityLevel))

	api.InitSqlDb(config.Config)
	services.ApplyMigrations(api.SqlDb)

	serv := api.NewAppServer(&config.Config)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	select {
	case err := <-serv.ListenAndServe():
		panic(err)
	case <-sigCh:
		constants.Logger.Info().Msg("Shutting down service...")
		os.Exit(1)
	}
}
