package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var StatusRequestDelay time.Duration = 30 * time.Second
var GetterNames []string = make([]string, 0)

var RelayIsLog bool = true
var RelayIsResend bool = true

var ErrorLogPath string = "..\\logs\\error.log"
var StatusLogPath string = "..\\logs\\status.log"

var KafkaAddrs []string = make([]string, 0)
var KafkaTopic string = ""

func GetIniParameters(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can't open ini file %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == ';' || line[0] == '#' {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("incorrect option")
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "StatusRequestDelay":
			StatusRequestDelay, err = time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("invalid duration format: %w", err)
			}

		case "GetterNames":
			GetterNames = strings.Split(value, ",")
			for i, name := range GetterNames {
				GetterNames[i] = strings.TrimSpace(name)
			}

		case "RelayIsLog":
			RelayIsLog, err = parseBool(value)
			if err != nil {
				return fmt.Errorf("invalid bool value for RelayIsLog: %w", err)
			}

		case "RelayIsResend":
			RelayIsResend, err = parseBool(value)
			if err != nil {
				return fmt.Errorf("invalid bool value for RelayIsResend: %w", err)
			}

		case "ErrorLogPath":
			ErrorLogPath = value

		case "StatusLogPath":
			StatusLogPath = value
		case "KafkaAddrs":
			KafkaAddrs = strings.Split(value, ",")
			for i, addr := range KafkaAddrs {
				KafkaAddrs[i] = strings.TrimSpace(addr)
			}
		case "KafkaTopic":
			KafkaTopic = value
		default:
			return fmt.Errorf("unknown option")
		}
	}

	return scanner.Err()

}

func parseBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "yes":
		return true, nil
	case "false", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", s)
	}
}
