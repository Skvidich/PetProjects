package server

import (
	"dataCollector/internal/app/server"
	"dataCollector/internal/config"
	"dataCollector/internal/core/coordinator"
	"dataCollector/internal/core/relay"
	"dataCollector/internal/core/storage"
	"dataCollector/internal/logger"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunServer(iniPath string) {
	cfg, err := config.LoadConfig(iniPath)
	log.Println("config parsing....")
	if err != nil {
		log.Fatalf("error while parsing ini %v ", err)
	}
	log.Println("config parsed")

	log.Println("logger starting....")
	errLogger, err := logger.NewErrLogger(cfg.ErrLog)
	if err != nil {
		log.Fatalf("can't create logger %v", err)
	}
	log.Println("logger started")

	log.Println("connecting to DB....")
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}
	log.Println("successfully connected to DB")
	log.Println("initializing storage")
	sqlStorage, err := storage.NewSQLStorage(db, cfg.ReportTable, cfg.ComponentTable)
	if err != nil {
		log.Fatalf("can't create logger %v", err)
	}
	log.Println("storage initialized")

	log.Println("initializing coordinator....")
	cord := coordinator.New(cfg.ReqDelay, cfg.Getters, errLogger)
	log.Println("coordinator initialized")
	log.Println("initializing relay....")
	rel := relay.NewRelay(cord.OutChan, cfg.RelayConfig.Save, cfg.RelayConfig.Resend, errLogger, sqlStorage)
	log.Println("relay initialized")

	log.Println("initializing kafka producer....")
	err = rel.SetupProducer(cfg.KafkaProducerConfig.Topic, cfg.KafkaProducerConfig.Brokers)
	if err != nil {
		log.Fatalf("error while initialising kafka producer %v", err)
	}
	log.Println("kafka producer initialized")

	log.Println("starting getters....")
	err = cord.StartAll()
	if err != nil {
		log.Fatalf("error can't start all getters %v", err)
	}
	log.Println("starting relay....")
	rel.Run()

	log.Println("starting gRPC server....")
	app := server.New(cord, errLogger, sqlStorage, cfg.Address)

	endChan := make(chan struct{})

	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
		<-stop
		log.Println("stopping gRPC server....")
		app.Stop()
		log.Println("stopping coordinator....")
		cord.Shutdown()
		log.Println("stopping relay....")
		rel.Close()
		_ = errLogger.Close()
		endChan <- struct{}{}
	}()
	err = app.Run()
	if err != nil {
		log.Fatalf("error while running server %v", err)
	}
	<-endChan
	log.Println("exiting....")
}
