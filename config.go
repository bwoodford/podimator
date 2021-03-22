package podimator

import(
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"
)

type Config struct {
    Location string `json:"location"`
    Podcasts []Podcast `json:"podcasts"`
}

func (config *Config) Setup() {
    err := os.MkdirAll(config.Location, 755)
    if err != nil {
        os.Exit(1)
    }

    for _, podcast := range config.Podcasts {
        // Create directories
        fullPath := fmt.Sprintf("%v/%v", config.Location, podcast.Name)
        err = os.MkdirAll(fullPath, 755)
        if err != nil {
            os.Exit(1)
        }
        // Initially unmark all podcasts
        podcast.Process = false
    }
}

func ConfigParse(path string) *Config {
    file, err := os.Open(path)
    if err != nil {
        os.Exit(1)
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        os.Exit(1)
    }

    var config Config
    err = json.Unmarshal(bytes, &config)
    if (err != nil) {
        os.Exit(1)
    }
    return &config
}
