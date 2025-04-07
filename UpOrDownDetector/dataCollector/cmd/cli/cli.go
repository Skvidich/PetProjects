package cli

import (
	"bufio"
	"dataCollector/internal/config"
	"dataCollector/internal/core/coordinator"
	"dataCollector/internal/core/relay"
	"dataCollector/internal/core/storage"
	"dataCollector/internal/logger"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"strings"
	"text/tabwriter"
)

func RunCLI(iniPath string) {

	cfg, err := config.LoadConfig(iniPath)
	if err != nil {
		fmt.Println("error while parsing ini ", err.Error())
		return
	}

	errLogger, err := logger.NewErrLogger(cfg.ErrLog)
	if err != nil {
		fmt.Println("can't create logger ", err.Error())
		return
	}

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

	sqlStorage, err := storage.NewSQLStorage(db, cfg.ReportTable, cfg.ComponentTable)
	if err != nil {
		fmt.Println("can't create logger ", err.Error())
		return
	}

	cord := coordinator.New(cfg.ReqDelay, cfg.Getters, errLogger)
	rel := relay.NewRelay(cord.OutChan, cfg.RelayConfig.Save, cfg.RelayConfig.Resend, errLogger, sqlStorage)
	err = rel.SetupProducer(cfg.KafkaProducerConfig.Topic, cfg.KafkaProducerConfig.Brokers)
	if err != nil {
		fmt.Println("error while initialising kafka producer ", err.Error())
		return
	}
	cord.StartAll()
	rel.Run()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Status coordinator cli")
	fmt.Println("Enter 'help' for available commands")

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Split(input, " ")

		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "help":
			printHelp()

		case "list":
			printAllGetters(cord)

		case "info":
			if len(parts) < 2 {
				fmt.Println("Error: missing getter name")
				break
			}
			printGetterInfo(cord, parts[1])

		case "start":
			handleStartCommand(cord, parts)

		case "stop":
			handleStopCommand(cord, parts)

		case "shutdown":
			cord.Shutdown()
			rel.Close()
			err = errLogger.Close()
			if err != nil {
				fmt.Println("error while closing logger ", err.Error())
			}
			fmt.Println("System shutdown completed")
			return

		default:
			fmt.Println("Unknown command. Type 'help' for available commands")
		}
	}
}

func printHelp() {
	fmt.Println(`
Available commands:
  list               - List all getters and their statuses
  info <name>        - Get info about specific getter
  start <name>       - Start getter(s)
  stop <name>    	 - Stop getter(s)
  shutdown           - Full system shutdown
  help               - Show this help`)
}

func printAllGetters(cord *coordinator.Coordinator) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")
	for _, info := range cord.GetterList() {
		fmt.Fprintf(w, "%s\t%s\n", info.Name, info.State)
	}
	w.Flush()
}

func printGetterInfo(cord *coordinator.Coordinator, name string) {

	info, err := cord.Getter(name)
	if err != nil {
		fmt.Println("Error", err.Error())
	} else {
		fmt.Printf("Getter: %s\nStatus: %s\n", info.Name, info.State)
	}

}

func handleStartCommand(cord *coordinator.Coordinator, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: missing target")
		return
	}

	target := parts[1]

	err := cord.Start(target)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	} else {
		fmt.Printf("Getter '%s' started\n", target)
	}

}

func handleStopCommand(cord *coordinator.Coordinator, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: missing target")
		return
	}

	target := parts[1]
	switch target {
	case "all":
		cord.StopAll()
		fmt.Println("All getters stopped")
	default:
		err := cord.Stop(target)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		} else {
			fmt.Printf("Getter '%s' stopped\n", target)
		}

	}
}
