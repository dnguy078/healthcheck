package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dnguy078/healthcheck/pkg/api"
	"github.com/dnguy078/healthcheck/pkg/service"
	"github.com/dnguy078/healthcheck/pkg/storage"
)

var (
	address   string
	sslCert   string
	sslKey    string
	runSSL    bool
	frequency string
	dataFile  string
)

func init() {
	flag.StringVar(&address, "bind", "127.0.0.1:8080", "address to bind to")
	flag.StringVar(&sslCert, "sslCert", "cert.pem", "ssl cert")
	flag.StringVar(&sslKey, "sslKey", "key.pem", "ssl key")
	flag.StringVar(&frequency, "checkfrequency", "3s", "frequency to run registered healthchecks")
	flag.StringVar(&dataFile, "datafile", "./pkg/storage/temp/data.json", "file containing existing healthchecks, loaded from disk")
	flag.BoolVar(&runSSL, "runSSL", false, "run with ssl")
	flag.Parse()
}

func main() {
	db := storage.NewCollection(dataFile)

	checkfrequency, err := time.ParseDuration(frequency)
	if err != nil {
		log.Fatal(err)
	}

	reporter, err := service.NewReporter(checkfrequency, db)
	if err != nil {
		log.Fatal(err)
	}
	reporter.Report()

	s, err := api.NewServer(address, sslCert, sslKey, db)
	if err != nil {
		log.Fatal(err)
	}
	s.Start()
	handleGracefulShutdown(s, db, reporter)
}

// handleGracefulShutdown listens for sig iterrupts, kills to gracefully shutdown. Existing healthchecks
// are written to disk
func handleGracefulShutdown(api *api.Server, db *storage.Collection, reporter *service.Reporter) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGHUP)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Dump(dataFile); err != nil {
		log.Printf("unable to write existing storage to disk, err: %s", err)
	}

	log.Printf("wrote existing healthchecks to disk, file: %s", dataFile)

	reporter.Stop()

	if err := api.Stop(ctx); err != nil {
		log.Printf("unable to stop http server, err: %s", err)
	}

}
