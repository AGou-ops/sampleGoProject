package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"greenlight.agou-ops.cn/internal/data"
	"greenlight.agou-ops.cn/internal/jsonlog"
	"greenlight.agou-ops.cn/internal/mailer"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	dialer mailer.Mailer
	wg     sync.WaitGroup
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "development/production/tmp mode")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:123@localhost/greenlight?sslmode=disable", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "max-open-conns", 25, "Postgres database connection limit")
	flag.IntVar(&cfg.db.maxIdleConns, "max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "max-idle-time", "15m", "Postgress max idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limter maximum of requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limter")

	// 	SendSmtp("smtp.163.com:25", da"i15628960878@163.com", "TGMXZWNPZOJPCLCL", []string{"ictw@qq.com"}, "hello", "world")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.163.com", "smtp server hostname")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "smtp server port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "dai15628960878@163.com", "smtp server email username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "TGMXZWNPZOJPCLCL", "smtp server password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "dai15628960878@163.com", "smtp sender email address")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	displayVersion := flag.Bool("version", false, "display current version")

	flag.Parse()

	if *displayVersion {
		fmt.Println(version)
		fmt.Println(buildTime)
		os.Exit(0)
	}

	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	// expose custome metrics
	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
		dialer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
