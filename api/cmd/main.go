package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/felipeparaujo/rinha-de-backend-24q1/api"
)

const maxRetries = 100

var (
	cpuProfile = flag.String("cpu-profile", "", "write cpu profile to file")
	port       = flag.Int("port", 8080, "api port")
)

func main() {
	flag.Parse()

	var cpuProfFile *os.File
	if *cpuProfile != "" {
		log.Println("CPU profiling enabled")
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		cpuProfFile = f
		pprof.StartCPUProfile(f)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
	go func() {
		<-c
		log.Print("Shutting down...")
		if *cpuProfile != "" {
			pprof.StopCPUProfile()
			cpuProfFile.Close()
		}

		os.Exit(1)
	}()

	app := api.App{
		Port: int(*port),
	}

	if err := app.Listen(); err != nil {
		os.Exit(1)
	}
}
