package config

import(
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"

    "github.com/IveGotNorto/podimator/podcast"
)

type Config struct {
    Location string `json:"location"`
    Podcasts []podcast.Podcast `json:"podcasts"`
}

func (config *Config) Setup() {
    err := os.MkdirAll(config.Location, 755)
    if err != nil {
        fmt.Fprintf(os.Stderr, "there was an issue creating podcast file directory: %v", err)
        os.Exit(1)
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
    err = json.Unmarshal(bytes, &config)
    if (err != nil) {
        fmt.Fprintf(os.Stderr, "unable to parse configuration file: %v", err)
        os.Exit(1)
    }
    return &config
}
