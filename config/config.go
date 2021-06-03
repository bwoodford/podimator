package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
)

type Podcast struct {
	URL  string
	Name string
}

type Config struct {
	Location string
	Podcasts []*Podcast
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

func Parse(path string) (*Config, []error) {

	errors := []error{}
	file, err := os.Open(path)
	if err != nil {
		errors = append(errors, fmt.Errorf("Unable to open configuration file: %v", err))
		return nil, errors
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		errors = append(errors, fmt.Errorf("Unable to read configuration file: %v", err))
		return nil, errors
	}

	var config Config
	err = toml.Unmarshal(bytes, &config)
	if err != nil {
		errors = append(errors, fmt.Errorf("Unable to parse configuration file: %v", err))
		return nil, errors
	}

	errors = append(errors, validate(&config)...)
	if len(errors) > 0 {
		return nil, errors
	}
	return &config, nil
}

func validate(config *Config) []error {
	errors := []error{}
	if config.Location == "" {
		errors = append(errors, fmt.Errorf("\tConfig location is empty"))
	}

	for i, pod := range config.Podcasts {
		if pod.Name == "" {
			errors = append(errors, fmt.Errorf("\tPodcast #%d: config name is empty", i))
		}
		if pod.URL == "" {
			errors = append(errors, fmt.Errorf("\tPodcast #%d: URL is empty", i))
		}
	}
	return errors
}
