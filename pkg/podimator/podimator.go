package podimator

import(
    "github.com/IveGotNorto/podimator/internal/podimator"
    "github.com/IveGotNorto/podimator/internal/config"

    "github.com/cavaliercoder/grab"
)

func Podimator() {
    podi := podimator.Podimator{
        Config: config.ConfigParse("podcasts.json"),
        Client: grab.NewClient(),
    }
    podi.Config.Setup()
    podi.Client.UserAgent = "Podimator"
    podi.Process()
}
