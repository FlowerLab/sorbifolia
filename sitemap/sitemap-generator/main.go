package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"go.x2ox.com/sorbifolia/sitemap"
)

var (
	filename string
	config   Config
)

func init() {
	flag.StringVar(&filename, "filename", "data.json", "")
	flag.StringVar(&filename, "f", "data.json", "")
}

func main() {
	flag.Parse()
	if err := parseConfig(); err != nil {
		log.Fatal(err)
	}

	if config.ImportURL == "" {
		sm := sitemap.New()
		for _, v := range findFile(config.ParseDirectory, "") {
			sm.Add(sitemap.NewURL(path.Join(config.DomainName, v)).
				SetLastMod(time.Now().Format("2006-01-02 15:04:05")).
				SetPriority(1.0),
			)
		}
		file, err := os.OpenFile(filepath.Join(config.OutputDirectory, "sitemap.xml"), os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err = file.Close(); err != nil {
				log.Printf("Error closing file: %s\n", err)
			}
		}()

		if err = sm.Write(file); err != nil {
			log.Fatal(err)
		}
	}
	sm := importURL(config.ImportURL)
	if sm == nil {
		sm = sitemap.New()
	}
	for _, v := range findFile(config.ParseDirectory, "") {
		sm.AddOrUpdate(
			sitemap.NewURL(path.Join(config.DomainName, v)).
				SetLastMod(time.Now().Format("2006-01-02 15:04:05")).
				SetChangeFreq(sitemap.Daily),
		)
	}
	file, err := os.OpenFile(filepath.Join(config.OutputDirectory, "sitemap.xml"), os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
		}
	}()

	if err = sm.Write(file); err != nil {
		log.Fatal(err)
	}
}

func parseConfig() error {
	/* #nosec G304 */
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("Error closing file: %s\n", err)
		}
	}()

	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return err
	}
	if config.DomainName == "" {
		return fmt.Errorf("domain name is empty")
	}
	if config.ParseDirectory == "" {
		return fmt.Errorf("parse directory is empty")
	}
	if config.OutputDirectory == "" {
		config.OutputDirectory = config.ParseDirectory
	}

	return nil
}

type Config struct {
	DomainName      string `json:"domain_name"`
	ParseDirectory  string `json:"parse_directory"`
	ImportURL       string `json:"import_url,omitempty"`
	OutputDirectory string `json:"output_directory,omitempty"` // default is ParseDirectory
}
