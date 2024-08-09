package main

import (
	"context"
	"ecom/api"
	"ecom/config"
	"ecom/pkg/constants"
	"flag"
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"

	_ "ecom/docs"
)

var (
	configPath     string
	verbosityLevel int
)

func init() {
	flag.StringVar(&configPath, "c", "config.yml", "config file path")
	flag.Parse()
}

// @title           Swagger Doc
// @version         1.0
func main() {
	if !config.LoadConfig(configPath) {
		os.Exit(-1)
	}

	if config.Config.Mode.Debug {
		fmt.Printf("\033[1;33m%s\033[0m", "----------------------------RUNNING IN DEBUG MODE----------------------------\n")
	}

	if verbosityLevel < 0 {
		verbosityLevel = config.Config.Log.Level
	}
	zerolog.SetGlobalLevel(zerolog.Level(verbosityLevel))

	ctx, cancel := context.WithCancel(context.Background())

	serv := api.NewAppServer(&config.Config, ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	)
	select {
	case err := <-serv.ListenAndServe():
		cancel()
		panic(err)
	case <-sigCh:
		constants.Logger.Info().Msg("Shutting down service...")
		cancel()
		os.Exit(1)
	}
}
