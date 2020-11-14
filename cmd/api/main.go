package main

import (
	"expvar"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	exporter "contrib.go.opencensus.io/exporter/zipkin"
	keygen "github.com/igomonov88/nimbler_key_generator/proto"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	logreporter "github.com/openzipkin/zipkin-go/reporter/log"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"

	"github.com/igomonov88/nimbler_writer/cmd/api/internal/handlers"
	"github.com/igomonov88/nimbler_writer/config"
	"github.com/igomonov88/nimbler_writer/internal/platform/database"
	contract "github.com/igomonov88/nimbler_writer/proto"
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

	reporter := logreporter.NewReporter(log)
	zipkinEndpoint, err := openzipkin.NewEndpoint(cfg.Zipkin.ServiceName,cfg.Zipkin.LocalEndpoint)
	if err != nil {
		return err
	}

	exp := exporter.NewExporter(reporter, zipkinEndpoint)
	trace.RegisterExporter(exp)
	trace.ApplyConfig(trace.Config{
		DefaultSampler: trace.ProbabilitySampler(cfg.Zipkin.Probability),
	})

	tracer, err := openzipkin.NewTracer(reporter, openzipkin.WithLocalEndpoint(zipkinEndpoint))
	if err != nil {
		return err
	}
	defer func() {
		log.Printf("main : Tracing Stopping : %s", cfg.Zipkin.LocalEndpoint)
		reporter.Close()
	}()

	// =========================================================================
	// Start ReaderServer Connection
	log.Println("main : Started : Initializing server grpc connection")


	keyGeneratorConn, err := grpc.Dial(cfg.KeyGenerator.APIHost, grpc.WithInsecure(), grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)))
	if err != nil {
		return errors.Wrap(err, "failed to connect to grpc writer server")
	}

	keyGenClient := keygen.NewKeyGeneratorClient(keyGeneratorConn)

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

	lis, err := net.Listen("tcp", cfg.Service.APIHost)
	if err != nil {
		return errors.Wrap(err, "failed to listen tcp connection")
	}

	grpcServer := grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)))
	srv := handlers.NewServer(db, keyGenClient)
	contract.RegisterWriterServer(grpcServer, srv)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API listening on %s", cfg.Service.APIHost)
		serverErrors <- grpcServer.Serve(lis)
	}()

	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown", sig)
		grpcServer.GracefulStop()
	}

	return nil
}