package main

import (
	"ecom/api"
	"ecom/config"
	"ecom/pkg/constants"
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
	debug          bool
)

func init() {
	flag.StringVar(&configPath, "c", "config.yml", "config file path")
	flag.IntVar(&verbosityLevel, "v", -1, "verbosity level, higher value - more logs")
	flag.BoolVar(&debug, "debug", false, "run service in debug mode")
	flag.Parse()
}

func main() {
	if debug {
		constants.DebugMode = debug
		fmt.Printf("\033[1;33m%s\033[0m", "Running in debug mode\n")
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if !config.LoadConfig(configPath) {
		os.Exit(-1)
	}

	if verbosityLevel < 0 {
		verbosityLevel = config.Config.Log.Level
	}
	zerolog.SetGlobalLevel(zerolog.Level(verbosityLevel))

	//api.InitMySqlCon(config.Config.MySql)
	//service.MakeSqlSchema(api.MySqlClient)
	//service.ApplyMigrations(api.MongoClient, api.MySqlClient)

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
