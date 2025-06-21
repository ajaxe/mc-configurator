package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	_ "time/tzdata"
)

var tzLocation *time.Location

func init() {
	tzLocation, _ = time.LoadLocation("America/New_York")
	if tzLocation == nil {
		log.Fatalf("error loading timezone: %s", "America/New_York")
	}
	fmt.Printf("Timezone: %s\n", tzLocation.String())
}

type GeneratorOptions struct {
	Source      string // Source file path
	Destination string // Destination file path, where the configuration will be written.
}

type ConfigGenerator struct {
	options GeneratorOptions // options holds the configuration options for the generator.
}

func (c *ConfigGenerator) Execute() {
	cfg := loadEnvConfig()
	src := c.options.Source
	dest := c.options.Destination

	f, err := os.OpenFile(src, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	fmt.Printf("opening file: %s\n", src)

	tmpFile := dest + ".out"
	out, err := os.Create(tmpFile)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer out.Close()

	fmt.Printf("creating file: %s\n", tmpFile)

	sc := bufio.NewScanner(f)
	wr := bufio.NewWriter(out)

	for sc.Scan() {
		line := sc.Text()

		updated := updateLineConfig(line, cfg)

		writeToFile(wr, updated)
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v", err)
	}
	err = out.Close()
	if err != nil {
		log.Fatalf("error closing file: %v", err)
	}

	err = os.Rename(tmpFile, dest)
	if err != nil {
		log.Fatalf("error renaming file: %v", err)
	}
}

// writeToFile writes the line to the file with a newline character.
func writeToFile(wr *bufio.Writer, line string) {
	_, err := wr.WriteString(line + "\n")
	if err != nil {
		log.Fatalf("error writing to file: %v", err)
	}
}

// loadEnvConfig loads the environment variables that start with "MC_" and returns them as a map.
func loadEnvConfig() map[string]string {
	config := map[string]string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "MC_") {
			fmt.Printf("env: %s\n", e)
			l := strings.TrimPrefix(e, "MC_")
			splits := strings.Split(l, "=")
			if len(splits) == 2 {
				k := ""
				switch splits[0] {
				case "RCON_PASSWORD", "RCON_PORT":
					k = strings.Replace(strings.ToLower(splits[0]), "_", ".", 1)
					config[k] = strings.Replace(strings.ToLower(splits[0]), "_", ".", 1)
				default:
					k = strings.ToLower(strings.Replace(splits[0], "_", "-", -1))
				}

				config[k] = splits[1]
			}
		}
	}

	return config
}

// updateLineConfig updates the line with the new value from the config map if it exists.
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
		line = fmt.Sprintf("%s=%s\n# %s", splits[0], v, time.Now().In(tzLocation).Format(time.RFC3339))
	}

	return line
}
