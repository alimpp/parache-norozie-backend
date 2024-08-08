package api

import (
	"context"
	"crypto/tls"
	"ecom/config"
	"ecom/pkg/constants"
	"ecom/pkg/services"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/swagger"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AppServer struct {
	app   *fiber.App
	cfg   *config.ConfStruct
	sms   services.SMS
	sqlDb *gorm.DB
	redis *redis.Client
}

var AppSrv *AppServer

var store = session.New()

func NewAppServer(cfg *config.ConfStruct) *AppServer {
	appSrv := &AppServer{cfg: cfg}

	if cfg.SMS.Url == "" {
		appSrv.sms = services.MockSMS{}
	}

	appSrv.sqlDb = InitSqlDb(config.Config)
	services.ApplyMigrations(appSrv.sqlDb)

	opts, err := redis.ParseURL(cfg.Redis.Url)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opts)

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	appSrv.redis = rdb

	app := fiber.New()

	// set up middlewares
	app.Use(requestIdMiddleware())
	app.Use(logMiddleware())
	app.Use(rateLimiterMiddleware())
	if appSrv.cfg.Observability.PProf.Enabled {
		app.Use(performanceMonitorMiddleware())
	}
	if appSrv.cfg.Observability.Prometheus.Enabled {
		app.Use(PrometheusMiddleware(app))
	}

	// REST endpoints
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/health", healthCheck)

	v1.Get(fmt.Sprintf("/swagger/%s/*", cfg.Swagger.Password), swagger.HandlerDefault)

	v1.Get("/login", login)
	v1.Post("/otp", verifyOtp)
	v1.Post("/password", password)

	appSrv.app = app
	AppSrv = appSrv
	return appSrv
}

func (s *AppServer) ListenAndServe() chan error {
	errCh := make(chan error)
	go func() {
		if s.cfg.Server.Cert != "" && s.cfg.Server.Key != "" {
			constants.Logger.Info().Msgf("Starting listening addr https://%s", s.cfg.Server.Addr)
			cer, err := tls.LoadX509KeyPair(s.cfg.Server.Cert, s.cfg.Server.Key)
			if err != nil {
				panic(err)
			}

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
			ln, err := tls.Listen("tcp", s.cfg.Server.Addr, tlsConfig)
			if err != nil {
				panic(err)
			}
			errCh <- s.app.Listener(ln)
		} else {
			constants.Logger.Info().Msgf("Starting listening addr http://%s", s.cfg.Server.Addr)
			errCh <- s.app.Listen(s.cfg.Server.Addr)
		}
	}()
	return errCh
}

func InitSqlDb(conf config.ConfStruct) *gorm.DB {
	if conf.DB.DriverName == "sqlite" {
		db, err := gorm.Open(sqlite.Open(conf.DB.DataSourceName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	} else if conf.DB.DriverName == "postgres" {
		db, err := gorm.Open(postgres.Open(conf.DB.DataSourceName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	} else {
		panic("unsupported database driver")
	}
}
