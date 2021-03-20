package podimator

import(
    "fmt"

    . "github.com/IveGotNorto/podimator/internal/podimator"
    . "github.com/IveGotNorto/podimator/internal/config"

    "github.com/cavaliercoder/grab"
)

var podi Podimator

func init() {
    podi = Podimator {
        Config: ConfigParse("podcasts.json"),
        Client: grab.NewClient(),
    }

    podi.Config.Setup()
    podi.Client.UserAgent = "Podimator"

    fmt.Println("Podimator setup complete...")
}

// Process all existing podcasts
func All(name string) {
    if name == "" || contains(name) {
          
    }
    fmt.Printf("Getting all episodes for show: %v\n", name)
}

// Get all current episodes of podcasts
func Update(name string) {
    if name == "" || contains(name) {

    }
    fmt.Printf("Updating show: %v\n", name)
}

func contains(name string) bool {
    for _, podcast := range podi.Config.Podcasts {
        if podcast.Name == name {
            return true
        }
    }
    return false
}
