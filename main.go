package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var tzLocation *time.Location

func init() {
	tzLocation, _ = time.LoadLocation("America/New_York")
}

func main() {
	cfg := loadConfig()
	// This is the main entry point for the CLI application.
	// The implementation will be added later.
	f, err := os.OpenFile("./server.properties", os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	out, err := os.Create("./server.properties.out")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer out.Close()

	sc := bufio.NewScanner(f)
	wr := bufio.NewWriter(out)

	for sc.Scan() {
		line := sc.Text()

		updated := updateLineConfig(line, cfg)

		writeToFile(wr, updated)
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
	}
	err = out.Close()
	if err != nil {
		log.Fatalf("Error closing file: %v", err)
	}

	err = os.Rename("./server.properties.out", "./server.properties")
	if err != nil {
		log.Fatalf("Error renaming file: %v", err)
	}
}

func updateLineConfig(line string, cfg map[string]string) string {
	if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
		return line
	}

	splits := strings.Split(line, "=")
	if len(splits) != 2 {
		return line
	}

	k := strings.ToLower(strings.TrimSpace(splits[0]))
	if v, ok := cfg[k]; ok {
		fmt.Printf("<%s>: current=%s new=%s\n", k, strings.TrimSpace(splits[1]), v)
		line = fmt.Sprintf("%s=%s # %s", splits[0], v, time.Now().In(tzLocation).Format(time.RFC3339))
	}

	return line
}

func writeToFile(wr *bufio.Writer, line string) {
	_, err := wr.WriteString(line + "\n")
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}
}
func loadConfig() map[string]string {
	config := map[string]string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "MC_") {
			l := strings.TrimPrefix(e, "MC_")
			splits := strings.Split(l, "=")
			if len(splits) == 2 {
				k := strings.ToLower(strings.Replace(splits[0], "_", "-", -1))
				config[k] = splits[1]
			}
		}
	}

	return config
}
