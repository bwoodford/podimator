package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/IveGotNorto/podimator/podcast"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Location string
	Podcasts []*podcast.Podcast
}

// Add error return
func (config *Config) Setup() {
	err := os.MkdirAll(config.Location, 755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "there was an issue creating the podcast file directory: %v", err)
	}

	for _, podcast := range config.Podcasts {
		// Create directories
		fullPath := fmt.Sprintf("%v/%v", config.Location, podcast.Name)
		err = os.MkdirAll(fullPath, 755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating %s directory: %v", podcast.Name, err)
			os.Exit(1)
		}
	}
}

// Add error return
func Parse(path string) *Config {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open configuration file: %v", err)
		os.Exit(1)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read configuration file: %v", err)
		os.Exit(1)
	}

	var config Config
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse configuration file: %v", err)
		os.Exit(1)
	}
	return &config
}
