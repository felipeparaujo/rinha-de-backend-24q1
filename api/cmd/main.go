package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxRetries = 100

var (
	cpuProfile   = flag.String("cpuprofile", "", "write cpu profile to file")
	port         = flag.Int("port", 8080, "api port")
	dbConnString = flag.String("db-conn-str", "postgres://admin:password@db:5432/rinha", "DB connection string")
)

type Config struct {
	Port           int
	CPUProfilePath string
	DBConnStr      string
}

func main() {
	flag.Parse()

	appConfig := Config{
		Port:           int(*port),
		CPUProfilePath: string(*cpuProfile),
		DBConnStr:      string(*dbConnString),
	}

	ctx := context.Background()
	config, err := pgxpool.ParseConfig(appConfig.DBConnStr)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	for pool.Ping(ctx) != nil {
		time.Sleep(250 * time.Millisecond)
	}

	if err := api.NewServer(pool).Listen(fmt.Sprintf(":%d", appConfig.Port)); err != nil {
		os.Exit(1)
	}
}
