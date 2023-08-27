package config

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml/v2"
)

var fsym fileSystem = osFS{}

type osFS struct{}

type fileSystem interface {
	Open(path string) (file, error)
	MkdirAll(path string, perm fs.FileMode) error
	ReadAll(r io.Reader) ([]byte, error)
}

type file interface {
	io.Closer
	io.Reader
}

type Podcast struct {
	URL  string
	Name string
}

type Config struct {
	Location string
	Podcasts []*Podcast
}

func (osFS) Open(path string) (file, error)               { return os.Open(path) }
func (osFS) MkdirAll(path string, perm fs.FileMode) error { return os.MkdirAll(path, perm) }
func (osFS) ReadAll(r io.Reader) ([]byte, error)          { return ioutil.ReadAll(r) }

func (config *Config) Setup() error {
	err := fsym.MkdirAll(config.Location, 755)
	if err != nil {
		return err
	}

	for _, podcast := range config.Podcasts {
		// Create directories
		fullPath := fmt.Sprintf("%v/%v", config.Location, podcast.Name)
		err = fsym.MkdirAll(fullPath, 755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating %s directory: %v\n", podcast.Name, err)
			os.Exit(1)
		}
	}
	return nil
}

func Parse(path string) (*Config, error) {
	file, err := fsym.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := fsym.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = toml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
