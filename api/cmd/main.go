package main

import (
	"flag"
	"os"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
)

const maxRetries = 100

var (
	cpuProfile   = flag.String("cpuprofile", "", "write cpu profile to file")
	port         = flag.Int("port", 8080, "api port")
	dbConnString = flag.String("db-conn-str", "postgres://admin:password@db:5432/rinha", "DB connection string")
)

func main() {
	flag.Parse()

	app := api.App{
		Port:      int(*port),
		DBConnStr: string(*dbConnString),
	}

	if err := app.Listen(); err != nil {
		os.Exit(1)
	}
}
