package main

import (
	"expvar"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"contrib.go.opencensus.io/exporter/zipkin"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"nimbler_writer/cmd/api/internal/handlers"
	"nimbler_writer/config"
	"nimbler_writer/internal/platform/database"
	contract "nimbler_writer/proto"
)

/*
ZipKin: http://localhost:9411
AddLoad: hey -m GET -c 10 -n 10000 "http://localhost:3000/v1/create_url"
expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,mem:memstats.Alloc"
*/

// build is the git version of this program. It is set using build flags in the
// makefile.
var build = "develop"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}


func run() error {
	// =========================================================================
	// Logging

	log := log.New(os.Stdout, "WRITER : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	cfg, err := config.Parse("config/config.yaml")
	if err != nil {
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting

	// Print the build version for our logs. Also expose it under /debug/vars.
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer log.Println("main : Completed")

	// =========================================================================
	// Start Tracing Support

	log.Println("main : Started : Initializing zipkin tracing support")

	localEndpoint, err := openzipkin.NewEndpoint(cfg.Zipkin.ServiceName, cfg.Zipkin.LocalEndpoint)
	if err != nil {
		return err
	}

	reporter := zipkinHTTP.NewReporter(cfg.Zipkin.ReporterURI)
	ze := zipkin.NewExporter(reporter, localEndpoint)

	trace.RegisterExporter(ze)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(cfg.Zipkin.Probability),
	})

	defer func() {
		log.Printf("main : Tracing Stopping : %s", cfg.Zipkin.LocalEndpoint)
		reporter.Close()
	}()

	// =========================================================================
	// Start Database

	log.Println("main : Started : Initializing database support")

	db, err := database.Open(database.Config{
		User:       cfg.Database.User,
		Password:   cfg.Database.Password,
		Host:       cfg.Database.Host,
		Name:       cfg.Database.Name,
		DisableTLS: cfg.Database.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		log.Printf("main : Database Stopping : %s", cfg.Database.Host)
		db.Close()
	}()

	// =========================================================================
	// Start API Service

	log.Println("main : Started : Initializing API support")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	listner, err := net.Listen("tcp", cfg.Service.APIHost)
	if err != nil {
		return errors.Wrap(err, "failed to listen tcp address")
	}
	s := grpc.NewServer()


	// Register Server
	srv := handlers.NewServer(db)
	contract.RegisterWriterServer(s, srv)
	servingErr:= s.Serve(listner)


	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", listner.Addr())
		serverErrors <- servingErr
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown", sig)
		// Asking listener to shutdown and load shed.
		err := listner.Close()
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.Service.ShutdownTimeout, err)
			s.GracefulStop()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}