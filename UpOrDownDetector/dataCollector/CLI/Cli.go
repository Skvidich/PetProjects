package CLI

import (
	"bufio"
	"dataCollector/StatusCoordinator"
	"dataCollector/StatusRelay"
	"dataCollector/common"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func RunCLI() {

	config, err := common.GetIni("./dataCollector.ini")
	if err != nil {
		fmt.Println("error while parsing ini ", err.Error())
		return
	}

	common.StatusLoggerInit(config.StatusLogPath)
	common.ErrorLoggerInit(config.ErrorLogPath)
	cord := StatusCoordinator.NewStatusCoordinator(config.StatusRequestDelay, config.GetterNames)
	relay := StatusRelay.NewStatusRelay(cord.OutChan, config.RelayIsLog, config.RelayIsResend)
	err = relay.InitProducer(config.KafkaTopic, config.KafkaAddrs)
	if err != nil {
		fmt.Println("error while initialising kafka producer ", err.Error())
		return
	}
	relay.InitProcess()
	cord.RunAll()
	relay.Run()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Status Coordinator CLI")
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
			relay.Close()
			common.KillErrorLogger()
			common.KillStatusLogger()

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

func printAllGetters(cord *StatusCoordinator.StatusCoordinator) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSTATUS")
	for _, info := range cord.GetListInfo() {
		fmt.Fprintf(w, "%s\t%s\n", info.Name, info.State)
	}
	w.Flush()
}

func printGetterInfo(cord *StatusCoordinator.StatusCoordinator, name string) {
	if !cord.IsGetterExist(name) {
		fmt.Printf("Getter '%s' not found\n", name)
		return
	}
	info := cord.GetInfo(name)
	fmt.Printf("Getter: %s\nStatus: %s\n", info.Name, info.State)
}

func handleStartCommand(cord *StatusCoordinator.StatusCoordinator, parts []string) {
	if len(parts) < 2 {
		fmt.Println("Error: missing target")
		return
	}

	target := parts[1]

	if !cord.IsGetterExist(target) {
		fmt.Printf("Getter '%s' not found\n", target)
		return
	}
	cord.RunGetter(target)
	fmt.Printf("Getter '%s' started\n", target)

}

func handleStopCommand(cord *StatusCoordinator.StatusCoordinator, parts []string) {
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
		if !cord.IsGetterExist(target) {
			fmt.Printf("Getter '%s' not found\n", target)
			return
		}
		cord.StopGetter(target)
		fmt.Printf("Getter '%s' stopped\n", target)
	}
}
