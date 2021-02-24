package config 

import(
    "encoding/json"
    "fmt"
    "os"
    "io/ioutil"

    "github.com/IveGotNorto/podimator/internal/podcast"
)

type Config struct {
    Location string `json:"location"`
    Podcasts []podcast.Podcast `json:"podcasts"`
}

func (config *Config) Setup() {
    err := os.MkdirAll(config.Location, 755)
    if err != nil {
        panic(err)
    }

    for _, podcast := range config.Podcasts {
        fullPath := fmt.Sprintf("%v/%v", config.Location, podcast.Name)
        err = os.MkdirAll(fullPath, 755)
        if err != nil {
            panic(err)
        }
    }
}

func ConfigParse(path string) Config {
    file, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        panic(err)
    }

    var config Config
    err = json.Unmarshal(bytes, &config)
    if (err != nil) {
        panic(err)
    }
    return config
}
