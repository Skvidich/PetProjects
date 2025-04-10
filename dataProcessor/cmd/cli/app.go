package cli

import (
	"dataProcessor/internal/config"
	"dataProcessor/internal/core/alert"
	"dataProcessor/internal/core/processor"
	"dataProcessor/internal/core/reader"
	"dataProcessor/internal/core/storage"
	"dataProcessor/internal/logger"
	"dataProcessor/pkg/errors"
	"database/sql"
	"fmt"
	"github.com/IBM/sarama"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunCLI(iniPath string) {
	cfg, err := config.LoadConfig(iniPath)
	if err != nil {
		log.Fatalf("can't load config: %v", err)
	}

	errLog, err := logger.NewErrLogger(cfg.ErrLog)
	if err != nil {
		log.Fatalf("can't create logger: %v", err)
	}

	app, err := New(cfg, errLog)
	if err != nil {
		log.Fatalf("can't create app: %v", err)
	}
	endChan := make(chan struct{})
	stop := make(chan os.Signal)
	go func() {

		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
		<-stop
		app.Close()
		endChan <- struct{}{}
	}()
	err = app.Run()
	if err != nil {
		errLog.LogError("runtime error", err)
		stop <- syscall.SIGTERM
	}
	<-endChan
	log.Println("exiting....")
}

type App struct {
	cfg    *config.Config
	errLog logger.Logger
	rep    *alert.KafkaReporter
	stor   *storage.SQLStorage
	cons   *reader.KafkaConsumer
	proc   *processor.StatusEngine
	db     *sql.DB
}

func New(cfg *config.Config, errLog logger.Logger) (*App, error) {

	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}
	return &App{
		cfg:    cfg,
		errLog: errLog,
	}, nil
}

func (a *App) CreateStorage() error {
	var err error
	if a.db != nil {
		err = a.db.Close()
		if err != nil {
			a.errLog.LogError("can't close DB connection", err)
		}

	}

	a.db, err = sql.Open("pgx", a.cfg.DSN)
	if err != nil {
		a.errLog.LogError("can't create connection", err)
		return fmt.Errorf("can't create connection: %v", err)
	}
	err = a.db.Ping()
	if err != nil {
		a.errLog.LogError("error while ping", err)
		return fmt.Errorf("error while ping: %v", err)
	}
	a.stor, err = storage.NewSQLStorage(a.db, a.cfg.PostgresStorageConfig.ReportTable, a.cfg.PostgresStorageConfig.IncidentTable, a.cfg.PostgresStorageConfig.ComponentTable)
	if err != nil {
		a.errLog.LogError("can't create storage entity", err)
		return fmt.Errorf("can't create storage entity: %v", err)
	}
	return nil
}

func (a *App) CreateReporter() error {

	var err error
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Net.MaxOpenRequests = 5
	cfg.Producer.Partitioner = sarama.NewHashPartitioner
	a.rep, err = alert.NewKafkaReporter(cfg, a.cfg.KafkaReporterConfig.Brokers, a.cfg.StartTopic, a.cfg.EndTopic)
	if err != nil {
		a.errLog.LogError("can't create reporter", err)
		return fmt.Errorf("can't create reporter: %v", err)
	}
	return nil
}

func (a *App) CreateConsumer() error {

	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	var err error

	var backup *reader.ConcurrentQueue
	if a.cons != nil {
		backup = a.cons.BackupQueue()
	}
	a.cons, err = reader.NewConsumer(cfg, a.cfg.KafkaConsumerConfig.Brokers, backup)
	if err != nil {
		a.errLog.LogError("can't create consumer", err)
		return fmt.Errorf("can't create consumer: %v", err)
	}
	a.cons.ConsumeTopic(a.cfg.KafkaConsumerConfig.Topic)
	return nil
}

func (a *App) CreateProcessor() {
	var backupIncidents *processor.IncidentBuffer
	var backupStatistic *processor.MetricsBuffer
	if a.proc != nil {
		backupIncidents, backupStatistic = a.proc.GetBackup()
	}

	a.proc = processor.NewEngine(a.stor, a.cons, a.rep, a.cfg.AggregationInterval, backupStatistic, backupIncidents)
	go a.proc.Start(a.cfg.ProcessEngineConfig.ReadTimeout)
}
func (a *App) Run() error {
	var err error
	var counter int
	err = a.CreateStorage()
	for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
		a.errLog.LogError("failed creating storage", err)
		err = a.CreateStorage()
	}
	if err != nil {
		a.errLog.LogError("failed creating storage", err)
		return fmt.Errorf("runout of retries while creating storage")
	}

	err = a.CreateConsumer()
	for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
		a.errLog.LogError("failed creating consumer", err)
		err = a.CreateConsumer()
	}
	if err != nil {
		a.errLog.LogError("failed creating consumer", err)
		return fmt.Errorf("runout of retries while creating consumer")
	}

	err = a.CreateReporter()
	for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
		a.errLog.LogError("failed creating reporter", err)
		err = a.CreateReporter()
	}
	if err != nil {
		a.errLog.LogError("failed creating reporter", err)
		return fmt.Errorf("runout of retries while creating reporter")
	}

	a.CreateProcessor()

	var appErr errors.TaggedError
	appErr = <-a.proc.ErrChan
	for appErr.Err != nil {
		a.errLog.LogError("runtime error: ", appErr.Err)
		switch appErr.Type {
		case errors.TypeProcessor:
			return fmt.Errorf("unknown error")
		case errors.TypeSaver:
			err = a.CreateStorage()
			for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
				a.errLog.LogError("failed creating storage", err)
				err = a.CreateStorage()
			}
			if err != nil {
				a.errLog.LogError("failed creating storage", err)
				return fmt.Errorf("runout of retries while creating storage")
			}
		case errors.TypeConsumer:
			err = a.CreateConsumer()
			for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
				a.errLog.LogError("failed creating consumer", err)
				err = a.CreateConsumer()
			}
			if err != nil {
				a.errLog.LogError("failed creating consumer", err)
				return fmt.Errorf("runout of retries while creating consumer")
			}
		case errors.TypeNotificator:
			err = a.CreateReporter()
			for counter = 0; counter < a.cfg.RetryMax && err != nil; counter++ {
				a.errLog.LogError("failed creating reporter", err)
				err = a.CreateReporter()
			}
			if err != nil {
				a.errLog.LogError("failed creating reporter", err)
				return fmt.Errorf("runout of retries while creating reporter")
			}

		}
		a.CreateProcessor()
		appErr = <-a.proc.ErrChan
	}
	return nil
}

func (a *App) Close() {
	if a.proc != nil {
		tagErr := a.proc.Close()
		if tagErr.Err != nil {
			a.errLog.LogError("error while closing engine", tagErr.Err)
		}
	}

	var err error
	if a.db != nil {
		err = a.db.Close()
		if err != nil {
			a.errLog.LogError("error while closing DB connection", err)
		}
	}

	if a.cons != nil {
		err = a.cons.Close()
		if err != nil {
			a.errLog.LogError("error while closing consumer", err)
		}
	}

	if a.rep != nil {
		err = a.rep.Close()
		if err != nil {
			a.errLog.LogError("error while closing consumer", err)
		}
	}

}
